package libra

import (
	"math/bits"
	"sort"
)

// AddQuietOrCapture adds a quiet move if the destination is empty, or a capture if occupied by an opponent's piece.
// Returns the new slice and true if a quiet move was added, false if a capture or blocked.
func (board *Board) AddQuietOrCapture(from, to byte, whiteToMove bool, moves []Move) ([]Move, bool) {
	if board.IsSquareEmpty(to) {
		moves = board.AddQuietMove(from, to, moves)
		return moves, true
	}

	if board.IsSquareKing(to) {
		return moves, false
	}

	if (whiteToMove && board.IsPieceAtSquareBlack(to)) || (!whiteToMove && board.IsPieceAtSquareWhite(to)) {
		moves = board.AddCapture(from, to, MoveCapture, whiteToMove, moves)
		return moves, false
	}

	return moves, false
}

// AddMove appends a move to the move list and returns the new slice.
func (board *Board) AddMove(move Move, moves []Move) []Move {
	return append(moves, move)
}

// AddQuietMove adds a non-capturing move to the move list and returns the new slice.
func (board *Board) AddQuietMove(from, to byte, moves []Move) []Move {
	move := NewMove(from, to, MoveQuiet, [2]byte{0, 0})
	return append(moves, move)
}

// AddCastleMove adds a castling move to the move list and returns the new slice.
func (board *Board) AddCastleMove(from, to byte, moves []Move) []Move {
	move := NewMove(from, to, MoveCastle, [2]byte{0, 0})
	return append(moves, move)
}

// getCapturedPiece returns the captured piece for a given move, handling en passant correctly.
func (board *Board) getCapturedPiece(moveType byte, to byte, whiteToMove bool) byte {
	if moveType == MoveEnPassant {
		if whiteToMove {
			return BlackPawn
		} else {
			return WhitePawn
		}
	}
	return board.PieceAtSquare(to)
}

// AddCapture adds a capturing move to the move list. Handles en passant as a special case. Returns the new slice.
func (board *Board) AddCapture(from, to, moveType byte, whiteToMove bool, moves []Move) []Move {
	captured := board.getCapturedPiece(moveType, to, whiteToMove)
	move := NewMove(from, to, moveType, [2]byte{captured, 0})
	return append(moves, move)
}

// AddPromotion adds all possible promotion moves (to queen, rook, bishop, knight) for a pawn reaching the last rank.
// If captured != 0, adds promotion-capture moves. Returns the new slice.
func (board *Board) AddPromotion(from, to, captured byte, whiteToMove bool, moves []Move) []Move {

	promotionPieces := []byte{WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight}
	if !whiteToMove {
		promotionPieces = []byte{BlackQueen, BlackRook, BlackBishop, BlackKnight}
	}
	for _, promo := range promotionPieces {
		moveType := MovePromotion
		if captured != 0 {
			moveType = MovePromotionCapture
		}
		move := NewMove(from, to, byte(moveType), [2]byte{promo, captured})
		moves = append(moves, move)
	}
	return moves
}

// GeneratePawnMoves generates all pawn moves (including promotions, captures, en passant) for the current side.
func (board *Board) GeneratePawnMoves(whiteToMove bool) []Move {
	moves := []Move{}
	var pawns uint64
	var dir int8
	var startRank, promotionRank byte
	if whiteToMove {
		pawns = board.WhitePawns
		dir = -8
		startRank = 6
		promotionRank = 0
	} else {
		pawns = board.BlackPawns
		dir = 8
		startRank = 1
		promotionRank = 7
	}
	for bb := pawns; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		file := board.SquareToFile(square)
		rank := board.SquareToRank(square)
		to := int8(square) + dir
		if to >= 0 && to < 64 && !board.IsSquareOccupied(byte(to)) {
			if byte(to/8) == promotionRank {
				moves = board.AddPromotion(square, byte(to), 0, whiteToMove, moves)
			} else {
				moves = board.AddQuietMove(square, byte(to), moves)
				if rank == startRank {
					twoForward := int8(square) + 2*dir
					if twoForward >= 0 && twoForward < 64 && !board.IsSquareOccupied(byte(twoForward)) {
						moves = board.AddQuietMove(square, byte(twoForward), moves)
					}
				}
			}
		}
		for _, df := range []int8{-1, 1} {
			captureFile := int8(file) + df
			if captureFile < 0 || captureFile > 7 {
				continue
			}
			captureTo := int8(square) + dir + df
			if captureTo < 0 || captureTo >= 64 {
				continue
			}
			if board.IsSquareOccupied(byte(captureTo)) && board.IsPieceAtSquareWhite(byte(captureTo)) != whiteToMove {
				if byte(captureTo/8) == promotionRank {
					moves = board.AddPromotion(square, byte(captureTo), board.PieceAtSquare(byte(captureTo)), whiteToMove, moves)
				} else {
					moves = board.AddCapture(square, byte(captureTo), MoveCapture, whiteToMove, moves)
				}
			}
			if board.IsSquareOnPassant(byte(captureTo)) {
				if (whiteToMove && rank == 3) || (!whiteToMove && rank == 4) {
					moves = board.AddCapture(square, byte(captureTo), MoveEnPassant, whiteToMove, moves)
				}
			}
		}
		bb &= bb - 1
	}
	return moves
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

