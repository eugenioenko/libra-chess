package libra

import (
	"fmt"
	"strings"
)

// Piece Types
const (
	SquareA2 = 48
	SquareA7 = 8
)

const BoardInitialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type KingLocations struct {
	White byte
	Black byte
}

type CastlingAvailability struct {
	BlackKingSide  bool
	BlackQueenSide bool
	WhiteKingSide  bool
	WhiteQueenSide bool
}

type Board struct {
	// 0: h1, 63: a8
	Position             [64]byte
	AttackedSquares      [64]bool
	CastlingAvailability *CastlingAvailability
	Pieces               *PieceColorLocation
	WhiteToMove          bool
	Moves                *BoardMoves
	OnPassant            byte
}

func NewBoard() *Board {
	board := &Board{}
	return board
}

func (board *Board) Initialize() {
	board.Position = [64]byte{}
	board.CastlingAvailability = &CastlingAvailability{
		BlackKingSide:  true,
		BlackQueenSide: true,
		WhiteKingSide:  true,
		WhiteQueenSide: true,
	}
	board.Pieces = NewPieceColorLocation()
	board.WhiteToMove = true
	board.Moves = NewBoardMoves()
	board.OnPassant = 0
}

// https://en.wikipedia.org/wiki/Forsyth–Edwards_Notation
func (board *Board) LoadInitial() (bool, error) {
	board.Initialize()
	return board.LoadFromFEN(BoardInitialFEN)
}

func (board *Board) LoadFromFEN(fen string) (bool, error) {
	board.Initialize()
	parts := strings.Split(fen, " ")
	if len(parts) == 0 {
		return false, fmt.Errorf("invalid FEN, missing blocks, at least piece list block required")
	}

	// 1. Piece placement data
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return false, fmt.Errorf("invalid FEN, missing ranks")
	}

	// TODO validate characters in FEN code are valid

	index := 0
	for _, pieces := range ranks {
		for _, piece := range pieces {
			if CharIsNumber(piece) {
				emptyCount := int(piece - '0')
				board.removePieces(index, emptyCount)
				index += emptyCount
			} else {
				board.Position[index] = byte(piece)
				index += 1
			}
		}
	}
	if index != 64 {
		return false, fmt.Errorf("invalid FEN, missing pieces")
	}

	// 2. Active Color
	if len(parts) > 1 {
		board.WhiteToMove = parts[1] == "w"
	}

	// 3. Castling
	if len(parts) > 2 && parts[2] != "-" {
		// TODO castling
		board.CastlingAvailability = &CastlingAvailability{
			BlackKingSide:  true,
			BlackQueenSide: true,
			WhiteKingSide:  true,
			WhiteQueenSide: true,
		}
	}

	// 4. On-passant
	if len(parts) > 3 && parts[3] != "-" {
		onPassant, ok := SquareNameToIndex(parts[3])
		if ok {
			board.OnPassant = onPassant
		}
	}

	// generate piece list table
	board.GeneratePiecesLocations()
	return true, nil
}

func (board *Board) removePieces(start int, count int) (bool, error) {
	if start+count > 64 {
		return false, fmt.Errorf("invalid remove pieces range, out of range")
	}
	for index := 0; index < count; index++ {
		board.Position[start+index] = 0
	}
	return true, nil
}

func (board *Board) Print() {
	fmt.Println()
	for index, piece := range board.Position {

		if index%8 == 0 {
			fmt.Print(8 - index/8)
			fmt.Print(" | ")
		}
		if piece != 0 {
			fmt.Print(PieceCodeToFont(piece))
		} else {
			fmt.Print(" ")
		}
		fmt.Print(" ")
		if index > 0 && ((index+1)%8) == 0 {
			fmt.Print("\n")
		}
	}
	fmt.Print("   ----------------\n    A B C D E F G H\n\n")
}

func (board *Board) IsSquareValid(square byte) bool {
	return square <= 63
}

func (board *Board) IsSquareEmpty(square byte) bool {
	return board.IsSquareValid(square) && board.Position[square] == 0
}

func (board *Board) IsSquareOccupied(square byte) bool {
	return board.IsSquareValid(square) && board.Position[square] > 0
}

func (board *Board) IsSquareOnPassant(square byte) bool {
	return board.OnPassant == square
}

func (board *Board) IsPieceAtSquareBlack(square byte) bool {
	return board.Position[square] >= 98
}

