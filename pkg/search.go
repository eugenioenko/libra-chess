package libra

import "sync"

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

	var bestMove *Move
	var bestMoveOriginalIndex int = -1 // Initialize with an invalid index
	numMoves := len(moves)
	moveChan := make(chan searchResult, numMoves)
	var wg sync.WaitGroup

	for i, currentMove := range moves {
		wg.Add(1)
		go func(m Move, index int) {
			defer wg.Done()
			clone := board.Clone()
			clone.Move(m)
			score := clone.AlphaBeta(
				depth-1, !maximizing,
				-MaxEvaluationScore, MaxEvaluationScore, tt,
			)
			moveChan <- searchResult{score: score, move: m, originalIndex: index}
		}(currentMove, i)
	}

	go func() {
		wg.Wait()
		close(moveChan)
	}()

	for result := range moveChan {
		if maximizing {
			if result.score > bestScore || (result.score == bestScore && (bestMove == nil || result.originalIndex < bestMoveOriginalIndex)) {
				bestScore = result.score
				// Assign a copy of the move from the result
				tempMove := result.move
				bestMove = &tempMove
				bestMoveOriginalIndex = result.originalIndex
			}
		} else {
			if result.score < bestScore || (result.score == bestScore && (bestMove == nil || result.originalIndex < bestMoveOriginalIndex)) {
				bestScore = result.score
				// Assign a copy of the move from the result
				tempMove := result.move
				bestMove = &tempMove
				bestMoveOriginalIndex = result.originalIndex
			}
		}
	}

	// Print debug info
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
