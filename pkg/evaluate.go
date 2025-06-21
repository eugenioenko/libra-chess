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

// EvaluateMaterialAndPST evaluates the material and piece-square table (PST) scores for both sides.
func (board *Board) EvaluateMaterialAndPST() (int, int) {
	whiteScore := 0
	blackScore := 0
	// Pawns
	for bb := board.WhitePawns; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		whiteScore += PieceCodeToValue[WhitePawn]
		whiteScore += pawnPST[sq]
		bb &= bb - 1
	}
	for bb := board.BlackPawns; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		blackScore += PieceCodeToValue[BlackPawn]
		blackScore += pawnPST[mirrorIndex(byte(sq))]
		bb &= bb - 1
	}
	// Knights
	for bb := board.WhiteKnights; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		whiteScore += PieceCodeToValue[WhiteKnight]
		whiteScore += knightPST[sq]
		bb &= bb - 1
	}
	for bb := board.BlackKnights; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		blackScore += PieceCodeToValue[BlackKnight]
		blackScore += knightPST[mirrorIndex(byte(sq))]
		bb &= bb - 1
	}
	// Bishops
	for bb := board.WhiteBishops; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		whiteScore += PieceCodeToValue[WhiteBishop]
		whiteScore += bishopPST[sq]
		bb &= bb - 1
	}
	for bb := board.BlackBishops; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		blackScore += PieceCodeToValue[BlackBishop]
		blackScore += bishopPST[mirrorIndex(byte(sq))]
		bb &= bb - 1
	}
	// Rooks
	for bb := board.WhiteRooks; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		whiteScore += PieceCodeToValue[WhiteRook]
		whiteScore += rookPST[sq]
		bb &= bb - 1
	}
	for bb := board.BlackRooks; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		blackScore += PieceCodeToValue[BlackRook]
		blackScore += rookPST[mirrorIndex(byte(sq))]
		bb &= bb - 1
	}
	// Queens
	for bb := board.WhiteQueens; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		whiteScore += PieceCodeToValue[WhiteQueen]
		whiteScore += queenPST[sq]
		bb &= bb - 1
	}
	for bb := board.BlackQueens; bb != 0; {
		sq := bits.TrailingZeros64(bb)
		blackScore += PieceCodeToValue[BlackQueen]
		blackScore += queenPST[mirrorIndex(byte(sq))]
		bb &= bb - 1
	}
	// King
	if board.WhiteKing != 0 {
		sq := bits.TrailingZeros64(board.WhiteKing)
		whiteScore += PieceCodeToValue[WhiteKing]
		whiteScore += kingPST[sq]
	}
	if board.BlackKing != 0 {
		sq := bits.TrailingZeros64(board.BlackKing)
		blackScore += PieceCodeToValue[BlackKing]
		blackScore += kingPST[mirrorIndex(byte(sq))]
	}

	// Encourage mating the king in the endgame
	material := 0
	material += bits.OnesCount64(board.WhitePawns) + bits.OnesCount64(board.BlackPawns)
	material += bits.OnesCount64(board.WhiteKnights)*3 + bits.OnesCount64(board.BlackKnights)*3
	material += bits.OnesCount64(board.WhiteBishops)*3 + bits.OnesCount64(board.BlackBishops)*3
	material += bits.OnesCount64(board.WhiteRooks)*5 + bits.OnesCount64(board.BlackRooks)*5
	material += bits.OnesCount64(board.WhiteQueens)*9 + bits.OnesCount64(board.BlackQueens)*9
	if material <= 14 && board.WhiteKing != 0 && board.BlackKing != 0 {
		wKing := byte(bits.TrailingZeros64(board.WhiteKing))
		bKing := byte(bits.TrailingZeros64(board.BlackKing))
		wRank := int(wKing / 8)
		wFile := int(wKing % 8)
		bRank := int(bKing / 8)
		bFile := int(bKing % 8)
		dist := abs(wRank-bRank) + abs(wFile-bFile)
		wMat := bits.OnesCount64(board.WhiteQueens)*9 + bits.OnesCount64(board.WhiteRooks)*5 + bits.OnesCount64(board.WhiteBishops)*3 + bits.OnesCount64(board.WhiteKnights)*3 + bits.OnesCount64(board.WhitePawns)
		bMat := bits.OnesCount64(board.BlackQueens)*9 + bits.OnesCount64(board.BlackRooks)*5 + bits.OnesCount64(board.BlackBishops)*3 + bits.OnesCount64(board.BlackKnights)*3 + bits.OnesCount64(board.BlackPawns)
		if wMat > bMat {
			whiteScore += (14 - dist) * 10 // Encourage white to approach black king
			blackScore -= (14 - dist) * 10
		} else if bMat > wMat {
			blackScore += (14 - dist) * 10 // Encourage black to approach white king
			whiteScore -= (14 - dist) * 10
		}
	}

	return whiteScore, blackScore
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