func (board *Board) IsPieceAtSquareWhite(square byte) bool {
	return board.Position[square] > 0 && board.Position[square] < 98
}

func (board *Board) SquareToRank(square byte) byte {
	return 8 - square/8
}

func (board *Board) IsSquareAtAFile(square byte) bool {
	return square%8 == 0
}

func (board *Board) IsSquareAtHFile(square byte) bool {
	return (square+1)%8 == 0
}

func (board *Board) AddQuietOrCapture(from, to byte, whiteToMove bool) bool {
	if board.IsSquareEmpty(to) {
		board.AddQuietMove(from, to)
		return true
	} else {
		if whiteToMove && board.IsPieceAtSquareBlack(to) {
			board.AddCapture(from, to, MoveCapture, whiteToMove)
		}
		return false
	}
}

func (board *Board) AddMove(move *Move) {
	board.Moves.All = append(board.Moves.All, move)
}

func (board *Board) AddQuietMove(from, to byte) {
	move := NewMove(from, to, MoveQuiet, [2]byte{0, 0})
	board.Moves.All = append(board.Moves.All, move)
	board.Moves.Quiet = append(board.Moves.Quiet, move)
}

func (board *Board) AddCastleMove(from, to byte) {
	move := NewMove(from, to, MoveCastle, [2]byte{0, 0})
	board.Moves.All = append(board.Moves.All, move)
	board.Moves.Quiet = append(board.Moves.Quiet, move)
}

func (board *Board) AddCapture(from, to, moveType byte, whiteToMove bool) {
	captured := board.Position[to]
	if moveType == MoveOnPassant {
		if whiteToMove {
			captured = BlackPawn
		} else {
			captured = WhitePawn
		}
	}
	move := NewMove(from, to, moveType, [2]byte{captured, 0})
	board.AddMove(move)
	board.Moves.Captures = append(board.Moves.Captures, move)
}

func (board *Board) AddPromotion(from, to, captured byte, whiteToMove bool) {
	var moveType byte = MovePromotion
	if captured != 0 {
		moveType = MovePromotionCapture
	}
	var queenPiece byte = WhiteQueen
	var rookPiece byte = WhiteRook
	var bishopPiece byte = WhiteBishop
	var knightPiece byte = WhiteKnight

	if !whiteToMove {
		queenPiece = BlackQueen
		rookPiece = BlackRook
		bishopPiece = BlackBishop
		knightPiece = BlackKnight
	}

	move := NewMove(from, to, moveType, [2]byte{queenPiece, captured})
	board.AddMove(move)
	board.Moves.Promotions = append(board.Moves.Promotions, move)

	move = NewMove(from, to, moveType, [2]byte{rookPiece, captured})
	board.AddMove(move)
	board.Moves.Promotions = append(board.Moves.Promotions, move)

	move = NewMove(from, to, moveType, [2]byte{bishopPiece, captured})
	board.AddMove(move)
	board.Moves.Promotions = append(board.Moves.Promotions, move)

	move = NewMove(from, to, moveType, [2]byte{knightPiece, captured})
	board.AddMove(move)
	board.Moves.Promotions = append(board.Moves.Promotions, move)
}

