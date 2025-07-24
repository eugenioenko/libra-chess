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
	NullMoveReduction   = 3         // Reduction for null move pruning (R)
	NullMoveMinDepth    = 3         // Minimum depth for null mov
)

type SearchOptions struct {
	MaxDepth           int                 // Maximum search depth (plies)
	RemainingTimeInMs  int                 // Time remaining for the game
	TimeLimitInMs      int                 // Maximum time allowed for the search. Defaults MaxEvaluationDepthTimeMs
	TranspositionTable *TranspositionTable // Optional transposition table to use for search
	UseBookMoves       bool                // Optional flag to use book moves
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

	// Limit depth search when running out of time
	if options.RemainingTimeInMs != 0 && options.RemainingTimeInMs < 2500 {
		maxDepth = 3
	}

	maxSearchTimeMs := MaxEvaluationTimeMs
	if options.TimeLimitInMs != 0 {
		maxSearchTimeMs = options.TimeLimitInMs
	}

	var bestMove *Move
	totalTimeSpentInMs := 0
	// Iterative deepening
	for depth := 1; depth <= maxDepth; depth++ {
		searchTimeLimit := maxSearchTimeMs - totalTimeSpentInMs
		if searchTimeLimit <= 0 {
			break
		}
		result := board.Search(depth, tt, searchTimeLimit, bestMove)
		result.PrintUCI()
		if result.BestMove != nil && (!result.IsInterrupted || bestMove == nil) {
			bestMove = result.BestMove
		}
		totalTimeSpentInMs += int(result.TimeSpentInMs)
	}

	return bestMove
}

func (board *Board) Search(depth int, tt *TranspositionTable, timeLimitInMs int, pvMove *Move) *SearchResult {
	result := &SearchResult{}
	result.StartTimer()
	result.SetMaxSearchDepth(int32(depth))
	moves := board.GenerateLegalMoves()
	result.IncMoveGeneration()
	ttMove := tt.BestMoveDeepest(board.ZobristHash())
	moves = board.SortMovesRoot(moves, pvMove, ttMove)
	ctx := &SearchContext{Done: make(chan struct{})}

	timeoutTime := MaxEvaluationTimeMs * time.Millisecond
	if timeLimitInMs != 0 {
		timeoutTime = time.Duration(time.Duration(timeLimitInMs)) * time.Millisecond
	}
	timeout := time.After(timeoutTime)

	var score int
	var move *Move
	finished := make(chan struct{})
	go func() {
		score, move = board.ParallelRootSearch(depth, tt, moves, result, ctx)
		close(finished)
	}()

	select {
	case <-timeout:
		close(ctx.Done)
		<-finished
		result.IsInterrupted = true
	case <-finished:
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
		return board.Evaluate()
	}

	hash := board.ZobristHash()
	if score, ok := tt.Get(hash, depth); ok {
		stats.IncTTHit()
		return score
	}

	moves := board.GenerateLegalMoves()
	stats.IncMoveGeneration()
	moves = board.SortMovesAlphaBeta(moves, depth, tt, hash, ctx, ply)
	if len(moves) == 0 {
		return board.MateOrStalemateScore(maximizing)
	}

	// --- Null Move Pruning ---
	if depth >= NullMoveMinDepth && ply > 0 && !board.IsInCheck(maximizing) && len(board.GenerateLegalMoves()) > 0 {
		nullBoard := board.Clone()
		nullBoard.WhiteToMove = !nullBoard.WhiteToMove // Switch side to move (null move)
		nullEval := 0
		newDepth := depth - NullMoveReduction - 1
		if newDepth < 1 {
			newDepth = 1
		}
		if maximizing {
			nullEval = nullBoard.AlphaBetaSearch(newDepth, false, alpha, beta, tt, stats, ctx, ply+1)
			if nullEval >= beta {
				return beta // Fail-hard beta cutoff
			}
		} else {
			nullEval = nullBoard.AlphaBetaSearch(newDepth, true, alpha, beta, tt, stats, ctx, ply+1)
			if nullEval <= alpha {
				return alpha // Fail-hard alpha cutoff
			}
		}
	}

	// --- Alpha-Beta Pruning ---
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
					// Update history heuristic for quiet moves
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
					// Update history heuristic for quiet moves
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

	stats.IncTTStore()
	tt.Set(hash, depth, result, bestMove)
	return result
}
