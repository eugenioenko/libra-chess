package libra

const (
	MoveQuiet = iota
	MoveCapture
	MoveOnPassant
	MovePromotion
	MovePromotionCapture
	MoveCastle
)

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

type BoardMoves struct {
	All        []*Move
	Quite      []*Move
	Captures   []*Move
	Promotions []*Move
}

func NewBoardMoves() *BoardMoves {
	return &BoardMoves{
		All:        []*Move{},
		Quite:      []*Move{},
		Captures:   []*Move{},
		Promotions: []*Move{},
	}
}

// N, S, E, W, NE, SE, SW, SE

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

func generateKnightJumps() [64][8]byte {
	squares := [64][8]byte{}
	jumpOffsets := [8][2]int8{{1, 2}, {-1, 2}, {2, -1}, {-2, -1}, {-1, -2}, {1, -2}, {-2, 1}, {2, 1}}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			squareFrom := y*8 + x
			for offsetIndex, offset := range jumpOffsets {
				x2 := int8(x) + offset[0]
				y2 := int8(y) + offset[1]
				if x2 >= 0 && y2 >= 0 {
					squares[squareFrom][offsetIndex] = byte(y2*8 + x2)
				} else {
					squares[squareFrom][offsetIndex] = 255
				}
			}
		}
	}
	return squares
}

var SquaresToEdge [64][8]byte = generateSquaresToEdge()
var BoardDirOffsets [8]int8 = [8]int8{-8, 1, 8, -1, -7, 9, 7, -9}
var SquareKnightJumps [64][8]byte = generateKnightJumps()
