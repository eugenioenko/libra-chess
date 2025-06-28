package libra

import (
	"runtime"
	"sync"
)

const (
	MaxEvaluationScore = 1000000
	MaxSearchDepth     = 64 // or whatever your engine's max depth is
)

type searchResult struct {
	score         int
	move          Move
	originalIndex int
}

func (board *Board) Search(maxDepth int, tt *TranspositionTable) (*Move, *SearchStats) {
	stats := &SearchStats{}
	stats.StartTimer()
	defer stats.StopTimer()

	maximizing := board.WhiteToMove
	var bestMove *Move
	bestScore := -MaxEvaluationScore
	if !maximizing {
		bestScore = MaxEvaluationScore
	}

	var pvMove *Move
	for depth := 1; depth <= maxDepth; depth++ {
		stats.SetMaxSearchDepth(int32(depth))
		moves := board.GenerateLegalMoves()
		stats.IncMoveGeneration()
		// Use SortMovesRoot at root
		ttMove := tt.BestMoveDeepest(board.ZobristHash())
		moves = board.SortMovesRoot(moves, pvMove, ttMove)
		ctx := &SearchContext{} // new context for each root search
		score, move := board.ParallelRootSearch(depth, tt, moves, stats, ctx)
		bestScore = score
		bestMove = move
		pvMove = move // update PV move for next iteration
	}
	stats.BestScore = bestScore
	return bestMove, stats
}

// ParallelRootSearch allows passing in a pre-sorted move list
func (board *Board) ParallelRootSearch(depth int, tt *TranspositionTable, moves []Move, stats *SearchStats, ctx *SearchContext) (int, *Move) {
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
	resultChan := make(chan searchResult, len(moves))
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range moveChan {
				clone := board.Clone()
				clone.Move(job.move)
				score := clone.AlphaBetaSearch(
					depth-1, !maximizing,
					-MaxEvaluationScore, MaxEvaluationScore, tt, stats, ctx, 1,
				)
				resultChan <- searchResult{score: score, move: job.move, originalIndex: job.index}
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
	var bestMoveOriginalIndex int = -1
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

func (board *Board) AlphaBetaSearch(depth int, maximizing bool, alpha int, beta int, tt *TranspositionTable, stats *SearchStats, ctx *SearchContext, ply int) int {
	stats.IncNodesSearched()
	if depth == 0 {
		return board.Evaluate()
	}

	hash := board.ZobristHash()
	if entry, ok := tt.Get(hash, depth); ok {
		stats.IncTTHit()
		return entry
	}

	moves := board.GenerateLegalMoves()
	stats.IncMoveGeneration()
	moves = board.SortMovesAlphaBeta(moves, depth, tt, hash, ctx, ply)
	if len(moves) == 0 {
		return board.MateOrStalemateScore(maximizing)
	}

	var result int
	var bestMove Move
	if maximizing {
		maxEval := -MaxEvaluationScore
		for i, move := range moves {
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
