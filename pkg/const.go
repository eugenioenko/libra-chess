package libra

const BoardInitialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var BoardSquareNames [64]string = [64]string{
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
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

var pieceCodeToFont = map[byte]string{
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

var PieceCodeToNotation = map[byte]string{
	WhitePawn:   "",
	WhiteKnight: "N",
	WhiteBishop: "B",
	WhiteRook:   "R",
	WhiteQueen:  "Q",
	WhiteKing:   "K",
	BlackPawn:   "",
	BlackKnight: "N",
	BlackBishop: "B",
	BlackRook:   "R",
	BlackQueen:  "Q",
	BlackKing:   "K",
}

var PieceCodeToValue = map[byte]int{
	WhitePawn:   100,
	WhiteKnight: 300,
	WhiteBishop: 300,
	WhiteRook:   500,
	WhiteQueen:  900,
	WhiteKing:   0,
	BlackPawn:   100,
	BlackKnight: 300,
	BlackBishop: 300,
	BlackRook:   500,
	BlackQueen:  900,
	BlackKing:   0,
}

var PieceToHistoryIndex = map[byte]int{
	WhitePawn:   1,
	WhiteKnight: 2,
	WhiteBishop: 3,
	WhiteRook:   4,
	WhiteQueen:  5,
	WhiteKing:   6,
	BlackPawn:   7,
	BlackKnight: 8,
	BlackBishop: 9,
	BlackRook:   10,
	BlackQueen:  11,
	BlackKing:   12,
}

var WhitePromotionMap = map[byte]byte{
	'q': WhiteQueen,
	'r': WhiteRook,
	'b': WhiteBishop,
	'n': WhiteKnight,
}

var BlackPromotionMap = map[byte]byte{
	'q': BlackQueen,
	'r': BlackRook,
	'b': BlackBishop,
	'n': BlackKnight,
}

// Piece-Square Tables (simplified, values in centipawns)
var pawnPST = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}
var knightPST = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}
var bishopPST = [64]int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}
var rookPST = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 0, 5, 5, 0, 0, 0,
}
var queenPST = [64]int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}
var kingPST = [64]int{
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	20, 20, 0, 0, 0, 0, 20, 20,
	20, 30, 10, 0, 0, 10, 30, 20,
}
