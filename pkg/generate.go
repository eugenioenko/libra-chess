package libra

import (
	"sort"
)

// AddQuietOrCapture adds a quiet move if the destination is empty, or a capture if occupied by an opponent's piece.
// Returns the new slice and true if a quiet move was added, false if a capture or blocked.
func (board *Board) AddQuietOrCapture(from, to byte, whiteToMove bool, moves []Move) ([]Move, bool) {
	if board.IsSquareEmpty(to) {
		moves = board.AddQuietMove(from, to, moves)
		return moves, true
	}

	if board.Position[to] == WhiteKing || board.Position[to] == BlackKing {
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
	return board.Position[to]
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
	var squares []byte
	var dir int8
	var startRank, promotionRank byte
	if whiteToMove {
		squares = board.Pieces.White.Pawns
		dir = -8
		startRank = 6
		promotionRank = 0
	} else {
		squares = board.Pieces.Black.Pawns
		dir = 8
		startRank = 1
		promotionRank = 7
	}
	for _, square := range squares {
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
					moves = board.AddPromotion(square, byte(captureTo), board.Position[byte(captureTo)], whiteToMove, moves)
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
	}
	return moves
}

// MarkSlidingAttacks marks all squares attacked by sliding pieces (rooks, bishops, queens) in the given directions.
// Used for attack maps and move generation.
func (board *Board) MarkSlidingAttacks(pieces []byte, startDir byte, endDir byte) {
	for _, square := range pieces {
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			amountToMove := int8(SquaresToEdge[square][dirOffset])
			for moveIndex := int8(1); moveIndex <= amountToMove; moveIndex++ {
				squareTo := int8(square) + offset*moveIndex
				if squareTo < 0 || squareTo >= 64 {
					break
				}
				board.AttackedSquares[byte(squareTo)] = true
				if board.IsSquareOccupied(byte(squareTo)) {
					break
				}
			}
		}
	}
}

// GenerateSlidingMoves generates all moves for sliding pieces (rooks, bishops, queens) in the given directions.
func (board *Board) GenerateSlidingMoves(pieces []byte, startDir byte, endDir byte, whiteToMove bool) []Move {
	moves := []Move{}
	for _, square := range pieces {
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
	}
	return moves
}

