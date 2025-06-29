package libra

import (
	"math/bits"
)

// AddQuietOrCapture adds a quiet move if the destination is empty, or a capture if occupied by an opponent's piece.
// Returns the new slice and true if a quiet move was added, false if a capture or blocked.
func (board *Board) AddQuietOrCapture(piece, from, to byte, whiteToMove bool, moves []Move) ([]Move, bool) {
	if board.IsSquareEmpty(to) {
		moves = board.AddQuietMove(piece, from, to, moves)
		return moves, true
	}

	if board.IsSquareKing(to) {
		return moves, false
	}

	if (whiteToMove && board.IsPieceAtSquareBlack(to)) || (!whiteToMove && board.IsPieceAtSquareWhite(to)) {
		moves = board.AddCapture(piece, from, to, MoveCapture, whiteToMove, moves)
		return moves, false
	}

	return moves, false
}

// AddMove appends a move to the move list and returns the new slice.
func (board *Board) AddMove(move Move, moves []Move) []Move {
	return append(moves, move)
}

// AddQuietMove adds a non-capturing move to the move list and returns the new slice.
func (board *Board) AddQuietMove(piece, from, to byte, moves []Move) []Move {
	move := NewMove(piece, from, to, MoveQuiet, 0, 0)
	return append(moves, move)
}

