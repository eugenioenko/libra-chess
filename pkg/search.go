package libra

import (
	"runtime"
	"sync"
)

const (
	MaxEvaluationScore = 1000000
)

type searchResult struct {
	score         int
	move          Move
	originalIndex int
}

func (board *Board) Search(depth int, tt *TranspositionTable) (int, *Move) {
	maximizing := board.WhiteToMove
	bestScore := -MaxEvaluationScore
	if !maximizing {
		bestScore = MaxEvaluationScore
	}
	moves := board.GenerateLegalMoves()
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
				score := clone.AlphaBeta(
					depth-1, !maximizing,
					-MaxEvaluationScore, MaxEvaluationScore, tt,
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

func (board *Board) AlphaBeta(depth int, maximizing bool, alpha int, beta int, tt *TranspositionTable) int {
	if depth == 0 {
		return board.Evaluate()
	}

	hash := board.ZobristHash()
	if entry, ok := tt.Get(hash, depth); ok {
		return entry
	}

	moves := board.GenerateLegalMoves()
	if len(moves) == 0 {
		return board.MateOrStalemateScore(maximizing)
	}

	var result int
	if maximizing {
		maxEval := -MaxEvaluationScore
		for _, move := range moves {
			prev := board.Move(move)
			eval := board.AlphaBeta(depth-1, false, alpha, beta, tt)
			board.UndoMove(prev)
			if eval > maxEval {
				maxEval = eval
			}
			if maxEval > alpha {
				alpha = maxEval
			}
			if beta <= alpha {
				break
			}
		}
		result = maxEval
	} else {
		minEval := MaxEvaluationScore
		for _, move := range moves {
			prev := board.Move(move)
			eval := board.AlphaBeta(depth-1, true, alpha, beta, tt)
			board.UndoMove(prev)
			if eval < minEval {
				minEval = eval
			}
			if minEval < beta {
				beta = minEval
			}
			if beta <= alpha {
				break
			}
		}
		result = minEval
	}

	tt.Set(hash, depth, result)
	return result
}
