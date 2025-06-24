package libra

import "sort"

// SortMoves orders the given moves slice in a deterministic way based on the move's properties and the board state.
// The sorting is done in-place and the function returns the sorted moves slice.
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
			victimA := moveA.Captured
			attackerA := moveA.Piece
			victimB := moveB.Captured
			attackerB := moveB.Piece
			scoreA := PieceCodeToValue[victimA] - PieceCodeToValue[attackerA]
			scoreB := PieceCodeToValue[victimB] - PieceCodeToValue[attackerB]
			if scoreA != scoreB {
				return scoreA > scoreB
			}
		}

		// For promotions, ensure consistent order by promotion piece
		if moveA.MoveType == MovePromotion || moveA.MoveType == MovePromotionCapture {
			if moveA.Promoted != moveB.Promoted {
				// For promotions, sort by piece value in ascending order: Knight < Bishop < Rook < Queen.
				// This ensures deterministic move ordering, so that when multiple promotions have equal evaluation,
				// the queen promotion (highest value) is preferred if all else is equal.
				return moveA.Promoted < moveB.Promoted
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
	ctx *SearchContext,
	ply int,
) []Move {
	type moveScore struct {
		move  Move
		score int
	}
	scored := make([]moveScore, len(moves))
	ttBestMove := tt.BestMoveDeepest(hash)

	for i, m := range moves {
		score := 0

		// 1. Transposition Table move gets highest priority
		if ttBestMove != nil && m == *ttBestMove {
			score += 1000000
		}

		// 2. MVV-LVA for captures
		switch m.MoveType {
		case MovePromotionCapture:
			victim := m.Captured
			attacker := m.Piece
			promoPiece := m.Promoted
			score += 20000 + 100*PieceCodeToValue[victim] + 100*PieceCodeToValue[promoPiece] - 100*PieceCodeToValue[attacker]
		case MoveCapture:
			victim := m.Captured
			attacker := m.Piece
			score += 30000 + 100*PieceCodeToValue[victim] - 100*PieceCodeToValue[attacker]
		case MovePromotion:
			promoPiece := m.Promoted
			score += 10000 + 100*PieceCodeToValue[promoPiece]
		}

		// 3. Killer moves
		if ctx != nil && ctx.IsKillerMove(m, ply) {
			score += 8000
		}

		// 4. History heuristic for quiet moves
		if ctx != nil && m.MoveType == MoveQuiet {
			score += ctx.HistoryHeuristic[PieceToHistoryIndex[m.Piece]][m.To]
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
		switch m.MoveType {
		case MovePromotionCapture:
			victim := m.Captured
			attacker := m.Piece
			promoPiece := m.Promoted
			score += 20000 + 100*PieceCodeToValue[victim] + 1000*PieceCodeToValue[promoPiece] - 100*PieceCodeToValue[attacker]
		case MoveCapture:
			victim := m.Captured
			attacker := m.Piece
			score += 10000 + 100*PieceCodeToValue[victim] - 100*PieceCodeToValue[attacker]
		case MovePromotion:
			promoPiece := m.Promoted
			score += 9000 + 100*PieceCodeToValue[promoPiece]
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
