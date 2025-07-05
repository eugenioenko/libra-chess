package libra

import (
	"math/bits"
	"math/rand"
)

// Use a deterministic RNG for Zobrist hashing
var zobristRNG = rand.New(rand.NewSource(0))

type ZobristCastlingState struct {
	BlackKingSide  uint64
	BlackQueenSide uint64
	WhiteKingSide  uint64
	WhiteQueenSide uint64
}

var zobristPieceTable = GenerateZobristPieceTable()
var zobristOnPassantTable = GenerateZobristOnPassantTable()

var zobristWhiteToMove uint64 = zobristRNG.Uint64()
var zobristCastling ZobristCastlingState = ZobristCastlingState{
	BlackKingSide:  zobristRNG.Uint64(),
	BlackQueenSide: zobristRNG.Uint64(),
	WhiteKingSide:  zobristRNG.Uint64(),
	WhiteQueenSide: zobristRNG.Uint64(),
}

func GenerateZobristPieceTable() [64]map[byte]uint64 {
	table := [64]map[byte]uint64{}
	for index := 0; index < 64; index++ {
		cell := map[byte]uint64{
			WhitePawn:   zobristRNG.Uint64(),
			WhiteKnight: zobristRNG.Uint64(),
			WhiteBishop: zobristRNG.Uint64(),
			WhiteRook:   zobristRNG.Uint64(),
			WhiteQueen:  zobristRNG.Uint64(),
			WhiteKing:   zobristRNG.Uint64(),
			BlackPawn:   zobristRNG.Uint64(),
			BlackKnight: zobristRNG.Uint64(),
			BlackBishop: zobristRNG.Uint64(),
			BlackRook:   zobristRNG.Uint64(),
			BlackQueen:  zobristRNG.Uint64(),
			BlackKing:   zobristRNG.Uint64(),
		}
		table[index] = cell
	}
	return table
}

func GenerateZobristOnPassantTable() [64]uint64 {
	table := [64]uint64{}
	for index := 0; index < 64; index++ {
		table[index] = zobristRNG.Uint64()
	}
	return table
}

func (board *Board) ZobristHash() uint64 {
	hash := board.ZobristHashWasm()
	if board.OnPassant != 0 {
		hash ^= zobristOnPassantTable[board.OnPassant]
	}
	return hash
}

func (board *Board) ZobristHashWasm() uint64 {
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
		hash ^= zobristCastling.BlackKingSide
	}
	if board.Castling.BlackQueenSide {
		hash ^= zobristCastling.BlackQueenSide
	}
	if board.Castling.WhiteKingSide {
		hash ^= zobristCastling.WhiteKingSide
	}
	if board.Castling.WhiteQueenSide {
		hash ^= zobristCastling.WhiteQueenSide
	}
	if board.WhiteToMove {
		hash ^= zobristWhiteToMove
	}
	return hash
}