// GenerateSlidingMoves generates all moves for sliding pieces (rooks, bishops, queens) in the given directions.
func (board *Board) GenerateSlidingMoves(bitboard uint64, startDir byte, endDir byte, whiteToMove bool) []Move {
	moves := []Move{}
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
				var isQuietMove bool
				moves, isQuietMove = board.AddQuietOrCapture(square, byte(squareTo), whiteToMove, moves)
				if !isQuietMove {
					break
				}
			}
		}
		bb &= bb - 1
	}
	return moves
}

// GenerateKingMoves generates all king moves (excluding castling) for the current side.
func (board *Board) GenerateKingMoves(whiteToMove bool) []Move {
	moves := []Move{}
	var kingSq byte
	if whiteToMove {
		if board.WhiteKing == 0 {
			return moves
		}
		kingSq = byte(bits.TrailingZeros64(board.WhiteKing))
	} else {
		if board.BlackKing == 0 {
			return moves
		}
		kingSq = byte(bits.TrailingZeros64(board.BlackKing))
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[kingSq][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(kingSq) + offset
			moves, _ = board.AddQuietOrCapture(kingSq, byte(squareTo), whiteToMove, moves)
		}
	}
	return moves
}

// GenerateCastleMoves generates castling moves if the king and rook have not moved and the path is clear and not attacked.
// Chess rules: King cannot castle out of, through, or into check; squares between must be empty.
func (board *Board) GenerateCastleMoves(whiteToMove bool) []Move {
	moves := []Move{}
	if whiteToMove {
		if board.CastlingAvailability.WhiteQueenSide &&
			board.IsSquareWhiteKing(SquareE1) &&
			board.IsSquareWhiteRook(SquareA1) &&
			board.IsSquareEmpty(SquareB1) &&
			board.IsSquareEmptyAndNotAttacked(SquareC1) &&
			board.IsSquareEmptyAndNotAttacked(SquareD1) &&
			!board.IsSquareAttacked(SquareE1) {
			moves = board.AddCastleMove(SquareE1, SquareC1, moves)
		}

		if board.CastlingAvailability.WhiteKingSide &&
			board.IsSquareWhiteKing(SquareE1) &&
			board.IsSquareWhiteRook(SquareH1) &&
			board.IsSquareEmptyAndNotAttacked(SquareF1) &&
			board.IsSquareEmptyAndNotAttacked(SquareG1) &&
			!board.IsSquareAttacked(SquareE1) {
			moves = board.AddCastleMove(SquareE1, SquareG1, moves)
		}
	} else {
		if board.CastlingAvailability.BlackQueenSide &&
			board.IsSquareBlackKing(SquareE8) &&
			board.IsSquareBlackRook(SquareA8) &&
			board.IsSquareEmpty(SquareB8) &&
			board.IsSquareEmptyAndNotAttacked(SquareC8) &&
			board.IsSquareEmptyAndNotAttacked(SquareD8) &&
			!board.IsSquareAttacked(SquareE8) {
			moves = board.AddCastleMove(SquareE8, SquareC8, moves)
		}

		if board.CastlingAvailability.BlackKingSide &&
			board.IsSquareBlackKing(SquareE8) &&
			board.IsSquareBlackRook(SquareH8) &&
			board.IsSquareEmptyAndNotAttacked(SquareF8) &&
			board.IsSquareEmptyAndNotAttacked(SquareG8) &&
			!board.IsSquareAttacked(SquareE8) {
			moves = board.AddCastleMove(SquareE8, SquareG8, moves)
		}
	}
	return moves
}

// GenerateRookMoves generates all rook moves for the current side.
func (board *Board) GenerateRookMoves(whiteToMove bool) []Move {
	var rooks uint64
	if whiteToMove {
		rooks = board.WhiteRooks
	} else {
		rooks = board.BlackRooks
	}
	return board.GenerateSlidingMoves(rooks, 0, 4, whiteToMove)
}

// GenerateBishopMoves generates all bishop moves for the current side.
func (board *Board) GenerateBishopMoves(whiteToMove bool) []Move {
	var bishops uint64
	if whiteToMove {
		bishops = board.WhiteBishops
	} else {
		bishops = board.BlackBishops
	}
	return board.GenerateSlidingMoves(bishops, 4, 8, whiteToMove)
}

// GenerateQueenMoves generates all queen moves for the current side.
func (board *Board) GenerateQueenMoves(whiteToMove bool) []Move {
	var queens uint64
	if whiteToMove {
		queens = board.WhiteQueens
	} else {
		queens = board.BlackQueens
	}
	return board.GenerateSlidingMoves(queens, 0, 8, whiteToMove)
}

// GenerateKnightMoves generates all knight moves for the current side.
func (board *Board) GenerateKnightMoves(whiteToMove bool) []Move {
	moves := []Move{}
	var knights uint64
	if whiteToMove {
		knights = board.WhiteKnights
	} else {
		knights = board.BlackKnights
	}
	for bb := knights; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				moves, _ = board.AddQuietOrCapture(square, squareTo, whiteToMove, moves)
			}
		}
		bb &= bb - 1
	}
	return moves
}

