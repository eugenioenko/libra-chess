package libra

import (
	"math/bits"
	"math/rand"
)

type ZobristCastlingState struct {
	BlackKingSide  uint64
	BlackQueenSide uint64
	WhiteKingSide  uint64
	WhiteQueenSide uint64
}

var zobristPieceTable = GenerateZobristPieceTable()
var zobristOnPassantTable = GenerateZobristOnPassantTable()
var zobristWhiteToMove uint64 = rand.Uint64()

var zobristCastlingState ZobristCastlingState = ZobristCastlingState{
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

func (board *Board) ZobristHash() uint64 {
	var hash uint64 = 0

	b := board.WhitePawns
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][WhitePawn]
		b &= b - 1
	}
	b = board.WhiteKnights
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][WhiteKnight]
		b &= b - 1
	}
	b = board.WhiteBishops
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][WhiteBishop]
		b &= b - 1
	}
	b = board.WhiteRooks
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][WhiteRook]
		b &= b - 1
	}
	b = board.WhiteQueens
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][WhiteQueen]
		b &= b - 1
	}
	b = board.WhiteKing
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][WhiteKing]
		b &= b - 1
	}
	b = board.BlackPawns
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][BlackPawn]
		b &= b - 1
	}
	b = board.BlackKnights
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][BlackKnight]
		b &= b - 1
	}
	b = board.BlackBishops
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][BlackBishop]
		b &= b - 1
	}
	b = board.BlackRooks
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][BlackRook]
		b &= b - 1
	}
	b = board.BlackQueens
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][BlackQueen]
		b &= b - 1
	}
	b = board.BlackKing
	for b != 0 {
		sq := bits.TrailingZeros64(b)
		hash ^= zobristPieceTable[sq][BlackKing]
		b &= b - 1
	}

	if board.Castling.BlackKingSide {
		hash ^= zobristCastlingState.BlackKingSide
	}
	if board.Castling.BlackQueenSide {
		hash ^= zobristCastlingState.BlackQueenSide
	}
	if board.Castling.WhiteKingSide {
		hash ^= zobristCastlingState.WhiteKingSide
	}
	if board.Castling.WhiteQueenSide {
		hash ^= zobristCastlingState.WhiteQueenSide
	}
	if board.OnPassant != 0 {
		hash ^= zobristOnPassantTable[board.OnPassant]
	}
	if board.WhiteToMove {
		hash ^= zobristWhiteToMove
	}
	return hash
}
