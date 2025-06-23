package libra

/*
import "math/bits"

func (board *Board) MarkPawnAttacks(whiteToMove bool) {
	var pawns uint64
	if whiteToMove {
		pawns = board.WhitePawns
	} else {
		pawns = board.BlackPawns
	}
	for bb := pawns; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		file := board.SquareToFile(square)
		rank := board.SquareToRank(square)

		if whiteToMove {
			if rank > 0 {
				if file > 0 {
					board.AttackedSquares |= (uint64(1) << (square - 9)) // Capture left
				}
				if file < 7 {
					board.AttackedSquares |= (uint64(1) << (square - 7)) // Capture right
				}
			}
		} else {
			if rank < 7 {
				if file > 0 {
					board.AttackedSquares |= (uint64(1) << (square + 7)) // Capture left
				}
				if file < 7 {
					board.AttackedSquares |= (uint64(1) << (square + 9)) // Capture right
				}
			}
		}
		bb &= bb - 1
	}
}

func (board *Board) MarkKnightAttacks(whiteToMove bool) {
	var knights uint64
	if whiteToMove {
		knights = board.WhiteKnights
	} else {
		knights = board.BlackKnights
	}
	for bb := knights; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightOffsets[square][moveIndex]
			if squareTo < 255 {
				board.AttackedSquares |= (uint64(1) << squareTo)
			}
		}
		bb &= bb - 1
	}
}

func (board *Board) MarkKingAttacks(whiteToMove bool) {
	var kings uint64
	if whiteToMove {
		kings = board.WhiteKing
	} else {
		kings = board.BlackKing
	}
	for bb := kings; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKingOffsets[square][moveIndex]
			if squareTo < 255 {
				board.AttackedSquares |= (uint64(1) << squareTo)
			}
		}
		bb &= bb - 1
	}
}

// MarkSlidingAttacks marks all squares attacked by sliding pieces (rooks, bishops, queens) in the given directions.
// Used for attack maps and move generation.
func (board *Board) MarkSlidingAttacks(bitboard uint64, startDir byte, endDir byte) {
	for bb := bitboard; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			amountToMove := int8(SquaresToEdge[square][dirOffset])
			for moveIndex := int8(1); moveIndex <= amountToMove; moveIndex++ {
				squareTo := int8(square) + offset*moveIndex
				if squareTo < 0 || squareTo >= 64 {
					break
				}
				board.AttackedSquares |= (uint64(1) << byte(squareTo))
				if board.IsSquareOccupied(byte(squareTo)) {
					break
				}
			}
		}
		bb &= bb - 1
	}
}

func (board *Board) MarkAttackedSquares(whiteToMove bool) {
	board.AttackedSquares = 0

	board.MarkKingAttacks(whiteToMove)
	board.MarkPawnAttacks(whiteToMove)
	board.MarkKnightAttacks(whiteToMove)

	if whiteToMove {
		board.MarkSlidingAttacks(board.WhiteQueens, 0, 8)
		board.MarkSlidingAttacks(board.WhiteBishops, 4, 8)
		board.MarkSlidingAttacks(board.WhiteRooks, 0, 4)
	} else {
		board.MarkSlidingAttacks(board.BlackQueens, 0, 8)
		board.MarkSlidingAttacks(board.BlackBishops, 4, 8)
		board.MarkSlidingAttacks(board.BlackRooks, 0, 4)
	}
}

func (board *Board) IsKingInCheckLegacy(whiteToMove bool) bool {
	board.MarkAttackedSquares(!whiteToMove)
	var kingSq byte
	if whiteToMove {
		if board.WhiteKing == 0 {
			return false
		}
		kingSq = byte(bits.TrailingZeros64(board.WhiteKing))
	} else {
		if board.BlackKing == 0 {
			return false
		}
		kingSq = byte(bits.TrailingZeros64(board.BlackKing))
	}
	return board.IsSquareAttacked(kingSq)
}
*/