func (board *Board) GeneratePseudoLegalMoves() []Move {
	moves := []Move{}
	board.GenerateAttackedSquares(!board.WhiteToMove)
	moves = append(moves, board.GeneratePawnMoves(board.WhiteToMove)...)
	moves = append(moves, board.GenerateKnightMoves(board.WhiteToMove)...)
	moves = append(moves, board.GenerateBishopMoves(board.WhiteToMove)...)
	moves = append(moves, board.GenerateRookMoves(board.WhiteToMove)...)
	moves = append(moves, board.GenerateQueenMoves(board.WhiteToMove)...)
	moves = append(moves, board.GenerateKingMoves(board.WhiteToMove)...)
	moves = append(moves, board.GenerateCastleMoves(board.WhiteToMove)...)
	return moves
}

func (board *Board) GenerateLegalMoves() []Move {
	legalMoves := []Move{}
	moves := board.GeneratePseudoLegalMoves()
	for _, move := range moves {
		if board.IsMoveLegal(move) {
			legalMoves = append(legalMoves, move)
		}
	}
	// Sort moves by MoveType descending, then by From, To, and promotion piece for full determinism
	sort.SliceStable(legalMoves, func(i, j int) bool {
		moveA := legalMoves[i]
		moveB := legalMoves[j]
		if moveA.MoveType != moveB.MoveType {
			// Sort by MoveType descending: captures and promotions first, then quiet moves
			return moveA.MoveType > moveB.MoveType
		}
		if moveA.From != moveB.From {
			// Sort by From square ascending
			return moveA.From < moveB.From
		}
		if moveA.To != moveB.To {
			// Sort by To square ascending
			return moveA.To < moveB.To
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
		// Sort by capture value if both moves are captures
		// This ensures that if two captures are available, the one with the higher value piece captured is preferred.
		if moveA.MoveType == MoveCapture && moveB.MoveType == MoveCapture {
			attackerA := board.PieceAtSquare(moveA.From)
			captureA := moveA.Data[0]
			attackerB := board.PieceAtSquare(moveB.From)
			captureB := moveB.Data[0]
			valueA := PieceCodeToValue[captureA] - PieceCodeToValue[attackerA]
			valueB := PieceCodeToValue[captureB] - PieceCodeToValue[attackerB]
			if valueA != valueB {
				return valueA > valueB
			}
		}
		return false
	})
	return legalMoves
}

func (board *Board) IsMoveLegal(move Move) bool {
	prev := board.Move(move)
	// Generate attacked squares for the opponent after the move
	inCheck := board.IsKingInCheck(!board.WhiteToMove)
	board.UndoMove(prev)
	return !inCheck
}

func (board *Board) IsKingInCheck(whiteToMove bool) bool {
	board.GenerateAttackedSquares(!whiteToMove)
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

func (board *Board) ResetAttackedSquares() {
	board.AttackedSquares = 0
}

func (board *Board) GenerateAttackedSquares(whiteToMove bool) {
	board.ResetAttackedSquares()

	if whiteToMove {
		board.MarkSlidingAttacks(board.WhiteQueens, 0, 8)
		board.MarkSlidingAttacks(board.WhiteBishops, 4, 8)
		board.MarkSlidingAttacks(board.WhiteRooks, 0, 4)
	} else {
		board.MarkSlidingAttacks(board.BlackQueens, 0, 8)
		board.MarkSlidingAttacks(board.BlackBishops, 4, 8)
		board.MarkSlidingAttacks(board.BlackRooks, 0, 4)
	}

	var knights uint64
	if whiteToMove {
		knights = board.WhiteKnights
	} else {
		knights = board.BlackKnights
	}
	for bb := knights; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				board.AttackedSquares |= (uint64(1) << squareTo)
			}
		}
		bb &= bb - 1
	}

	var kingSq byte
	if whiteToMove {
		if board.WhiteKing == 0 {
			return
		}
		kingSq = byte(bits.TrailingZeros64(board.WhiteKing))
	} else {
		if board.BlackKing == 0 {
			return
		}
		kingSq = byte(bits.TrailingZeros64(board.BlackKing))
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[kingSq][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(kingSq) + offset
			if squareTo >= 0 && squareTo < 64 {
				board.AttackedSquares |= (uint64(1) << byte(squareTo))
			}
		}
	}

	var pawns uint64
	var dir int8
	if whiteToMove {
		pawns = board.WhitePawns
		dir = -8
	} else {
		pawns = board.BlackPawns
		dir = 8
	}
	for bb := pawns; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		file := board.SquareToFile(square)
		for _, df := range []int8{-1, 1} {
			attackFile := int8(file) + df
			if attackFile < 0 || attackFile > 7 {
				continue
			}
			attack := int8(square) + dir + df
			if attack >= 0 && attack < 64 {
				board.AttackedSquares |= (uint64(1) << byte(attack))
			}
		}
		bb &= bb - 1
	}
}
