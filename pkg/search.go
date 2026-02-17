package libra

import (
	"runtime"
	"sync"
	"time"
)

const (
	SearchMaxDepth      = 16        // Maximum depth to search
	MaxEvaluationScore  = 1_000_000 // Maximum score for wining
	MaxEvaluationTimeMs = 3_000     // Maximum time for a search at the root level
)

type SearchOptions struct {
	MaxDepth           int                 // Maximum search depth (plies)
	TimeLimitInMs      int                 // Soft time limit: stop deepening after this (ms)
	MaxTimeLimitInMs   int                 // Hard time limit: abort in-flight search after this (ms)
	TranspositionTable *TranspositionTable // Optional transposition table to use for search
	UseBookMoves       bool                // Optional flag to use book moves
	StopChan           chan struct{}        // External stop signal (e.g. UCI "stop" command)
}

func (board *Board) IterativeDeepeningSearch(options SearchOptions) *Move {
	tt := options.TranspositionTable
	if tt == nil {
		tt = NewTranspositionTable()
	}

	maxDepth := SearchMaxDepth
	if options.MaxDepth != 0 {
		maxDepth = options.MaxDepth
	}

	// Soft limit: stop starting new depths after this
	softLimit := options.TimeLimitInMs
	// Hard limit: abort in-flight search after this
	hardLimit := options.MaxTimeLimitInMs

	// Fall back to defaults when no time info is provided
	if softLimit == 0 && hardLimit == 0 && options.MaxDepth == 0 {
		softLimit = MaxEvaluationTimeMs
		hardLimit = MaxEvaluationTimeMs
	}

	// If only one is set, use it for both
	if softLimit == 0 && hardLimit > 0 {
		softLimit = hardLimit
	}
	if hardLimit == 0 && softLimit > 0 {
		hardLimit = softLimit
	}

	var bestMove *Move
	totalTimeSpentInMs := 0
	// Iterative deepening
	for depth := 1; depth <= maxDepth; depth++ {
		// Stop deepening if we've exceeded the soft limit
		if softLimit > 0 && totalTimeSpentInMs >= softLimit {
			break
		}
		// Give in-flight search up to the hard limit remaining
		searchTimeLimit := 0
		if hardLimit > 0 {
			searchTimeLimit = hardLimit - totalTimeSpentInMs
			if searchTimeLimit <= 0 {
				break
			}
		}
		result := board.Search(depth, tt, searchTimeLimit, bestMove, options.StopChan)
		result.PrintUCI()
		if result.BestMove != nil && (!result.IsInterrupted || bestMove == nil) {
			bestMove = result.BestMove
		}
		totalTimeSpentInMs += int(result.TimeSpentInMs)
		// If search was interrupted (timeout or stop), don't start next depth
		if result.IsInterrupted {
			break
		}
	}

	return bestMove
}

func (board *Board) Search(depth int, tt *TranspositionTable, timeLimitInMs int, pvMove *Move, stopChan chan struct{}) *SearchResult {
	result := &SearchResult{}
	result.StartTimer()
	result.SetMaxSearchDepth(int32(depth))
	moves := board.GenerateLegalMoves()
	result.IncMoveGeneration()
	ttMove := tt.BestMoveDeepest(board.ZobristHash())
	moves = board.SortMovesRoot(moves, pvMove, ttMove)
	ctx := &SearchContext{Done: make(chan struct{})}

	var score int
	var move *Move
	finished := make(chan struct{})
	go func() {
		score, move = board.ParallelRootSearch(depth, tt, moves, result, ctx)
		close(finished)
	}()

	if timeLimitInMs > 0 {
		// Timed search: respect both timeout and external stop
		timeout := time.After(time.Duration(timeLimitInMs) * time.Millisecond)
		select {
		case <-timeout:
			close(ctx.Done)
			<-finished
			result.IsInterrupted = true
		case <-finished:
		case <-stopChan:
			close(ctx.Done)
			<-finished
			result.IsInterrupted = true
		}
	} else if stopChan != nil {
		// Infinite/depth-only: wait for finish or external stop
		select {
		case <-finished:
		case <-stopChan:
			close(ctx.Done)
			<-finished
			result.IsInterrupted = true
		}
	} else {
		// No time limit and no stop channel: just wait
		<-finished
	}

	result.BestScore = score
	result.StopTimer()
	result.BestMove = move
	return result
}

