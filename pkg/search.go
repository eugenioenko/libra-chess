package libra

const (
	MaxEvaluationScore = 1000000
)

type SearchDebug struct {
	EvalCount    int
	TTWriteCount int
	TTHitCount   int
	TTMissCount  int
}

func (board *Board) Search(depth int, tt *TranspositionTable) (int, *Move) {
	debug := &SearchDebug{}
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
	for _, move := range moves {
		clone := board.Clone()
		clone.MakeMove(move)
		score := minimax(
			clone, depth-1, !maximizing,
			-MaxEvaluationScore, MaxEvaluationScore, tt,
			debug,
		)
		if maximizing {
			if score > bestScore || bestMove == nil {
				bestScore = score
				m := move
				bestMove = &m
			}
		} else {
			if score < bestScore || bestMove == nil {
				bestScore = score
				m := move
				bestMove = &m
			}
		}
	}
	// Print debug info
	return bestScore, bestMove
}

func minimax(board *Board, depth int, maximizing bool, alpha int, beta int, tt *TranspositionTable, debug *SearchDebug) int {
	if depth == 0 {
		debug.EvalCount++
		return board.Evaluate()
	}

	hash := board.ZobristHash()
	if entry, ok := tt.Get(hash, depth); ok {
		debug.TTHitCount++
		return entry
	} else {
		debug.TTMissCount++
	}

	moves := board.GenerateLegalMoves()
	if len(moves) == 0 {
		return board.MateOrStalemateScore(maximizing)
	}

	var result int
	if maximizing {
		maxEval := -MaxEvaluationScore
		for _, move := range moves {
			prev := board.MakeMove(move)
			eval := minimax(board, depth-1, false, alpha, beta, tt, debug)
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
			prev := board.MakeMove(move)
			eval := minimax(board, depth-1, true, alpha, beta, tt, debug)
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

	debug.TTWriteCount++
	tt.Set(hash, depth, result)
	return result
}