func (board *Board) GeneratePawnMoves(whiteToMove bool) {
	squares := board.Pieces.White.Pawns
	if !whiteToMove {
		squares = board.Pieces.Black.Pawns
	}

	for _, square := range squares {

		// two space forward, pawns are on initial square
		if whiteToMove {
			if board.SquareToRank(square) == 2 {
				squareToMove := square - 16
				if board.IsSquareEmpty(squareToMove) {
					board.AddQuietMove(square, squareToMove)
				}
			}
		} else {
			if board.SquareToRank(square) == 7 {
				squareToMove := square + 16
				if board.IsSquareEmpty(squareToMove) {
					board.AddQuietMove(square, squareToMove)
				}
			}
		}

		// one move forward and promotion
		if whiteToMove {
			squareToMove := square - 8
			if board.IsSquareEmpty(squareToMove) {
				if board.SquareToRank(squareToMove) == 8 {
					board.AddPromotion(square, squareToMove, 0, whiteToMove)
				} else {
					board.AddQuietMove(square, squareToMove)
				}
			}
		} else {
			squareToMove := square + 8
			if board.IsSquareEmpty(squareToMove) {
				if board.SquareToRank(squareToMove) == 1 {
					board.AddPromotion(square, squareToMove, 0, whiteToMove)
				} else {
					board.AddQuietMove(square, squareToMove)
				}
			}
		}

		// captures and promotion captures
		if whiteToMove {
			// left capture with white
			leftSquare := square - 8 - 1
			if !board.IsSquareAtHFile(square) && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareBlack(leftSquare) {
				if board.SquareToRank(leftSquare) == 8 {
					// promotion capture
					captured := board.Position[leftSquare]
					board.AddPromotion(square, leftSquare, captured, whiteToMove)
				} else {
					// normal capture
					board.AddCapture(square, leftSquare, MoveCapture, whiteToMove)
				}
			}
			// en-passant capture left
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(leftSquare) && board.IsPieceAtSquareBlack(leftSquare+8) {
				board.AddCapture(square, leftSquare, MoveOnPassant, whiteToMove)
			}
			// right capture with white
			rightSquare := square - 8 + 1
			if !board.IsSquareAtAFile(square) && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareBlack(rightSquare) {
				if board.SquareToRank(rightSquare) == 1 {
					// promotion capture
					captured := board.Position[rightSquare]
					board.AddPromotion(square, rightSquare, captured, whiteToMove)
				} else {
					// normal capture
					board.AddCapture(square, rightSquare, MoveCapture, whiteToMove)
				}
			}
			// en-passant capture right
			if !board.IsSquareAtAFile(square) && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare+8) {
				board.AddCapture(square, rightSquare, MoveOnPassant, whiteToMove)
			}
		} else {
			// right capture with black
			rightSquare := square + 8 - 1
			if !board.IsSquareAtHFile(square) && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareWhite(rightSquare) {
				if board.SquareToRank(rightSquare) == 1 {
					captured := board.Position[rightSquare]
					board.AddPromotion(square, rightSquare, captured, whiteToMove)
				} else {
					board.AddCapture(square, rightSquare, MoveCapture, whiteToMove)
				}
			}
			// en-passant capture right
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareWhite(rightSquare-8) {
				board.AddCapture(square, rightSquare, MoveOnPassant, whiteToMove)
			}
			// left capture with black
			leftSquare := square + 8 + 1
			if !board.IsSquareAtAFile(square) && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareWhite(leftSquare) {
				if board.SquareToRank(leftSquare) == 1 {
					captured := board.Position[leftSquare]
					board.AddPromotion(square, leftSquare, captured, whiteToMove)
				} else {
					board.AddCapture(square, leftSquare, MoveCapture, whiteToMove)
				}
			}
			// en passant capture right
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(leftSquare) && board.IsPieceAtSquareWhite(leftSquare-8) {
				board.AddCapture(square, leftSquare, MoveOnPassant, whiteToMove)
			}
		}
	}
}

func (board *Board) GenerateSlidingMoves(pieces []byte, startDir byte, endDir byte, whiteToMove bool) {
	for _, square := range pieces {
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			amountToMove := int8(SquaresToEdge[square][dirOffset])
			for moveIndex := int8(1); moveIndex <= amountToMove; moveIndex++ {
				squareTo := int8(square) + (offset * moveIndex)
				isQuietMove := board.AddQuietOrCapture(square, byte(squareTo), whiteToMove)
				if !isQuietMove {
					break
				}
			}
		}
	}
}

func (board *Board) GenerateKingMoves(whiteToMove bool) {
	square := board.Pieces.White.King
	if !whiteToMove {
		square = board.Pieces.Black.King
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[square][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(square) + offset
			board.AddQuietOrCapture(square, byte(squareTo), whiteToMove)
		}
	}
}

func (board *Board) GenerateCastleMoves(whiteToMove bool) {
	if whiteToMove {
		if (board.CastlingAvailability.WhiteQueenSide) && board.IsSquareEmpty(57) && board.IsSquareEmpty(58) && board.IsSquareEmpty(59) {
			board.AddCastleMove(60, 58)
		} else if (board.CastlingAvailability.WhiteKingSide) && board.IsSquareEmpty(61) && board.IsSquareEmpty(62) {
			board.AddCastleMove(60, 62)
		}
	} else {
		if (board.CastlingAvailability.BlackQueenSide) && board.IsSquareEmpty(1) && board.IsSquareEmpty(2) && board.IsSquareEmpty(3) {
			board.AddCastleMove(4, 2)
		} else if (board.CastlingAvailability.BlackKingSide) && board.IsSquareEmpty(5) && board.IsSquareEmpty(6) {
			board.AddCastleMove(4, 6)
		}
	}
}