type ConcurrentSearch struct {
	score         int
	move          Move
	originalIndex int
}

// ParallelRootSearch allows passing in a pre-sorted move list
func (board *Board) ParallelRootSearch(depth int, tt *TranspositionTable, moves []Move, stats *SearchResult, ctx *SearchContext) (int, *Move) {
	maximizing := board.WhiteToMove
	bestScore := -MaxEvaluationScore
	if !maximizing {
		bestScore = MaxEvaluationScore
	}
	if len(moves) == 0 {
		return board.MateOrStalemateScore(maximizing), nil
	}

	numWorkers := runtime.GOMAXPROCS(0)
	moveChan := make(chan struct {
		move  Move
		index int
	}, len(moves))
	resultChan := make(chan ConcurrentSearch, len(moves))
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range moveChan {
				if runtime.GOARCH == "wasm" {
					runtime.Gosched()
				}
				select {
				case <-ctx.Done:
					return
				default:
				}
				clone := board.Clone()
				clone.Move(job.move)
				score := clone.AlphaBetaSearch(
					depth-1, !maximizing,
					-MaxEvaluationScore, MaxEvaluationScore, tt, stats, ctx, 1,
				)
				resultChan <- ConcurrentSearch{score: score, move: job.move, originalIndex: job.index}
			}
		}()
	}

	// Send jobs
	for i, m := range moves {
		moveChan <- struct {
			move  Move
			index int
		}{move: m, index: i}
	}
	close(moveChan)

	// Wait for workers to finish and close resultChan
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var bestMove *Move
	bestMoveOriginalIndex := -1
	for result := range resultChan {
		if maximizing {
			if result.score > bestScore || (result.score == bestScore && (bestMove == nil || result.originalIndex < bestMoveOriginalIndex)) {
				bestScore = result.score
				tempMove := result.move
				bestMove = &tempMove
				bestMoveOriginalIndex = result.originalIndex
			}
		} else {
			if result.score < bestScore || (result.score == bestScore && (bestMove == nil || result.originalIndex < bestMoveOriginalIndex)) {
				bestScore = result.score
				tempMove := result.move
				bestMove = &tempMove
				bestMoveOriginalIndex = result.originalIndex
			}
		}
	}
	return bestScore, bestMove
}

func (board *Board) QuiescenceSearch(maximizing bool, alpha int, beta int, stats *SearchResult, ctx *SearchContext) int {
	select {
	case <-ctx.Done:
		return 0
	default:
	}

	if runtime.GOARCH == "wasm" {
		runtime.Gosched()
	}

	stats.IncNodesSearched()

	standPat := board.Evaluate()

	if maximizing {
		if standPat >= beta {
			return beta
		}
		if standPat > alpha {
			alpha = standPat
		}
	} else {
		if standPat <= alpha {
			return alpha
		}
		if standPat < beta {
			beta = standPat
		}
	}

	captures := board.GenerateLegalCaptures()
	captures = board.SortCaptures(captures)

	if maximizing {
		for _, move := range captures {
			prev := board.Move(move)
			score := board.QuiescenceSearch(false, alpha, beta, stats, ctx)
			board.UndoMove(prev)
			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				break
			}
		}
		return alpha
	}

	for _, move := range captures {
		prev := board.Move(move)
		score := board.QuiescenceSearch(true, alpha, beta, stats, ctx)
		board.UndoMove(prev)
		if score < beta {
			beta = score
		}
		if beta <= alpha {
			break
		}
	}
	return beta
}

