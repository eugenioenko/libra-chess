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

// EvaluatePawnStructure evaluates pawn structure for both sides: doubled, isolated, and passed pawns.
func (board *Board) EvaluatePawnStructure() (int, int) {
	whitePenalty := 0
	blackPenalty := 0
	whiteBonus := 0
	blackBonus := 0

	// Helper: get file mask for a file
	fileMask := func(file int) uint64 {
		return 0x0101010101010101 << file
	}

	// Doubled pawns
	for file := 0; file < 8; file++ {
		wPawnsOnFile := bits.OnesCount64(board.WhitePawns & fileMask(file))
		bPawnsOnFile := bits.OnesCount64(board.BlackPawns & fileMask(file))
		if wPawnsOnFile > 1 {
			whitePenalty += (wPawnsOnFile - 1) * 10 // Penalty per extra pawn
		}
		if bPawnsOnFile > 1 {
			blackPenalty += (bPawnsOnFile - 1) * 10
		}
	}

	// Isolated pawns
	for file := 0; file < 8; file++ {
		wFilePawns := board.WhitePawns & fileMask(file)
		bFilePawns := board.BlackPawns & fileMask(file)
		adjWhite := uint64(0)
		adjBlack := uint64(0)
		if file > 0 {
			adjWhite |= board.WhitePawns & fileMask(file-1)
			adjBlack |= board.BlackPawns & fileMask(file-1)
		}
		if file < 7 {
			adjWhite |= board.WhitePawns & fileMask(file+1)
			adjBlack |= board.BlackPawns & fileMask(file+1)
		}
		if wFilePawns != 0 && adjWhite == 0 {
			whitePenalty += bits.OnesCount64(wFilePawns) * 15 // Penalty per isolated pawn
		}
		if bFilePawns != 0 && adjBlack == 0 {
			blackPenalty += bits.OnesCount64(bFilePawns) * 15
		}
	}

	// Passed pawns
	for file := 0; file < 8; file++ {
		wPawns := board.WhitePawns & fileMask(file)
		bPawns := board.BlackPawns & fileMask(file)
		for wPawns != 0 {
			sq := bits.TrailingZeros64(wPawns)
			rank := sq / 8
			// Passed if no black pawns on same or adjacent files ahead
			mask := uint64(0)
			for f := MathMinByte(byte(file+1), 7); f >= MathMinByte(byte(file-1), 0) && f <= 7; f-- {
				mask |= board.BlackPawns & fileMask(int(f))
			}
			mask &= ^((1 << (sq + 1)) - 1) // Only pawns ahead
			if mask == 0 {
				whiteBonus += (7 - rank) * 12 // More bonus as pawn advances
			}
			wPawns &= wPawns - 1
		}
		for bPawns != 0 {
			sq := bits.TrailingZeros64(bPawns)
			rank := sq / 8
			mask := uint64(0)
			for f := MathMinByte(byte(file+1), 7); f >= MathMinByte(byte(file-1), 0) && f <= 7; f-- {
				mask |= board.WhitePawns & fileMask(int(f))
			}
			mask &= (1 << sq) - 1 // Only pawns ahead for black
			if mask == 0 {
				blackBonus += rank * 12
			}
			bPawns &= bPawns - 1
		}
	}

	return whiteBonus - whitePenalty, blackBonus - blackPenalty
}

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

// EvaluateKingSafety evaluates king safety for both sides in the middlegame.
func (board *Board) EvaluateKingSafety() (int, int) {
	whitePenalty := 0
	blackPenalty := 0

	// Only apply in middlegame (if both sides have queens or enough material)
	material := 0
	material += bits.OnesCount64(board.WhitePawns) + bits.OnesCount64(board.BlackPawns)
	material += bits.OnesCount64(board.WhiteKnights)*3 + bits.OnesCount64(board.BlackKnights)*3
	material += bits.OnesCount64(board.WhiteBishops)*3 + bits.OnesCount64(board.BlackBishops)*3
	material += bits.OnesCount64(board.WhiteRooks)*5 + bits.OnesCount64(board.BlackRooks)*5
	material += bits.OnesCount64(board.WhiteQueens)*9 + bits.OnesCount64(board.BlackQueens)*9
	if material <= 20 { // skip in endgame
		return 0, 0
	}

	// Helper: count friendly pawns in 3x3 area around king
	countKingPawnShield := func(kingSq int, pawns uint64, isWhite bool) int {
		file := kingSq % 8
		rank := kingSq / 8
		count := 0
		for df := -1; df <= 1; df++ {
			for dr := 0; dr <= 1; dr++ { // only in front and same rank
				f := file + df
				r := rank + dr*func() int {
					if isWhite {
						return -1
					} else {
						return 1
					}
				}()
				if f < 0 || f > 7 || r < 0 || r > 7 {
					continue
				}
				sq := r*8 + f
				if (pawns & (1 << sq)) != 0 {
					count++
				}
			}
		}
		return count
	}

	if board.WhiteKing != 0 {
		wKingSq := bits.TrailingZeros64(board.WhiteKing)
		shield := countKingPawnShield(wKingSq, board.WhitePawns, true)
		whitePenalty -= shield * 12       // reward for pawn cover
		whitePenalty += (3 - shield) * 18 // penalty for missing pawns
	}
	if board.BlackKing != 0 {
		bKingSq := bits.TrailingZeros64(board.BlackKing)
		shield := countKingPawnShield(bKingSq, board.BlackPawns, false)
		blackPenalty -= shield * 12
		blackPenalty += (3 - shield) * 18
	}

	// TODO: Add open file and enemy piece proximity checks for more accuracy
	return whitePenalty, blackPenalty
}

// EvaluateBishopPair returns a bonus for having both bishops.
func (board *Board) EvaluateBishopPair() (int, int) {
	whiteBonus := 0
	blackBonus := 0
	if bits.OnesCount64(board.WhiteBishops) >= 2 {
		whiteBonus = 35 // typical value, can be tuned
	}
	if bits.OnesCount64(board.BlackBishops) >= 2 {
		blackBonus = 35
	}
	return whiteBonus, blackBonus
}

// EvaluateCenterControl rewards control of the central squares (d4, d5, e4, e5).
func (board *Board) EvaluateCenterControl() (int, int) {
	centerSquares := [4]byte{27, 28, 35, 36} // d4, e4, d5, e5 (0-based)
	whiteControl := 0
	blackControl := 0

	for _, sq := range centerSquares {
		if board.IsSquareAttacked(sq, true) {
			whiteControl++
		}
		if board.IsSquareAttacked(sq, false) {
			blackControl++
		}
	}
	// Each control gets a bonus (tune as needed)
	return whiteControl * 10, blackControl * 10
}

func (board *Board) Evaluate() int {
	whiteScore, blackScore := board.EvaluateMaterialAndPST()
	pawnStructWhite, pawnStructBlack := board.EvaluatePawnStructure()
	kingSafeWhite, kingSafeBlack := board.EvaluateKingSafety()
	bishopPairWhite, bishopPairBlack := board.EvaluateBishopPair()
	centerWhite, centerBlack := board.EvaluateCenterControl()
	whiteScore += pawnStructWhite + kingSafeWhite + bishopPairWhite + centerWhite
	blackScore += pawnStructBlack + kingSafeBlack + bishopPairBlack + centerBlack
	return whiteScore - blackScore
}
