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
	CastlingAvailability *CastlingAvailability
	Pieces               *PieceColorLocation
	WhiteToMove          bool
	Moves                []*Move
	Captures             []*Move
	Promotions           []*Move
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
	board.Moves = []*Move{}
	board.Captures = []*Move{}
	board.Promotions = []*Move{}
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
				if byte(piece) == WhiteKing {
					board.Pieces.White.King = byte(index)
				} else if byte(piece) == BlackKing {
					board.Pieces.Black.King = byte(index)
				}
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

func (board *Board) IsPieceAt8thRank(square byte) bool {
	return square < 8
}

func (board *Board) IsPieceAt1tstRank(square byte) bool {
	return square >= 56
}

func (board *Board) IsPieceAtAFile(square byte) bool {
	return square%8 == 0
}

func (board *Board) IsPieceAtHFile(square byte) bool {
	return (square+1)%8 == 0
}

func (board *Board) AddQuiteOrCapture(from, to byte) bool {
	if board.IsSquareEmpty(to) {
		board.AddMove(NewMove(from, to, MoveQuiet))
		return true
	} else {
		if board.WhiteToMove && board.IsPieceAtSquareBlack(to) {
			board.AddCapture(NewMove(from, to, MoveCapture))
		}
		return false
	}
}

func (board *Board) AddMove(move *Move) {
	board.Moves = append(board.Moves, move)
}

func (board *Board) AddCapture(move *Move) {
	board.AddMove(move)
	board.Captures = append(board.Captures, move)
}

func (board *Board) AddPromotion(move *Move) {
	board.Moves = append(board.Moves, move)
	board.Promotions = append(board.Promotions, move)
}

func (board *Board) GenerateMoves() {
	board.GeneratePawnMoves()
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
			board.Pieces.Black.King = byte(index)
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

func (board *Board) GeneratePawnMoves() {
	squares := board.Pieces.White.Pawns
	if !board.WhiteToMove {
		squares = board.Pieces.Black.Pawns
	}

	for _, square := range squares {
		// two squares forward
		var leftRankSquare byte = SquareA2
		if !board.WhiteToMove {
			leftRankSquare = SquareA7
		}
		if square >= leftRankSquare && square <= leftRankSquare+8 {
			var amountToMove int8 = 16
			if !board.WhiteToMove {
				amountToMove = -16
			}
			squareToMove := square - byte(amountToMove)
			if board.IsSquareEmpty(squareToMove) {
				board.AddMove(NewMove(square, squareToMove, MoveQuiet))
			}
		}

		// one move forward and promotion
		var amountToMove int8 = 8
		var dirToMove int8 = 1
		if !board.WhiteToMove {
			dirToMove = -1
		}
		squareToMove := square - byte(amountToMove*dirToMove)
		if board.IsSquareEmpty(squareToMove) {
			if (board.WhiteToMove && squareToMove <= 7) || (!board.WhiteToMove && squareToMove >= 56) {
				board.AddMove(NewMove(square, squareToMove, MovePromotion))
			} else {
				board.AddMove(NewMove(square, squareToMove, MoveQuiet))
			}
		}

		// captures
		if board.WhiteToMove {
			leftSquare := square - 8 - 1
			if square%8 != 0 && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareBlack(leftSquare) {

				if leftSquare <= 7 {
					board.AddPromotion(NewMove(square, leftSquare, MovePromotionCapture))
				} else {
					board.AddCapture(NewMove(square, leftSquare, MoveCapture))
				}
			}
			if square%8 != 0 && board.IsSquareOnPassant(leftSquare) && board.IsPieceAtSquareBlack(leftSquare+8) {
				board.AddCapture(NewMove(square, leftSquare, MoveOnPassant))
			}
			rightSquare := square - 8 + 1
			if square%7 != 0 && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareBlack(rightSquare) {
				if leftSquare <= 7 {
					board.AddPromotion(NewMove(square, rightSquare, MovePromotionCapture))
				} else {
					board.AddCapture(NewMove(square, rightSquare, MoveCapture))
				}
			}
			if square%7 != 0 && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare+8) {
				board.AddCapture(NewMove(square, rightSquare, MoveOnPassant))
			}
		} else {
			rightSquare := square + 8 - 1
			if square%8 != 0 && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareBlack(rightSquare) {
				if rightSquare >= 56 {
					board.AddCapture(NewMove(square, rightSquare, MovePromotionCapture))
				} else {
					board.AddCapture(NewMove(square, rightSquare, MoveCapture))
				}
			}
			if square%8 != 0 && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare-8) {
				board.AddCapture(NewMove(square, rightSquare, MoveOnPassant))
			}
			leftSquare := square + 8 + 1
			if square%7 != 0 && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareBlack(leftSquare) {
				if rightSquare >= 56 {
					board.AddCapture(NewMove(square, leftSquare, MovePromotionCapture))
				} else {
					board.AddCapture(NewMove(square, leftSquare, MoveCapture))
				}
			}
			if square%8 != 0 && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare-8) {
				board.AddCapture(NewMove(square, rightSquare, MoveOnPassant))
			}
		}
	}
}

func (board *Board) GenerateRookMoves() {
	rooks := board.Pieces.White.Rooks
	if !board.WhiteToMove {
		rooks = board.Pieces.Black.Rooks
	}
	for _, rook := range rooks {
		var square int8 = int8(rook) - 8
		// up
		if !board.IsPieceAt8thRank(rook) {
			for {
				canContinue := board.AddQuiteOrCapture(rook, byte(square))
				square -= 8
				if square < 0 || !canContinue {
					break
				}
			}
		}
		// down
		if !board.IsPieceAt1tstRank(rook) {
			square = int8(rook) + 8
			for {
				canContinue := board.AddQuiteOrCapture(rook, byte(square))
				square += 8
				if square >= 64 || !canContinue {
					break
				}
			}
		}
		// left
		if !board.IsPieceAtAFile(rook) {
			square = int8(rook) - 1
			for {
				canContinue := board.AddQuiteOrCapture(rook, byte(square))
				if square%8 == 0 || !canContinue {
					break
				}
				square -= 1
			}
		}
		// right
		if !board.IsPieceAtHFile(rook) {
			square = int8(rook) + 1
			for {
				canContinue := board.AddQuiteOrCapture(rook, byte(square))
				if (square+1)%8 == 0 || !canContinue {
					break
				}
				square += 1
			}
		}
	}
}
