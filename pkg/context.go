package libra

// SearchContext holds per-search context (killer moves, history heuristic, etc.)
type SearchContext struct {
	KillerMoves [MaxSearchDepth][2]Move
	// HistoryHeuristic[piece][toSquare] (simple version)
	HistoryHeuristic [16][64]int // 16 piece types, 64 squares
}

// IsKillerMove returns true if the move is a killer move at the given ply
func (info *SearchContext) IsKillerMove(move Move, ply int) bool {
	return (move == info.KillerMoves[ply][0]) || (move == info.KillerMoves[ply][1])
}

// AddKillerMove adds a move to the killer moves for the given ply (if it's not already present)
func (info *SearchContext) AddKillerMove(move Move, ply int) {
	if move == info.KillerMoves[ply][0] || move == info.KillerMoves[ply][1] {
		return
	}
	info.KillerMoves[ply][1] = info.KillerMoves[ply][0]
	info.KillerMoves[ply][0] = move
}