func (board *Board) GenerateRookMoves(whiteToMove bool) {
	rooks := board.Pieces.White.Rooks
	if !whiteToMove {
		rooks = board.Pieces.Black.Rooks
	}
	board.GenerateSlidingMoves(rooks, 0, 4, whiteToMove)
}

func (board *Board) GenerateBishopMoves(whiteToMove bool) {
	bishops := board.Pieces.White.Bishops
	if !whiteToMove {
		bishops = board.Pieces.Black.Bishops
	}
	board.GenerateSlidingMoves(bishops, 4, 8, whiteToMove)
}

func (board *Board) GenerateQueenMoves(whiteToMove bool) {
	queues := board.Pieces.White.Queens
	if !whiteToMove {
		queues = board.Pieces.Black.Queens
	}
	board.GenerateSlidingMoves(queues, 0, 8, whiteToMove)
}

func (board *Board) GenerateKnightMoves(whiteToMove bool) {
	knights := board.Pieces.White.Knights
	if !whiteToMove {
		knights = board.Pieces.Black.Knights
	}
	for _, square := range knights {
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				board.AddQuietOrCapture(square, squareTo, whiteToMove)
			}
		}
	}
}

func (board *Board) GenerateMoves() {
	board.Moves = NewBoardMoves()
	board.AttackedSquares = [64]bool{}
	whiteToMove := !board.WhiteToMove

	// generate opposing color attacked squares
	board.GenerateKnightMoves(whiteToMove)
	board.GenerateBishopMoves(whiteToMove)
	board.GenerateRookMoves(whiteToMove)
	board.GenerateQueenMoves(whiteToMove)
	board.GenerateKingMoves(whiteToMove)
	for _, move := range board.Moves.Quiet {
		board.AttackedSquares[move.To] = true
	}

	// generate opposing color pawn attacked squares
	board.Moves = NewBoardMoves()
	// todo generate pawn attack vectors
	board.GeneratePawnMoves(whiteToMove)
	for _, move := range board.Moves.Captures {
		board.AttackedSquares[move.To] = true
	}
	for _, move := range board.Moves.Promotions {
		if move.MoveType == MovePromotionCapture {
			board.AttackedSquares[move.To] = true
		}
	}

	// generate moving color all moves
	board.Moves = NewBoardMoves()
	board.GeneratePawnMoves(board.WhiteToMove)
	board.GenerateKnightMoves(board.WhiteToMove)
	board.GenerateBishopMoves(board.WhiteToMove)
	board.GenerateRookMoves(board.WhiteToMove)
	board.GenerateQueenMoves(board.WhiteToMove)
	board.GenerateKingMoves(board.WhiteToMove)
	board.GenerateCastleMoves(board.WhiteToMove)
}

func (board *Board) GeneratePiecesLocations() {
	for index, piece := range board.Position {
		switch {
		case piece == WhitePawn:
			board.Pieces.White.Pawns = append(board.Pieces.White.Pawns, byte(index))
		case piece == WhiteKnight:
			board.Pieces.White.Knights = append(board.Pieces.White.Knights, byte(index))
		case piece == WhiteBishop:
			board.Pieces.White.Bishops = append(board.Pieces.White.Bishops, byte(index))
		case piece == WhiteRook:
			board.Pieces.White.Rooks = append(board.Pieces.White.Rooks, byte(index))
		case piece == WhiteQueen:
			board.Pieces.White.Queens = append(board.Pieces.White.Queens, byte(index))
		case piece == WhiteKing:
			board.Pieces.White.King = byte(index)
		case piece == BlackPawn:
			board.Pieces.Black.Pawns = append(board.Pieces.Black.Pawns, byte(index))
		case piece == BlackKnight:
			board.Pieces.Black.Knights = append(board.Pieces.Black.Knights, byte(index))
		case piece == BlackBishop:
			board.Pieces.Black.Bishops = append(board.Pieces.Black.Bishops, byte(index))
		case piece == BlackRook:
			board.Pieces.Black.Rooks = append(board.Pieces.Black.Rooks, byte(index))
		case piece == BlackQueen:
			board.Pieces.Black.Queens = append(board.Pieces.Black.Queens, byte(index))
		case piece == BlackKing:
			board.Pieces.Black.King = byte(index)
		}
	}
}