func (board *Board) AlphaBetaSearch(depth int, maximizing bool, alpha int, beta int, tt *TranspositionTable, stats *SearchResult, ctx *SearchContext, ply int) int {
	// Check for cancellation at every node
	select {
	case <-ctx.Done:
		return 0 // or a sentinel value indicating cancellation
	default:
	}

	// Yield to scheduler in WASM to allow cancellation
	if runtime.GOARCH == "wasm" {
		runtime.Gosched()
	}

	stats.IncNodesSearched()

	if depth == 0 {
		return board.QuiescenceSearch(maximizing, alpha, beta, stats, ctx)
	}

	hash := board.ZobristHash()
	if entry, ok := tt.Get(hash, depth); ok {
		stats.IncTTHit()
		switch entry.Bound {
		case BoundExact:
			return entry.Score
		case BoundLower:
			if entry.Score >= beta {
				return entry.Score
			}
			if entry.Score > alpha && maximizing {
				alpha = entry.Score
			}
		case BoundUpper:
			if entry.Score <= alpha {
				return entry.Score
			}
			if entry.Score < beta && !maximizing {
				beta = entry.Score
			}
		}
	}

	moves := board.GenerateLegalMoves()
	stats.IncMoveGeneration()
	moves = board.SortMovesAlphaBeta(moves, depth, tt, hash, ctx, ply)
	if len(moves) == 0 {
		return board.MateOrStalemateScore(maximizing)
	}

	origAlpha := alpha
	origBeta := beta
	var result int
	var bestMove Move
	if maximizing {
		maxEval := -MaxEvaluationScore
		for i, move := range moves {
			if runtime.GOARCH == "wasm" {
				runtime.Gosched()
			}
			prev := board.Move(move)
			eval := board.AlphaBetaSearch(depth-1, false, alpha, beta, tt, stats, ctx, ply+1)
			board.UndoMove(prev)
			if eval > maxEval {
				maxEval = eval
				bestMove = move
			}
			if maxEval > alpha {
				alpha = maxEval
			}
			if beta <= alpha {
				stats.IncBetaCutoff()
				if move.MoveType != MoveCapture && move.MoveType != MovePromotionCapture {
					ctx.AddKillerMove(move, ply)
					ctx.HistoryHeuristic[PieceToHistoryIndex[move.Piece]][move.To] += depth * depth
				}
				nodesPruned := len(moves) - (i + 1)
				for j := 0; j < nodesPruned; j++ {
					stats.IncNodesPruned()
				}
				break
			}
		}
		result = maxEval
	} else {
		minEval := MaxEvaluationScore
		for i, move := range moves {
			if runtime.GOARCH == "wasm" {
				runtime.Gosched()
			}
			prev := board.Move(move)
			eval := board.AlphaBetaSearch(depth-1, true, alpha, beta, tt, stats, ctx, ply+1)
			board.UndoMove(prev)
			if eval < minEval {
				minEval = eval
				bestMove = move
			}
			if minEval < beta {
				beta = minEval
			}
			if beta <= alpha {
				stats.IncBetaCutoff()
				if move.MoveType != MoveCapture && move.MoveType != MovePromotionCapture {
					ctx.AddKillerMove(move, ply)
					ctx.HistoryHeuristic[PieceToHistoryIndex[move.Piece]][move.To] += depth * depth
				}
				nodesPruned := len(moves) - (i + 1)
				for j := 0; j < nodesPruned; j++ {
					stats.IncNodesPruned()
				}
				break
			}
		}
		result = minEval
	}

	// Determine bound type
	var bound byte = BoundExact
	if maximizing {
		if result <= origAlpha {
			bound = BoundUpper
		} else if result >= beta {
			bound = BoundLower
		}
	} else {
		if result >= origBeta {
			bound = BoundLower
		} else if result <= alpha {
			bound = BoundUpper
		}
	}

	stats.IncTTStore()
	tt.Set(hash, depth, result, bestMove, bound)
	return result
}
