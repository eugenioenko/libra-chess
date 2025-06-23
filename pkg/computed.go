package libra

// RookRays [square][direction]
var RookRays [64][4]uint64

// BishopRays [square][direction]
var BishopRays [64][4]uint64

// KnightOffsets are the offsets for knight moves from a given square.
// A value of 255 indicates an invalid move (off the board).
var KnightOffsets [64][8]byte

// KingOffsets are the offsets for king moves from a given square.
// A value of 255 indicates an invalid move (off the board).
var KingOffsets [64][8]byte

// BoardDirOffsets are the offsets for each direction on the board.
var BoardDirOffsets = [8]int8{-8, 1, 8, -1, -7, 9, 7, -9} // N, E, S, W, NE, SE, SW, NW

// SquaresToEdge [square][direction]
var SquaresToEdge [64][8]byte

func init() {
	// Initialize SquaresToEdge
	for i := 0; i < 64; i++ {
		rank := i / 8
		file := i % 8
		SquaresToEdge[i][0] = byte(rank)                      // N
		SquaresToEdge[i][1] = byte(7 - file)                  // E
		SquaresToEdge[i][2] = byte(7 - rank)                  // S
		SquaresToEdge[i][3] = byte(file)                      // W
		SquaresToEdge[i][4] = min(byte(rank), byte(7-file))   // NE
		SquaresToEdge[i][5] = min(byte(7-rank), byte(7-file)) // SE
		SquaresToEdge[i][6] = min(byte(7-rank), byte(file))   // SW
		SquaresToEdge[i][7] = min(byte(rank), byte(file))     // NW
	}

	// Initialize KnightOffsets
	for i := 0; i < 64; i++ {
		for j := 0; j < 8; j++ {
			KnightOffsets[i][j] = 255 // Default to invalid
		}
		rank := i / 8
		file := i % 8
		// NNE
		if rank > 1 && file < 7 {
			KnightOffsets[i][0] = byte(i - 15)
		}
		// ENE
		if rank > 0 && file < 6 {
			KnightOffsets[i][1] = byte(i - 6)
		}
		// ESE
		if rank < 7 && file < 6 {
			KnightOffsets[i][2] = byte(i + 10)
		}
		// SSE
		if rank < 6 && file < 7 {
			KnightOffsets[i][3] = byte(i + 17)
		}
		// SSW
		if rank < 6 && file > 0 {
			KnightOffsets[i][4] = byte(i + 15)
		}
		// WSW
		if rank < 7 && file > 1 {
			KnightOffsets[i][5] = byte(i + 6)
		}
		// WNW
		if rank > 0 && file > 1 {
			KnightOffsets[i][6] = byte(i - 10)
		}
		// NNW
		if rank > 1 && file > 0 {
			KnightOffsets[i][7] = byte(i - 17)
		}
	}

	// Initialize KingOffsets
	for i := 0; i < 64; i++ {
		for j := 0; j < 8; j++ {
			KingOffsets[i][j] = 255 // Default to invalid
		}
		rank := i / 8
		file := i % 8
		// N
		if rank > 0 {
			KingOffsets[i][0] = byte(i - 8)
		}
		// E
		if file < 7 {
			KingOffsets[i][1] = byte(i + 1)
		}
		// S
		if rank < 7 {
			KingOffsets[i][2] = byte(i + 8)
		}
		// W
		if file > 0 {
			KingOffsets[i][3] = byte(i - 1)
		}
		// NE
		if rank > 0 && file < 7 {
			KingOffsets[i][4] = byte(i - 7)
		}
		// SE
		if rank < 7 && file < 7 {
			KingOffsets[i][5] = byte(i + 9)
		}
		// SW
		if rank < 7 && file > 0 {
			KingOffsets[i][6] = byte(i + 7)
		}
		// NW
		if rank > 0 && file > 0 {
			KingOffsets[i][7] = byte(i - 9)
		}
	}

	// Initialize RookRays and BishopRays
	for sq := 0; sq < 64; sq++ {
		// Rook Rays
		// North
		for r := sq - 8; r >= 0; r -= 8 {
			RookRays[sq][0] |= (1 << r)
		}
		// East
		for r := sq + 1; r%8 != 0; r++ {
			RookRays[sq][1] |= (1 << r)
		}
		// South
		for r := sq + 8; r < 64; r += 8 {
			RookRays[sq][2] |= (1 << r)
		}
		// West
		for r := sq - 1; r%8 != 7 && r >= 0; r-- {
			RookRays[sq][3] |= (1 << r)
		}

		// Bishop Rays
		// NE
		for r := sq - 7; r >= 0 && r%8 != 0; r -= 7 {
			BishopRays[sq][0] |= (1 << r)
		}
		// SE
		for r := sq + 9; r < 64 && r%8 != 0; r += 9 {
			BishopRays[sq][1] |= (1 << r)
		}
		// SW
		for r := sq + 7; r < 64 && r%8 != 7; r += 7 {
			BishopRays[sq][2] |= (1 << r)
		}
		// NW
		for r := sq - 9; r >= 0 && r%8 != 7; r -= 9 {
			BishopRays[sq][3] |= (1 << r)
		}
	}
}

func min(a, b byte) byte {
	if a < b {
		return a
	}
	return b
}
