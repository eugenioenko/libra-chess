package libra

var RookRays [64][4]uint64 = generateRookRays()
var BishopRays [64][4]uint64 = generateBishopRays()
var SquaresToEdge [64][8]byte = generateSquaresToEdge()
var KnightOffsets [64][8]byte = generateKnightOffsets()
var KingOffsets [64][8]byte = generateKingOffsets()
var BoardDirOffsets [8]int8 = [8]int8{-8, 1, 8, -1, -7, 9, 7, -9}

func generateSquaresToEdge() [64][8]byte {
	squares := [64][8]byte{}
	for i := range squares {
		index := byte(i)
		y := index / 8
		x := index - y*8
		south := 7 - y
		north := y
		west := x
		east := 7 - x
		squares[index][0] = north
		squares[index][1] = east
		squares[index][2] = south
		squares[index][3] = west
		squares[index][4] = MathMinByte(north, east)
		squares[index][5] = MathMinByte(south, east)
		squares[index][6] = MathMinByte(south, west)
		squares[index][7] = MathMinByte(north, west)
	}
	return squares
}

func generateKnightOffsets() [64][8]byte {
	squares := [64][8]byte{}
	jumpOffsets := [8][2]int8{{1, 2}, {-1, 2}, {2, -1}, {-2, -1}, {-1, -2}, {1, -2}, {-2, 1}, {2, 1}}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			squareFrom := y*8 + x
			for offsetIndex, offset := range jumpOffsets {
				x2 := int8(x) + offset[0]
				y2 := int8(y) + offset[1]
				if x2 >= 0 && y2 >= 0 && x2 < 8 && y2 < 8 {
					squares[squareFrom][offsetIndex] = byte(y2*8 + x2)
				} else {
					squares[squareFrom][offsetIndex] = 255
				}
			}
		}
	}
	return squares
}

func generateKingOffsets() [64][8]byte {
	squares := [64][8]byte{}
	jumpOffsets := [8][2]int8{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			squareFrom := y*8 + x
			for offsetIndex, offset := range jumpOffsets {
				x2 := int8(x) + offset[0]
				y2 := int8(y) + offset[1]
				if x2 >= 0 && y2 >= 0 && x2 < 8 && y2 < 8 {
					squares[squareFrom][offsetIndex] = byte(y2*8 + x2)
				} else {
					squares[squareFrom][offsetIndex] = 255
				}
			}
		}
	}
	return squares
}

// Precomputed rays for rook and bishop directions for each square
// Rook: 0=N, 1=E, 2=S, 3=W
// Bishop: 0=NE, 1=SE, 2=SW, 3=NW

func generateRookRays() [64][4]uint64 {
	rookRays := [64][4]uint64{}
	for sq := 0; sq < 64; sq++ {
		file := sq % 8
		rank := sq / 8
		for dir := 0; dir < 4; dir++ {
			var ray uint64
			switch dir {
			case 0: // North
				for r := rank - 1; r >= 0; r-- {
					ray |= 1 << (r*8 + file)
				}
			case 1: // East
				for f := file + 1; f < 8; f++ {
					ray |= 1 << (rank*8 + f)
				}
			case 2: // South
				for r := rank + 1; r < 8; r++ {
					ray |= 1 << (r*8 + file)
				}
			case 3: // West
				for f := file - 1; f >= 0; f-- {
					ray |= 1 << (rank*8 + f)
				}
			}
			rookRays[sq][dir] = ray
		}
	}
	return rookRays
}

func generateBishopRays() [64][4]uint64 {
	bishopRays := [64][4]uint64{}
	for sq := 0; sq < 64; sq++ {
		file := sq % 8
		rank := sq / 8
		for dir := 0; dir < 4; dir++ {
			var ray uint64
			switch dir {
			case 0: // NE
				for r, f := rank-1, file+1; r >= 0 && f < 8; r, f = r-1, f+1 {
					ray |= 1 << (r*8 + f)
				}
			case 1: // SE
				for r, f := rank+1, file+1; r < 8 && f < 8; r, f = r+1, f+1 {
					ray |= 1 << (r*8 + f)
				}
			case 2: // SW
				for r, f := rank+1, file-1; r < 8 && f >= 0; r, f = r+1, f-1 {
					ray |= 1 << (r*8 + f)
				}
			case 3: // NW
				for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
					ray |= 1 << (r*8 + f)
				}
			}
			bishopRays[sq][dir] = ray
		}
	}
	return bishopRays
}
