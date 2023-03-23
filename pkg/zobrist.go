package libra

import "math/rand"

type ZobristCastlingAvailability struct {
	BlackKingSide  uint64
	BlackQueenSide uint64
	WhiteKingSide  uint64
	WhiteQueenSide uint64
}

var zobristPieceTable = GenerateZobristPieceTable()
var zobristOnPassantTable = GenerateZobristOnPassantTable()
var zobristWhiteToMove uint64 = rand.Uint64()
var zobristBlackToMove uint64 = rand.Uint64()

var zobristCastlingAvailability ZobristCastlingAvailability = ZobristCastlingAvailability{
	BlackKingSide:  rand.Uint64(),
	BlackQueenSide: rand.Uint64(),
	WhiteKingSide:  rand.Uint64(),
	WhiteQueenSide: rand.Uint64(),
}

func GenerateZobristPieceTable() [64]map[byte]uint64 {
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

func GenerateZobristOnPassantTable() [64]uint64 {
	table := [64]uint64{}
	for index := 0; index < 64; index++ {
		table[index] = rand.Uint64()
	}
	return table
}

func ZobristHash(board *Board) uint64 {
	var hash uint64 = 0

	// hash pieces
	for index, piece := range board.Position {
		if piece != 0 {
			code, ok := zobristPieceTable[index][piece]
			if ok {
				hash ^= code
			}
		}
	}

	// hash castling availability
	if board.CastlingAvailability.BlackKingSide {
		hash ^= zobristCastlingAvailability.BlackKingSide
	}
	if board.CastlingAvailability.BlackQueenSide {
		hash ^= zobristCastlingAvailability.BlackQueenSide
	}
	if board.CastlingAvailability.WhiteKingSide {
		hash ^= zobristCastlingAvailability.WhiteKingSide
	}
	if board.CastlingAvailability.WhiteKingSide {
		hash ^= zobristCastlingAvailability.WhiteKingSide
	}

	// hash on passant
	if board.OnPassant != 0 {
		hash ^= zobristOnPassantTable[board.OnPassant]
	}

	// hash color to move
	if board.WhiteToMove {
		hash ^= zobristWhiteToMove
	} else {
		hash ^= zobristBlackToMove
	}

	return hash
}
