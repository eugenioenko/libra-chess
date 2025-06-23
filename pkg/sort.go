package libra

import "sort"

func (board *Board) SortMoves(moves []Move, depth int, tt *TranspositionTable) []Move {
	if len(moves) == 0 {
		return moves
	}

	// Sort moves by MoveType preferring captures, then by From, To, and promotion piece for full determinism
	sort.Slice(moves, func(i, j int) bool {
		moveA := moves[i]
		moveB := moves[j]

		isCaptureA := moveA.MoveType == MoveCapture || moveA.MoveType == MovePromotionCapture
		isCaptureB := moveB.MoveType == MoveCapture || moveB.MoveType == MovePromotionCapture
		if isCaptureA != isCaptureB {
			return isCaptureA
		}

		isPromoA := moveA.MoveType == MovePromotion
		isPromoB := moveB.MoveType == MovePromotion
		if isPromoA != isPromoB {
			return isPromoA
		}

		// Sort by capture value if both moves are captures
		// This ensures that if two captures are available, the one with the higher value piece captured is preferred.
		if isCaptureA && isCaptureB {
			victimA := moveA.Data[0]
			attackerA := board.PieceAtSquare(moveA.From)
			victimB := moveB.Data[0]
			attackerB := board.PieceAtSquare(moveB.From)
			scoreA := PieceCodeToValue[victimA] - PieceCodeToValue[attackerA]
			scoreB := PieceCodeToValue[victimB] - PieceCodeToValue[attackerB]
			if scoreA != scoreB {
				return scoreA > scoreB
			}
		}

		// For promotions, ensure consistent order by promotion piece
		if moveA.MoveType == MovePromotion || moveA.MoveType == MovePromotionCapture {
			if moveA.Data[0] != moveB.Data[0] {
				// For promotions, sort by piece value in ascending order: Knight < Bishop < Rook < Queen.
				// This ensures deterministic move ordering, so that when multiple promotions have equal evaluation,
				// the queen promotion (highest value) is preferred if all else is equal.
				return moveA.Data[0] < moveB.Data[0]
			}
		}

		if moveA.From != moveB.From {
			return moveA.From < moveB.From
		}
		return moveA.To < moveB.To
	})

	return moves
}

// SortMovesAlphaBeta orders moves for alpha-beta search using TT, killer moves, MVV-LVA, and history heuristic.
// ttMove: best move from transposition table (if any)
func (board *Board) SortMovesAlphaBeta(
	moves []Move,
	depth int,
	tt *TranspositionTable,
	hash uint64,
) []Move {
	type moveScore struct {
		move  Move
		score int
	}
	scored := make([]moveScore, len(moves))

	for i, m := range moves {
		score := 0

		// 1. Transposition Table move gets highest priority
		ttBestMove := tt.BestMoveDeepest(hash)
		if ttBestMove != nil && m == *ttBestMove {
			score += 1000000
		}

		// 2. MVV-LVA for captures
		if m.MoveType == MoveCapture || m.MoveType == MovePromotionCapture {
			victim := m.Data[0]
			attacker := board.PieceAtSquare(m.From)
			score += 10000 + 100*PieceCodeToValue[victim] - PieceCodeToValue[attacker]
		}

		scored[i] = moveScore{move: m, score: score}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	for i := range moves {
		moves[i] = scored[i].move
	}
	return moves
}

// SortMovesRoot orders root moves using previous PV move, TT move, and MVV-LVA
// pvMove: principal variation move from previous iteration (if any)
// ttMove: best move from transposition table (if any)
func (board *Board) SortMovesRoot(
	moves []Move,
	pvMove *Move,
	ttMove *Move,
) []Move {
	type moveScore struct {
		move  Move
		score int
	}
	scored := make([]moveScore, len(moves))

	for i, m := range moves {
		score := 0

		// 1. Previous PV move gets highest priority
		if pvMove != nil && m == *pvMove {
			score += 1000000
		}

		// 2. Transposition Table move (if not PV)
		if ttMove != nil && m == *ttMove {
			score += 900000
		}

		// 3. MVV-LVA for captures
		if m.MoveType == MoveCapture || m.MoveType == MovePromotionCapture {
			victim := m.Data[0]
			attacker := board.PieceAtSquare(m.From)
			score += 10000 + 100*PieceCodeToValue[victim] - PieceCodeToValue[attacker]
		}

		scored[i] = moveScore{move: m, score: score}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	for i := range moves {
		moves[i] = scored[i].move
	}
	return moves
}
