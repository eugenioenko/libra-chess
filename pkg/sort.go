package libra

import "sort"

/*
	SortMovesAlphaBeta orders moves for alpha-beta search using TT, killer moves, MVV-LVA, and history heuristic.

	# SortMovesAlphaBeta Scoring Table

| Move Type      | Formula                               | Min   | Max   |
|----------------|---------------------------------------|-------|-------|
| TT Move        | 70_000                                |70_000 |70_000 |
| Killer Move    | 60_000                                |60_000 |60_000 |
| Promo Capture  | 50_000+10Victim+10Promo-10Attacker    |48_000 |67_000 |
| Capture        | 30_000+10Victim-10Attacker            |29_000 |38_000 |
| Promo (quiet)  | 10_000+10Promo                        |11_000 |19_000 |
| History        | history[Piece][To]                    |   0   | < 6k  |
| Quiet          | 0                                     |   0   |   0   |
*/
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
			score += 90_0000 // High value for TT move
		}

		// 2. MVV-LVA for captures, now with SEE
		switch m.MoveType {
		case MoveCapture, MovePromotionCapture:
			see := board.SEE(m)
			victim := m.Captured
			attacker := m.Piece
			if see >= 0 {
				// Good capture: high score
				score += 70_000 + 10*PieceCodeToValue[victim] - 10*PieceCodeToValue[attacker] + 1000*see
			} else {
				// Bad capture: demote, but still above quiets
				score += 10_000 + see
			}
		case MovePromotion:
			promoPiece := m.Promoted
			score += 30_000 + 10*PieceCodeToValue[promoPiece]
		}

		// 3. Killer moves
		if ctx != nil && ctx.IsKillerMove(m, ply) {
			score += 10_000
		}

		// 4. History heuristic for quiet moves
		// Value is set with depth^2, so deeper moves are prioritized
		// Example with max depth 9 it can reach 810
		if ctx != nil && m.IsQuiet() {
			score += ctx.HistoryHeuristic[PieceToHistoryIndex[m.Piece]][m.To] * 10
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

/*
	SortMovesRoot orders root moves using previous PV move, TT move, and MVV-LVA
	pvMove: principal variation move from previous iteration (if any)
	ttMove: best move from transposition table (if any)

# SortMovesRoot Scoring Table

| Move Type      | Formula                         | Min   | Max   |
|----------------|---------------------------------|-------|-------|
| PV Move        | 90_000                          |90_000 |90_000 |
| TT Move        | 70_000                          |70_000 |70_000 |
| Promo Capture  | 50k+10Victim+10Promo-10Attacker |48_000 |67_000 |
| Capture        | 30k+10Victim-10Attacker         |29_000 |38_000 |
| Promo (quiet)  | 10k+10Promo                     |11_000 |19_000 |
| Quiet/Other    | 0                               |   0   |   0   |
*/
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
			score += 90_000
		}

		// 2. Transposition Table move (if not PV)
		if ttMove != nil && m == *ttMove {
			score += 70_000
		}

		// 3. MVV-LVA for captures, now with SEE
		switch m.MoveType {
		case MovePromotionCapture, MoveCapture:
			see := board.SEE(m)
			victim := m.Captured
			attacker := m.Piece
			if see >= 0 {
				// Good capture: high score
				if m.MoveType == MovePromotionCapture {
					promoPiece := m.Promoted
					score += 50_000 + 10*PieceCodeToValue[victim] + 10*PieceCodeToValue[promoPiece] - 10*PieceCodeToValue[attacker] + 1000*see
				} else {
					score += 30_000 + 10*PieceCodeToValue[victim] - 10*PieceCodeToValue[attacker] + 1000*see
				}
			} else {
				// Bad capture: demote, but still above quiets
				score += 10_000 + see
			}
		case MovePromotion:
			promoPiece := m.Promoted
			score += 10_000 + 10*PieceCodeToValue[promoPiece]
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

// SEE (Static Exchange Evaluation) estimates the net material gain/loss for a capture move.
// Returns the net gain (positive for winning, negative for losing, 0 for even).
func (board *Board) SEE(move Move) int {
	// Only evaluate captures
	if move.MoveType != MoveCapture && move.MoveType != MovePromotionCapture {
		return 0
	}
	victim := move.Captured
	attacker := move.Piece
	// Simple SEE: value of captured piece minus value of attacker
	// (does not recurse, but is fast and good enough for move ordering)
	return PieceCodeToValue[victim] - PieceCodeToValue[attacker]
}
