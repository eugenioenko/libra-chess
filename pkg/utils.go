package libra

func CharIsNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

var squareNameToIndex map[string]byte = map[string]byte{
	"a8": 0,
	"b8": 1,
	"c8": 2,
	"d8": 3,
	"e8": 4,
	"f8": 5,
	"g8": 6,
	"h8": 7,
	"a7": 8,
	"b7": 9,
	"c7": 10,
	"d7": 11,
	"e7": 12,
	"f7": 13,
	"g7": 14,
	"h7": 15,
	"a6": 16,
	"b6": 17,
	"c6": 18,
	"d6": 19,
	"e6": 20,
	"f6": 21,
	"g6": 22,
	"h6": 23,
	"a5": 24,
	"b5": 25,
	"c5": 26,
	"d5": 27,
	"e5": 28,
	"f5": 29,
	"g5": 30,
	"h5": 31,
	"a4": 32,
	"b4": 33,
	"c4": 34,
	"d4": 35,
	"e4": 36,
	"f4": 37,
	"g4": 38,
	"h4": 39,
	"a3": 40,
	"b3": 41,
	"c3": 42,
	"d3": 43,
	"e3": 44,
	"f3": 45,
	"g3": 46,
	"h3": 47,
	"a2": 48,
	"b2": 49,
	"c2": 50,
	"d2": 51,
	"e2": 52,
	"f2": 53,
	"g2": 54,
	"h2": 55,
	"a1": 56,
	"b1": 57,
	"c1": 58,
	"d1": 59,
	"e1": 60,
	"f1": 61,
	"g1": 62,
	"h1": 63,
}

func SquareNameToIndex(name string) (byte, bool) {
	val, ok := squareNameToIndex[name]
	return val, ok
}

func MathMinByte(a byte, b byte) byte {
	if a > b {
		return b
	}
	return a
}

const (
	SquareA8 = 0
	SquareB8 = 1
	SquareC8 = 2
	SquareD8 = 3
	SquareE8 = 4
	SquareF8 = 5
	SquareG8 = 6
	SquareH8 = 7
	SquareA7 = 8
	SquareB7 = 9
	SquareC7 = 10
	SquareD7 = 11
	SquareE7 = 12
	SquareF7 = 13
	SquareG7 = 14
	SquareH7 = 15
	SquareA6 = 16
	SquareB6 = 17
	SquareC6 = 18
	SquareD6 = 19
	SquareE6 = 20
	SquareF6 = 21
	SquareG6 = 22
	SquareH6 = 23
	SquareA5 = 24
	SquareB5 = 25
	SquareC5 = 26
	SquareD5 = 27
	SquareE5 = 28
	SquareF5 = 29
	SquareG5 = 30
	SquareH5 = 31
	SquareA4 = 32
	SquareB4 = 33
	SquareC4 = 34
	SquareD4 = 35
	SquareE4 = 36
	SquareF4 = 37
	SquareG4 = 38
	SquareH4 = 39
	SquareA3 = 40
	SquareB3 = 41
	SquareC3 = 42
	SquareD3 = 43
	SquareE3 = 44
	SquareF3 = 45
	SquareG3 = 46
	SquareH3 = 47
	SquareA2 = 48
	SquareB2 = 49
	SquareC2 = 50
	SquareD2 = 51
	SquareE2 = 52
	SquareF2 = 53
	SquareG2 = 54
	SquareH2 = 55
	SquareA1 = 56
	SquareB1 = 57
	SquareC1 = 58
	SquareD1 = 59
	SquareE1 = 60
	SquareF1 = 61
	SquareG1 = 62
	SquareH1 = 63
)

// SquareIndexToName converts a square index (0-63) to its algebraic name (e.g., 0 -> "a8").
func SquareIndexToName(idx byte) (string, bool) {
	if idx >= 64 {
		return "", false
	}
	return BoardSquareNames[idx], true
}
