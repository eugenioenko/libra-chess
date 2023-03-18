package libra

import (
	"fmt"
	"math/rand"
	"strings"
)

// Piece Types
const (
	WhitePawn   = 80  // P
	WhiteKnight = 78  // N
	WhiteBishop = 66  // B
	WhiteRook   = 82  // R
	WhiteQueen  = 81  // Q
	WhiteKing   = 75  // K
	BlackPawn   = 112 // p
	BlackKnight = 110 // n
	BlackBishop = 98  // b
	BlackRook   = 114 // r
	BlackQueen  = 113 // q
	BlackKing   = 107 // k

	SquareA2 = 48
	SquareA7 = 8
)

// Move Types
const (
	MoveQuiet = iota
	MoveCapture
	MoveCastle
	MovePromotion
)

var PieceCodeToFont = map[byte]string{
	WhitePawn:   "♟︎",
	WhiteKnight: "♞",
	WhiteBishop: "♝",
	WhiteRook:   "♜",
	WhiteQueen:  "♛",
	WhiteKing:   "♚",
	BlackPawn:   "♙",
	BlackKnight: "♘",
	BlackBishop: "♗",
	BlackRook:   "♖",
	BlackQueen:  "♕",
	BlackKing:   "♔",
}

func GenerateZobristTable() [64]map[byte]uint64 {
	table := [64]map[byte]uint64{}
	for index := 0; index < 64; index++ {
		cell := map[byte]uint64{
			WhitePawn:   rand.Uint64(),
			WhiteKnight: rand.Uint64(),
			WhiteBishop: rand.Uint64(),
			WhiteRook:   rand.Uint64(),
			WhiteQueen:  rand.Uint64(),
			WhiteKing:   rand.Uint64(),
			BlackPawn:   rand.Uint64(),
			BlackKnight: rand.Uint64(),
			BlackBishop: rand.Uint64(),
			BlackRook:   rand.Uint64(),
			BlackQueen:  rand.Uint64(),
			BlackKing:   rand.Uint64(),
		}
		table[index] = cell
	}
	return table
}

var ZobristTable = GenerateZobristTable()

const BoardInitialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type KingLocations struct {
	White byte
	Black byte
}

type Move struct {
	From byte
	To   byte
	Kind byte
	Code byte
}

func NewMove(from byte, to byte, kind byte) *Move {
	return &Move{
		From: from,
		To:   to,
		Kind: kind,
		Code: 0,
	}
}

func NewPromotionMove(from byte, to byte, kind byte, code byte) *Move {
	return &Move{
		From: from,
		To:   to,
		Kind: kind,
		Code: code,
	}
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

func (board *Board) ZobristHash() uint64 {
	var hash uint64 = 0
	for index, piece := range board.Position {
		if piece != 0 {
			code, ok := ZobristTable[index][piece]
			if ok {
				hash ^= code
			}
		}
	}
	return hash
}

func (board *Board) Print() {
	fmt.Println()
	for index, piece := range board.Position {

		if index%8 == 0 {
			fmt.Print(8 - index/8)
			fmt.Print(" | ")
		}
		fmt.Print(PieceCodeToFont[piece])
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

func (board *Board) IsPieceAtSquareWhite(square byte) bool {
	return board.Position[square] < 98
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
		leftSquare := squareToMove - 1
		rightSquare := squareToMove + 1
		captureWhite := false
		if board.WhiteToMove {
			captureWhite = true
		}
		if board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareWhite(leftSquare) == captureWhite {
			board.AddMove(NewMove(square, leftSquare, MoveCapture))
		}
		if board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareWhite(rightSquare) == captureWhite {
			board.AddMove(NewMove(square, rightSquare, MoveCapture))
		}

	}
}
