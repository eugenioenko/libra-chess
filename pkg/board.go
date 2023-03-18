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
}

func NewBoard() *Board {
	return &Board{
		Position: [64]byte{},
		CastlingAvailability: CastlingAvailability{
			BlackKingSide:  true,
			BlackQueenSide: true,
			WhiteKingSide:  true,
			WhiteQueenSide: true,
		},
		Kings: KingLocations{
			White: 0,
			Black: 0,
		},
		WhiteToMove: true,
		Moves:       []*Move{},
	}
}

// https://en.wikipedia.org/wiki/Forsyth–Edwards_Notation
func (board *Board) LoadInitial() (bool, error) {
	return board.LoadFromFEN(BoardInitialFEN)
}

func (board *Board) LoadFromFEN(fen string) (bool, error) {
	board.removePieces(0, 64)
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
		fmt.Print(PieceCodeToFont(piece))
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

func (board *Board) IsPieceAtSquareBlack(square byte) bool {
	return board.Position[square] >= 98
}

func (board *Board) AddMove(move *Move) {
	board.Moves = append(board.Moves, move)
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

		// one move forward
		var amountToMove int8 = 8
		var dirToMove int8 = 1
		if !board.WhiteToMove {
			dirToMove = -1
		}
		squareToMove := square - byte(amountToMove*dirToMove)
		if board.IsSquareEmpty(squareToMove) {
			board.AddMove(NewMove(square, squareToMove, MoveQuiet))
		}

		// captures
		var leftCaptureDir int8 = -1
		var rightCaptureDir int8 = 1

		if board.WhiteToMove {
			leftCaptureDir = 1
			rightCaptureDir = -1
		}

		leftSquare := squareToMove + byte(leftCaptureDir)
		rightSquare := squareToMove + byte(rightCaptureDir)

		if board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareBlack(leftSquare) == board.WhiteToMove {
			board.AddMove(NewMove(square, leftSquare, MoveCapture))
		}
		if board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareBlack(rightSquare) == board.WhiteToMove {
			board.AddMove(NewMove(square, rightSquare, MoveCapture))
		}

	}
}
