package libra

import "math/rand"

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

func ZobristHash(position [64]byte) uint64 {
	var hash uint64 = 0
	for index, piece := range position {
		if piece != 0 {
			code, ok := ZobristTable[index][piece]
			if ok {
				hash ^= code
			}
		}
	}
	return hash
}

var ZobristTable = GenerateZobristTable()
