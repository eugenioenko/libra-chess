package libra

import (
	"math/bits"
)

// mirrorIndex mirrors a square index for black's perspective
func mirrorIndex(idx byte) byte {
	return 56 ^ (idx & 56) | (idx & 7)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// EvaluateMaterialAndPST evaluates using tapered PeSTO piece-square tables.
// Returns separate white and black scores after phase interpolation.
func (board *Board) EvaluateMaterialAndPST() (int, int) {
	mgWhite, mgBlack := 0, 0
	egWhite, egBlack := 0, 0
	phase := 0

	// Pawns (phase weight = 0, so no phase contribution)
	for bb := board.WhitePawns; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgWhite += mgPawnPST[sq]
		egWhite += egPawnPST[sq]
		bb &= bb - 1
	}
	for bb := board.BlackPawns; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgBlack += mgPawnPST[mirrorIndex(byte(sq))]
		egBlack += egPawnPST[mirrorIndex(byte(sq))]
		bb &= bb - 1
	}

	// Knights
	for bb := board.WhiteKnights; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgWhite += mgKnightPST[sq]
		egWhite += egKnightPST[sq]
		phase += KnightPhase
		bb &= bb - 1
	}
	for bb := board.BlackKnights; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgBlack += mgKnightPST[mirrorIndex(byte(sq))]
		egBlack += egKnightPST[mirrorIndex(byte(sq))]
		phase += KnightPhase
		bb &= bb - 1
	}

	// Bishops
	for bb := board.WhiteBishops; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgWhite += mgBishopPST[sq]
		egWhite += egBishopPST[sq]
		phase += BishopPhase
		bb &= bb - 1
	}
	for bb := board.BlackBishops; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgBlack += mgBishopPST[mirrorIndex(byte(sq))]
		egBlack += egBishopPST[mirrorIndex(byte(sq))]
		phase += BishopPhase
		bb &= bb - 1
	}

	// Rooks
	for bb := board.WhiteRooks; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgWhite += mgRookPST[sq]
		egWhite += egRookPST[sq]
		phase += RookPhase
		bb &= bb - 1
	}
	for bb := board.BlackRooks; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgBlack += mgRookPST[mirrorIndex(byte(sq))]
		egBlack += egRookPST[mirrorIndex(byte(sq))]
		phase += RookPhase
		bb &= bb - 1
	}

	// Queens
	for bb := board.WhiteQueens; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgWhite += mgQueenPST[sq]
		egWhite += egQueenPST[sq]
		phase += QueenPhase
		bb &= bb - 1
	}
	for bb := board.BlackQueens; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		mgBlack += mgQueenPST[mirrorIndex(byte(sq))]
		egBlack += egQueenPST[mirrorIndex(byte(sq))]
		phase += QueenPhase
		bb &= bb - 1
	}

	// Kings (no phase contribution)
	if board.WhiteKing != 0 {
		sq := bits.TrailingZeros64(board.WhiteKing)
		mgWhite += mgKingPST[sq]
		egWhite += egKingPST[sq]
	}
	if board.BlackKing != 0 {
		sq := bits.TrailingZeros64(board.BlackKing)
		mgBlack += mgKingPST[mirrorIndex(byte(sq))]
		egBlack += egKingPST[mirrorIndex(byte(sq))]
	}

	// Clamp phase to TotalPhase
	if phase > TotalPhase {
		phase = TotalPhase
	}

	// Tapered interpolation: phase=TotalPhase means full middlegame, phase=0 means full endgame
	mgScore := mgWhite - mgBlack
	egScore := egWhite - egBlack
	whiteScore := (mgScore*phase + egScore*(TotalPhase-phase)) / TotalPhase

	// Endgame king proximity heuristic
	if phase <= 6 && board.WhiteKing != 0 && board.BlackKing != 0 {
		wKing := byte(bits.TrailingZeros64(board.WhiteKing))
		bKing := byte(bits.TrailingZeros64(board.BlackKing))
		wRank := int(wKing / 8)
		wFile := int(wKing % 8)
		bRank := int(bKing / 8)
		bFile := int(bKing % 8)
		dist := abs(wRank-bRank) + abs(wFile-bFile)
		if egWhite > egBlack {
			whiteScore += (14 - dist) * 10
		} else if egBlack > egWhite {
			whiteScore -= (14 - dist) * 10
		}
	}

	return whiteScore, 0
}

// Mobility: count the number of legal moves for each side
func (board *Board) MateOrStalemateScore(maximizing bool) int {
	kingSq := board.ActiveKingSquare()
	if board.IsSquareAttacked(kingSq, board.WhiteToMove) {
		if maximizing {
			return -MaxEvaluationScore
		} else {
			return MaxEvaluationScore
		}
	} else {
		return 0
	}
}

func (board *Board) Evaluate() int {
	whiteScore, blackScore := board.EvaluateMaterialAndPST()

	return whiteScore - blackScore
}
