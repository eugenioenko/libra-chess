package libra

import (
	"fmt"
	"math/rand"
	"strings"
)

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

type Board struct {
	// 0: h1, 63: a8
	position [64]byte
	// black king side, black queen side, white king side, white queen side
	castling [4]bool
	// white king, black king
	kings       [2]byte
	whiteToMove bool
}

func NewBoard() *Board {
	return &Board{
		position:    [64]byte{},
		castling:    [4]bool{true, true, true, true},
		kings:       [2]byte{},
		whiteToMove: true,
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
				board.position[index] = byte(piece)
				index += 1
				if byte(piece) == WhiteKing {
					board.kings[0] = byte(index)
				} else if byte(piece) == BlackKing {
					board.kings[1] = byte(index)
				}
			}
		}
	}
	if index != 64 {
		return false, fmt.Errorf("invalid FEN, missing pieces")
	}

	// 2. Active Color
	board.whiteToMove = parts[1] == "w"

	return true, nil
}

func (board *Board) removePieces(start int, count int) (bool, error) {
	if start+count > 64 {
		return false, fmt.Errorf("invalid remove pieces range, out of range")
	}
	for index := 0; index < count; index++ {
		board.position[start+index] = 0
	}
	return true, nil
}

func (board *Board) ZobristHash() uint64 {
	var hash uint64 = 0
	for index, piece := range board.position {
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
	for index, piece := range board.position {

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
