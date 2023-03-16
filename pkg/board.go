package libra

import (
	"fmt"
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

const BoardInitialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type Board struct {
	position    [64]byte
	whiteToMove bool
}

func NewBoard() *Board {
	return &Board{
		position:    [64]byte{},
		whiteToMove: true,
	}
}

// https://en.wikipedia.org/wiki/Forsyth%E2%80%93Edwards_Notation
func (board *Board) LoadInitial() (bool, error) {
	return board.LoadFromFEN(BoardInitialFEN)
}

func (board *Board) LoadFromFEN(fen string) (bool, error) {
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
				board.RemovePieces(index, emptyCount)
				index += emptyCount
			} else {
				board.position[index] = byte(piece)
				index += 1
			}
		}
	}
	if index != 64 {
		return false, fmt.Errorf("invalid FEN, missing pieces")
	}

	// 2. Active Color
	if parts[1] == "w" {
		board.whiteToMove = true
	} else {
		board.whiteToMove = false
	}

	return true, nil
}

func (board *Board) RemovePieces(start int, count int) (bool, error) {
	if start+count >= 64 {
		return false, fmt.Errorf("invalid remove pieces range, out of range")
	}
	for index := 0; index < count; index++ {
		board.position[start+index] = 0
	}
	return true, nil
}