// GenerateKingMoves generates all king moves (excluding castling) for the current side.
func (board *Board) GenerateKingMoves(whiteToMove bool) []Move {
	moves := []Move{}
	square := board.Pieces.White.King
	if !whiteToMove {
		square = board.Pieces.Black.King
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[square][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(square) + offset
			moves, _ = board.AddQuietOrCapture(square, byte(squareTo), whiteToMove, moves)
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
			board.Position[SquareE1] == WhiteKing &&
			board.Position[SquareA1] == WhiteRook &&
			board.IsSquareEmpty(SquareB1) &&
			board.IsSquareEmptyAndNotAttacked(SquareC1) &&
			board.IsSquareEmptyAndNotAttacked(SquareD1) &&
			!board.AttackedSquares[SquareE1] {
			moves = board.AddCastleMove(SquareE1, SquareC1, moves)
		}

		if board.CastlingAvailability.WhiteKingSide &&
			board.Position[SquareE1] == WhiteKing &&
			board.Position[SquareH1] == WhiteRook &&
			board.IsSquareEmptyAndNotAttacked(SquareF1) &&
			board.IsSquareEmptyAndNotAttacked(SquareG1) &&
			!board.AttackedSquares[SquareE1] {
			moves = board.AddCastleMove(SquareE1, SquareG1, moves)
		}
	} else {
		if board.CastlingAvailability.BlackQueenSide &&
			board.Position[SquareE8] == BlackKing &&
			board.Position[SquareA8] == BlackRook &&
			board.IsSquareEmpty(SquareB8) &&
			board.IsSquareEmptyAndNotAttacked(SquareC8) &&
			board.IsSquareEmptyAndNotAttacked(SquareD8) &&
			!board.AttackedSquares[SquareE8] {
			moves = board.AddCastleMove(SquareE8, SquareC8, moves)
		}

		if board.CastlingAvailability.BlackKingSide &&
			board.Position[SquareE8] == BlackKing &&
			board.Position[SquareH8] == BlackRook &&
			board.IsSquareEmptyAndNotAttacked(SquareF8) &&
			board.IsSquareEmptyAndNotAttacked(SquareG8) &&
			!board.AttackedSquares[SquareE8] {
			moves = board.AddCastleMove(SquareE8, SquareG8, moves)
		}
	}
	return moves
}

// GenerateRookMoves generates all rook moves for the current side.
func (board *Board) GenerateRookMoves(whiteToMove bool) []Move {
	rooks := board.Pieces.White.Rooks
	if !whiteToMove {
		rooks = board.Pieces.Black.Rooks
	}
	return board.GenerateSlidingMoves(rooks, 0, 4, whiteToMove)
}

// GenerateBishopMoves generates all bishop moves for the current side.
func (board *Board) GenerateBishopMoves(whiteToMove bool) []Move {
	bishops := board.Pieces.White.Bishops
	if !whiteToMove {
		bishops = board.Pieces.Black.Bishops
	}
	return board.GenerateSlidingMoves(bishops, 4, 8, whiteToMove)
}

// GenerateQueenMoves generates all queen moves for the current side.
func (board *Board) GenerateQueenMoves(whiteToMove bool) []Move {
	queens := board.Pieces.White.Queens
	if !whiteToMove {
		queens = board.Pieces.Black.Queens
	}
	return board.GenerateSlidingMoves(queens, 0, 8, whiteToMove)
}

// GenerateKnightMoves generates all knight moves for the current side.
func (board *Board) GenerateKnightMoves(whiteToMove bool) []Move {
	moves := []Move{}
	knights := board.Pieces.White.Knights
	if !whiteToMove {
		knights = board.Pieces.Black.Knights
	}
	for _, square := range knights {
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				moves, _ = board.AddQuietOrCapture(square, squareTo, whiteToMove, moves)
			}
		}
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
		if legalMoves[i].MoveType != legalMoves[j].MoveType {
			// Promotion and captures are prioritized
			return legalMoves[i].MoveType > legalMoves[j].MoveType
		}
		if legalMoves[i].From != legalMoves[j].From {
			// Lower 'From' square first
			return legalMoves[i].From < legalMoves[j].From
		}
		if legalMoves[i].To != legalMoves[j].To {
			// Lower 'To' square first
			return legalMoves[i].To < legalMoves[j].To
		}
		// For promotions, ensure consistent order by promotion piece
		if legalMoves[i].MoveType == MovePromotion || legalMoves[i].MoveType == MovePromotionCapture {
			if legalMoves[i].Data[0] != legalMoves[j].Data[0] {
				// For promotions, sort by piece value in ascending order: Knight < Bishop < Rook < Queen.
				// This ensures deterministic move ordering, so that when multiple promotions have equal evaluation,
				// the queen promotion (highest value) is preferred if all else is equal.
				return legalMoves[i].Data[0] < legalMoves[j].Data[0]
			}
		}
		// Moves are considered equal for sorting if all criteria match
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
	king := board.Pieces.White.King
	if !whiteToMove {
		king = board.Pieces.Black.King
	}
	return board.AttackedSquares[king]
}

func (board *Board) ResetAttackedSquares() {
	for i := range board.AttackedSquares {
		board.AttackedSquares[i] = false
	}
}

func (board *Board) GenerateAttackedSquares(whiteToMove bool) {
	board.ResetAttackedSquares()

	if whiteToMove {
		board.MarkSlidingAttacks(board.Pieces.White.Queens, 0, 8)
		board.MarkSlidingAttacks(board.Pieces.White.Bishops, 4, 8)
		board.MarkSlidingAttacks(board.Pieces.White.Rooks, 0, 4)
	} else {
		board.MarkSlidingAttacks(board.Pieces.Black.Queens, 0, 8)
		board.MarkSlidingAttacks(board.Pieces.Black.Bishops, 4, 8)
		board.MarkSlidingAttacks(board.Pieces.Black.Rooks, 0, 4)
	}

	knights := board.Pieces.White.Knights
	if !whiteToMove {
		knights = board.Pieces.Black.Knights
	}
	for _, square := range knights {
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				board.AttackedSquares[squareTo] = true
			}
		}
	}

	king := board.Pieces.White.King
	if !whiteToMove {
		king = board.Pieces.Black.King
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[king][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(king) + offset
			if squareTo >= 0 && squareTo < 64 {
				board.AttackedSquares[byte(squareTo)] = true
			}
		}
	}

	var pawns []byte
	var dir int8
	if whiteToMove {
		pawns = board.Pieces.White.Pawns
		dir = -8
	} else {
		pawns = board.Pieces.Black.Pawns
		dir = 8
	}
	for _, square := range pawns {
		file := board.SquareToFile(square)
		for _, df := range []int8{-1, 1} {
			attackFile := int8(file) + df
			if attackFile < 0 || attackFile > 7 {
				continue
			}
			attack := int8(square) + dir + df
			if attack >= 0 && attack < 64 {
				board.AttackedSquares[byte(attack)] = true
			}
		}
	}
}
