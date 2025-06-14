package libra

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
	for _, sq := range board.Pieces.White.Pawns {
		whiteScore += PieceCodeToValue[WhitePawn]
		whiteScore += pawnPST[sq]
	}
	for _, sq := range board.Pieces.Black.Pawns {
		blackScore += PieceCodeToValue[BlackPawn]
		blackScore += pawnPST[mirrorIndex(sq)]
	}
	for _, sq := range board.Pieces.White.Knights {
		whiteScore += PieceCodeToValue[WhiteKnight]
		whiteScore += knightPST[sq]
	}
	for _, sq := range board.Pieces.Black.Knights {
		blackScore += PieceCodeToValue[BlackKnight]
		blackScore += knightPST[mirrorIndex(sq)]
	}
	for _, sq := range board.Pieces.White.Bishops {
		whiteScore += PieceCodeToValue[WhiteBishop]
		whiteScore += bishopPST[sq]
	}
	for _, sq := range board.Pieces.Black.Bishops {
		blackScore += PieceCodeToValue[BlackBishop]
		blackScore += bishopPST[mirrorIndex(sq)]
	}
	for _, sq := range board.Pieces.White.Rooks {
		whiteScore += PieceCodeToValue[WhiteRook]
		whiteScore += rookPST[sq]
	}
	for _, sq := range board.Pieces.Black.Rooks {
		blackScore += PieceCodeToValue[BlackRook]
		blackScore += rookPST[mirrorIndex(sq)]
	}
	for _, sq := range board.Pieces.White.Queens {
		whiteScore += PieceCodeToValue[WhiteQueen]
		whiteScore += queenPST[sq]
	}
	for _, sq := range board.Pieces.Black.Queens {
		blackScore += PieceCodeToValue[BlackQueen]
		blackScore += queenPST[mirrorIndex(sq)]
	}
	if board.Pieces.White.King != 0 {
		whiteScore += PieceCodeToValue[WhiteKing]
		whiteScore += kingPST[board.Pieces.White.King]
	}
	if board.Pieces.Black.King != 0 {
		blackScore += PieceCodeToValue[BlackKing]
		blackScore += kingPST[mirrorIndex(board.Pieces.Black.King)]
	}

	// Encourage mating the king in the endgame
	material := 0
	material += len(board.Pieces.White.Pawns) + len(board.Pieces.Black.Pawns)
	material += len(board.Pieces.White.Knights)*3 + len(board.Pieces.Black.Knights)*3
	material += len(board.Pieces.White.Bishops)*3 + len(board.Pieces.Black.Bishops)*3
	material += len(board.Pieces.White.Rooks)*5 + len(board.Pieces.Black.Rooks)*5
	material += len(board.Pieces.White.Queens)*9 + len(board.Pieces.Black.Queens)*9
	if material <= 14 && board.Pieces.White.King != 0 && board.Pieces.Black.King != 0 {
		wKing := board.Pieces.White.King
		bKing := board.Pieces.Black.King
		wRank := int(wKing / 8)
		wFile := int(wKing % 8)
		bRank := int(bKing / 8)
		bFile := int(bKing % 8)
		dist := abs(wRank-bRank) + abs(wFile-bFile)
		// If one side has more material, encourage reducing the distance between kings
		wMat := len(board.Pieces.White.Queens)*9 + len(board.Pieces.White.Rooks)*5 + len(board.Pieces.White.Bishops)*3 + len(board.Pieces.White.Knights)*3 + len(board.Pieces.White.Pawns)
		bMat := len(board.Pieces.Black.Queens)*9 + len(board.Pieces.Black.Rooks)*5 + len(board.Pieces.Black.Bishops)*3 + len(board.Pieces.Black.Knights)*3 + len(board.Pieces.Black.Pawns)
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
	if board.IsKingInCheck(board.WhiteToMove) {
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