// AddCastleMove adds a castling move to the move list and returns the new slice.
func (board *Board) AddCastleMove(piece, from, to byte, moves []Move) []Move {
	move := NewMove(piece, from, to, MoveCastle, 0, 0)
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
func (board *Board) AddCapture(piece, from, to, moveType byte, whiteToMove bool, moves []Move) []Move {
	captured := board.getCapturedPiece(moveType, to, whiteToMove)
	move := NewMove(piece, from, to, moveType, 0, captured)
	return append(moves, move)
}

// AddPromotion adds all possible promotion moves (to queen, rook, bishop, knight) for a pawn reaching the last rank.
// If captured != 0, adds promotion-capture moves. Returns the new slice.
func (board *Board) AddPromotion(piece, from, to, captured byte, whiteToMove bool, moves []Move) []Move {

	promotionPieces := []byte{WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight}
	if !whiteToMove {
		promotionPieces = []byte{BlackQueen, BlackRook, BlackBishop, BlackKnight}
	}
	for _, promo := range promotionPieces {
		moveType := MovePromotion
		if captured != 0 {
			moveType = MovePromotionCapture
		}
		move := NewMove(piece, from, to, byte(moveType), promo, captured)
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
	var piece byte
	if whiteToMove {
		pawns = board.WhitePawns
		dir = -8
		startRank = 6
		promotionRank = 0
		piece = WhitePawn
	} else {
		pawns = board.BlackPawns
		dir = 8
		startRank = 1
		promotionRank = 7
		piece = BlackPawn
	}
	for bb := pawns; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		file := board.SquareToFile(square)
		rank := board.SquareToRank(square)
		to := int8(square) + dir
		if to >= 0 && to < 64 && !board.IsSquareOccupied(byte(to)) {
			if byte(to/8) == promotionRank {
				moves = board.AddPromotion(piece, square, byte(to), 0, whiteToMove, moves)
			} else {
				moves = board.AddQuietMove(piece, square, byte(to), moves)
				if rank == startRank {
					twoForward := int8(square) + 2*dir
					if twoForward >= 0 && twoForward < 64 && !board.IsSquareOccupied(byte(twoForward)) {
						moves = board.AddQuietMove(piece, square, byte(twoForward), moves)
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
					capture := board.PieceAtSquare(byte(captureTo))
					moves = board.AddPromotion(piece, square, byte(captureTo), capture, whiteToMove, moves)
				} else {
					moves = board.AddCapture(piece, square, byte(captureTo), MoveCapture, whiteToMove, moves)
				}
			}
			if board.IsSquareOnPassant(byte(captureTo)) {
				if (whiteToMove && rank == 3) || (!whiteToMove && rank == 4) {
					moves = board.AddCapture(piece, square, byte(captureTo), MoveEnPassant, whiteToMove, moves)
				}
			}
		}
		bb &= bb - 1
	}
	return moves
}

// GenerateSlidingMoves generates all moves for sliding pieces (rooks, bishops, queens) using precomputed rays.
func (board *Board) GenerateSlidingMoves(piece byte, bitboard uint64, startDir byte, endDir byte, whiteToMove bool) []Move {
	moves := []Move{}
	for bb := bitboard; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			var ray uint64
			if dirOffset < 4 { // Rook directions
				ray = RookRays[square][dirOffset]
			} else { // Bishop directions
				ray = BishopRays[square][dirOffset-4]
			}

			for s := int(square) + int(offset); s >= 0 && s < 64 && (ray&(uint64(1)<<s)) != 0; s += int(offset) {
				var isQuietMove bool
				moves, isQuietMove = board.AddQuietOrCapture(piece, square, byte(s), whiteToMove, moves)
				if !isQuietMove {
					break // Stop after a capture or own piece
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
	var piece byte
	if whiteToMove {
		kingSq = byte(bits.TrailingZeros64(board.WhiteKing))
		piece = WhiteKing
	} else {
		kingSq = byte(bits.TrailingZeros64(board.BlackKing))
		piece = BlackKing
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[kingSq][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(kingSq) + offset
			moves, _ = board.AddQuietOrCapture(piece, kingSq, byte(squareTo), whiteToMove, moves)
		}
	}
	return moves
}

// Generates all castling moves
// King cannot castle out of, through, or into check;
// squares between must be empty.
func (board *Board) GenerateCastleMoves(whiteToMove bool) []Move {
	moves := []Move{}
	var piece byte
	if whiteToMove {
		piece = WhiteKing
		if board.Castling.WhiteQueenSide &&
			board.IsSquareWhiteKing(SquareE1) &&
			board.IsSquareWhiteRook(SquareA1) &&
			board.IsSquareEmpty(SquareB1) &&
			board.IsSquareEmpty(SquareC1) &&
			board.IsSquareEmpty(SquareD1) &&
			!board.IsSquareAttacked(SquareC1, whiteToMove) &&
			!board.IsSquareAttacked(SquareD1, whiteToMove) &&
			!board.IsSquareAttacked(SquareE1, whiteToMove) {
			moves = board.AddCastleMove(piece, SquareE1, SquareC1, moves)
		}

		if board.Castling.WhiteKingSide &&
			board.IsSquareWhiteKing(SquareE1) &&
			board.IsSquareWhiteRook(SquareH1) &&
			board.IsSquareEmpty(SquareF1) &&
			board.IsSquareEmpty(SquareG1) &&
			!board.IsSquareAttacked(SquareF1, whiteToMove) &&
			!board.IsSquareAttacked(SquareG1, whiteToMove) &&
			!board.IsSquareAttacked(SquareE1, whiteToMove) {
			moves = board.AddCastleMove(piece, SquareE1, SquareG1, moves)
		}
	} else {
		piece = BlackKing
		if board.Castling.BlackQueenSide &&
			board.IsSquareBlackKing(SquareE8) &&
			board.IsSquareBlackRook(SquareA8) &&
			board.IsSquareEmpty(SquareB8) &&
			board.IsSquareEmpty(SquareC8) &&
			board.IsSquareEmpty(SquareD8) &&
			!board.IsSquareAttacked(SquareC8, whiteToMove) &&
			!board.IsSquareAttacked(SquareD8, whiteToMove) &&
			!board.IsSquareAttacked(SquareE8, whiteToMove) {
			moves = board.AddCastleMove(piece, SquareE8, SquareC8, moves)
		}

		if board.Castling.BlackKingSide &&
			board.IsSquareBlackKing(SquareE8) &&
			board.IsSquareBlackRook(SquareH8) &&
			board.IsSquareEmpty(SquareF8) &&
			board.IsSquareEmpty(SquareG8) &&
			!board.IsSquareAttacked(SquareF8, whiteToMove) &&
			!board.IsSquareAttacked(SquareG8, whiteToMove) &&
			!board.IsSquareAttacked(SquareE8, whiteToMove) {
			moves = board.AddCastleMove(piece, SquareE8, SquareG8, moves)
		}
	}
	return moves
}

// GenerateRookMoves generates all rook moves for the current side.
func (board *Board) GenerateRookMoves(whiteToMove bool) []Move {
	var rooks uint64
	var piece byte
	if whiteToMove {
		piece = WhiteRook
		rooks = board.WhiteRooks
	} else {
		piece = BlackRook
		rooks = board.BlackRooks
	}
	return board.GenerateSlidingMoves(piece, rooks, 0, 4, whiteToMove)
}

// GenerateBishopMoves generates all bishop moves for the current side.
func (board *Board) GenerateBishopMoves(whiteToMove bool) []Move {
	var bishops uint64
	var piece byte
	if whiteToMove {
		piece = WhiteBishop
		bishops = board.WhiteBishops
	} else {
		piece = BlackBishop
		bishops = board.BlackBishops
	}
	return board.GenerateSlidingMoves(piece, bishops, 4, 8, whiteToMove)
}

// GenerateQueenMoves generates all queen moves for the current side.
func (board *Board) GenerateQueenMoves(whiteToMove bool) []Move {
	var queens uint64
	var piece byte
	if whiteToMove {
		piece = WhiteQueen
		queens = board.WhiteQueens
	} else {
		piece = BlackQueen
		queens = board.BlackQueens
	}
	return board.GenerateSlidingMoves(piece, queens, 0, 8, whiteToMove)
}

// GenerateKnightMoves generates all knight moves for the current side.
func (board *Board) GenerateKnightMoves(whiteToMove bool) []Move {
	moves := []Move{}
	var knights uint64
	var piece byte
	if whiteToMove {
		knights = board.WhiteKnights
		piece = WhiteKnight
	} else {
		knights = board.BlackKnights
		piece = BlackKnight
	}
	for bb := knights; bb != 0; {
		square := byte(bits.TrailingZeros64(bb))
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := KnightOffsets[square][moveIndex]
			if squareTo < 255 {
				moves, _ = board.AddQuietOrCapture(piece, square, squareTo, whiteToMove, moves)
			}
		}
		bb &= bb - 1
	}
	return moves
}

func (board *Board) GeneratePseudoLegalMoves() []Move {
	moves := []Move{}
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
	return legalMoves
}

// Checks if the move leaves the king in check and undoes the move.
func (board *Board) IsMoveLegal(move Move) bool {
	prev := board.Move(move)
	kingSq := board.PassiveKingSquare()

	inCheck := board.IsSquareAttacked(kingSq, !board.WhiteToMove)
	board.UndoMove(prev)
	return !inCheck
}

func (board *Board) IsSquareAttacked(square byte, whiteToMove bool) bool {
	if board.IsSquareAttackedByPawns(square, whiteToMove) ||
		board.IsSquareAttackedByKnights(square, whiteToMove) ||
		board.IsSquareAttackedByKing(square, whiteToMove) ||
		board.IsSquareAttackedBySlidingPieces(square, whiteToMove) {
		return true
	}
	return false
}

// IsSquareAttackedByPawns returns true if the king of the given color is attacked by any enemy pawn.
func (board *Board) IsSquareAttackedByPawns(square byte, whiteToMove bool) bool {
	squareBB := uint64(1) << square
	var pawnAttackersBB uint64

	if whiteToMove {
		// Black pawns attack
		if square%8 != 0 {
			pawnAttackersBB |= squareBB >> 9
		}
		if square%8 != 7 {
			pawnAttackersBB |= squareBB >> 7
		}
		return (pawnAttackersBB & board.BlackPawns) != 0
	} else {
		// White pawns attack
		if square%8 != 0 {
			pawnAttackersBB |= squareBB << 7
		}
		if square%8 != 7 {
			pawnAttackersBB |= squareBB << 9
		}
		return (pawnAttackersBB & board.WhitePawns) != 0
	}
}

// IsSquareAttackedByKing returns true if the square is attacked by the king
func (board *Board) IsSquareAttackedByKing(square byte, whiteToMove bool) bool {
	var enemyKingBB uint64
	if whiteToMove {
		enemyKingBB = board.BlackKing
	} else {
		enemyKingBB = board.WhiteKing
	}

	for i := 0; i < 8; i++ {
		sq := KingOffsets[square][i]
		if sq < 64 && ((enemyKingBB & (uint64(1) << sq)) != 0) {
			return true
		}
	}
	return false
}

// IsSquareAttackedByKnights returns true if the square is attacked by the knights
func (board *Board) IsSquareAttackedByKnights(square byte, whiteToMove bool) bool {
	var enemyKnights uint64
	if whiteToMove {
		enemyKnights = board.BlackKnights
	} else {
		enemyKnights = board.WhiteKnights
	}

	for i := 0; i < 8; i++ {
		sq := KnightOffsets[square][i]
		if sq < 64 && ((enemyKnights & (uint64(1) << sq)) != 0) {
			return true
		}
	}
	return false
}

// Optimized sliding piece attack detection using precomputed rays (corrected)
func (board *Board) IsSquareAttackedBySlidingPieces(square byte, whiteToMove bool) bool {
	var rooksAndQueens, bishopsAndQueens uint64
	if whiteToMove {
		rooksAndQueens = board.BlackRooks | board.BlackQueens
		bishopsAndQueens = board.BlackBishops | board.BlackQueens
	} else {
		rooksAndQueens = board.WhiteRooks | board.WhiteQueens
		bishopsAndQueens = board.WhiteBishops | board.WhiteQueens
	}
	occupied := board.OccupiedSquares()
	// Rook/Queen directions
	for dir := 0; dir < 4; dir++ {
		ray := RookRays[square][dir]
		attackers := ray & rooksAndQueens
		if attackers == 0 {
			continue
		}
		// Find the closest attacker in this direction
		var sqStep int
		switch dir {
		case 0:
			sqStep = -8 // North
		case 1:
			sqStep = 1 // East
		case 2:
			sqStep = 8 // South
		case 3:
			sqStep = -1 // West
		}
		for s := int(square) + sqStep; s >= 0 && s < 64 && (ray&(1<<s)) != 0; s += sqStep {
			mask := uint64(1) << s
			if (occupied & mask) != 0 {
				if (rooksAndQueens & mask) != 0 {
					return true
				}
				break
			}
		}
	}
	// Bishop/Queen directions
	for dir := 0; dir < 4; dir++ {
		ray := BishopRays[square][dir]
		attackers := ray & bishopsAndQueens
		if attackers == 0 {
			continue
		}
		// Find the closest attacker in this direction
		var sqStep int
		switch dir {
		case 0:
			sqStep = -7 // NE
		case 1:
			sqStep = 9 // SE
		case 2:
			sqStep = 7 // SW
		case 3:
			sqStep = -9 // NW
		}
		for s := int(square) + sqStep; s >= 0 && s < 64 && (ray&(1<<s)) != 0; s += sqStep {
			mask := uint64(1) << s
			if (occupied & mask) != 0 {
				if (bishopsAndQueens & mask) != 0 {
					return true
				}
				break
			}
		}
	}
	return false
}
