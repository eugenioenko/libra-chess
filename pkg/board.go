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
	CastlingAvailability CastlingAvailability
	Kings                KingLocations
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
	board.CastlingAvailability = CastlingAvailability{
		BlackKingSide:  true,
		BlackQueenSide: true,
		WhiteKingSide:  true,
		WhiteQueenSide: true,
	}
	board.Kings = KingLocations{
		White: 0,
		Black: 0,
	}
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
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid FEN, missing blocks")
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
					board.Kings.White = byte(index)
				} else if byte(piece) == BlackKing {
					board.Kings.Black = byte(index)
				}
			}
		}
	}
	if index != 64 {
		return false, fmt.Errorf("invalid FEN, missing pieces")
	}

	// 2. Active Color
	board.WhiteToMove = parts[1] == "w"

	// 3. On-passant
	if parts[3] != "-" {
		onPassant, ok := SquareNameToIndex(parts[3])
		if ok {
			board.OnPassant = onPassant
		}
	}
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

func (board *Board) GenerateMoves() {
	board.GeneratePawnMoves()
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

func (board *Board) GeneratePawnMoves() {
	squares := []byte{}
	for index, piece := range board.Position {
		if board.WhiteToMove && piece == WhitePawn {
			squares = append(squares, byte(index))
		} else if !board.WhiteToMove && piece == BlackPawn {
			squares = append(squares, byte(index))
		}
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
