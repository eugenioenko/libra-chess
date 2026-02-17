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

// Game phase weights for tapered evaluation
const (
	PawnPhase   = 0
	KnightPhase = 1
	BishopPhase = 1
	RookPhase   = 2
	QueenPhase  = 4
	TotalPhase  = 24 // 4*1(knights) + 4*1(bishops) + 4*2(rooks) + 2*4(queens)
)

// PeSTO base material values (added to PST in init)
var mgPieceValue = [6]int{82, 337, 365, 477, 1025, 0} // pawn, knight, bishop, rook, queen, king
var egPieceValue = [6]int{94, 281, 297, 512, 936, 0}

// PeSTO Piece-Square Tables (positional component, material added in init)
// Indexed from white's perspective: a8=0, h1=63

// Middlegame tables
var mgPawnPST = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	98, 134, 61, 95, 68, 126, 34, -11,
	-6, 7, 26, 31, 65, 56, 25, -20,
	-14, 13, 6, 21, 23, 12, 17, -23,
	-27, -2, -5, 12, 17, 6, 10, -25,
	-26, -4, -4, -10, 3, 3, 33, -12,
	-35, -1, -20, -23, -15, 24, 38, -22,
	0, 0, 0, 0, 0, 0, 0, 0,
}
var egPawnPST = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	178, 173, 158, 134, 147, 132, 165, 187,
	94, 100, 85, 67, 56, 53, 82, 84,
	32, 24, 13, 5, -2, 4, 17, 17,
	13, 9, -3, -7, -7, -8, 3, -1,
	4, 7, -6, 1, 0, -5, -1, -8,
	13, 8, 8, 10, 13, 0, 2, -7,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var mgKnightPST = [64]int{
	-167, -89, -34, -49, 61, -97, -15, -107,
	-73, -41, 72, 36, 23, 62, 7, -17,
	-47, 60, 37, 65, 84, 129, 73, 44,
	-9, 17, 19, 53, 37, 69, 18, 22,
	-13, 4, 16, 13, 28, 19, 21, -8,
	-23, -9, 12, 10, 19, 17, 25, -16,
	-29, -53, -12, -3, -1, 18, -14, -19,
	-105, -21, -58, -33, -17, -28, -19, -23,
}
var egKnightPST = [64]int{
	-58, -38, -13, -28, -31, -27, -63, -99,
	-25, -8, -25, -2, -9, -25, -24, -52,
	-24, -20, 10, 9, -1, -9, -19, -41,
	-17, 3, 22, 22, 22, 11, 8, -18,
	-18, -6, 16, 25, 16, 17, 4, -18,
	-23, -3, -1, 15, 10, -3, -20, -22,
	-42, -20, -10, -5, -2, -20, -23, -44,
	-29, -51, -23, -15, -22, -18, -50, -64,
}

var mgBishopPST = [64]int{
	-29, 4, -82, -37, -25, -42, 7, -8,
	-26, 16, -18, -13, 30, 59, 18, -47,
	-16, 37, 43, 40, 35, 50, 37, -2,
	-4, 5, 19, 50, 37, 37, 7, -2,
	-6, 13, 13, 26, 34, 12, 10, 4,
	0, 15, 15, 15, 14, 27, 18, 10,
	4, 15, 16, 0, 7, 21, 33, 1,
	-33, -3, -14, -21, -13, -12, -39, -21,
}
var egBishopPST = [64]int{
	-14, -21, -11, -8, -7, -9, -17, -24,
	-8, -4, 7, -12, -3, -13, -4, -14,
	2, -8, 0, -1, -2, 6, 0, 4,
	-3, 9, 12, 9, 14, 10, 3, 2,
	-6, 3, 13, 19, 7, 10, -3, -9,
	-12, -3, 8, 10, 13, 3, -7, -15,
	-14, -18, -7, -1, 4, -9, -15, -27,
	-23, -9, -23, -5, -9, -16, -5, -17,
}

var mgRookPST = [64]int{
	32, 42, 32, 51, 63, 9, 31, 43,
	27, 32, 58, 62, 80, 67, 26, 44,
	-5, 19, 26, 36, 17, 45, 61, 16,
	-24, -11, 7, 26, 24, 35, -8, -20,
	-36, -26, -12, -1, 9, -7, 6, -23,
	-45, -25, -16, -17, 3, 0, -5, -33,
	-44, -16, -20, -9, -1, 11, -6, -71,
	-19, -13, 1, 17, 16, 7, -37, -26,
}
var egRookPST = [64]int{
	13, 10, 18, 15, 12, 12, 8, 5,
	11, 13, 13, 11, -3, 3, 8, 3,
	7, 7, 7, 5, 4, -3, -5, -3,
	4, 3, 13, 1, 2, 1, -1, 2,
	3, 5, 8, 4, -5, -6, -8, -11,
	-4, 0, -5, -1, -7, -12, -8, -16,
	-6, -6, 0, 2, -9, -9, -11, -3,
	-9, 2, 3, -1, -5, -13, 4, -20,
}

var mgQueenPST = [64]int{
	-28, 0, 29, 12, 59, 44, 43, 45,
	-24, -39, -5, 1, -16, 57, 28, 54,
	-13, -17, 7, 8, 29, 56, 47, 57,
	-27, -27, -16, -16, -1, 17, -2, 1,
	-9, -26, -9, -10, -2, -4, 3, -3,
	-14, 2, -11, -2, -5, 2, 14, 5,
	-35, -8, 11, 2, 8, 15, -3, 1,
	-1, -18, -9, 10, -15, -25, -31, -50,
}
var egQueenPST = [64]int{
	-9, 22, 22, 27, 27, 19, 10, 20,
	-17, 20, 32, 41, 58, 25, 30, 0,
	-20, 6, 9, 49, 47, 35, 19, 9,
	3, 22, 24, 45, 57, 40, 57, 36,
	-18, 28, 19, 47, 31, 34, 39, 23,
	-16, -27, 15, 6, 9, 17, 10, 5,
	-22, -23, -30, -16, -16, -23, -36, -32,
	-33, -28, -22, -43, -5, -32, -20, -41,
}

var mgKingPST = [64]int{
	-65, 23, 16, -15, -56, -34, 2, 13,
	29, -1, -20, -7, -8, -4, -38, -29,
	-9, 24, 2, -16, -20, 6, 22, -22,
	-17, -20, -12, -27, -30, -25, -14, -36,
	-49, -1, -27, -39, -46, -44, -33, -51,
	-14, -14, -22, -46, -44, -30, -15, -27,
	1, 7, -8, -64, -43, -16, 9, 8,
	-15, 36, 12, -54, 8, -28, 24, 14,
}
var egKingPST = [64]int{
	-74, -35, -18, -18, -11, 15, 4, -17,
	-12, 17, 14, 17, 17, 38, 23, 11,
	10, 17, 23, 15, 20, 45, 44, 13,
	-8, 22, 24, 27, 26, 33, 26, 3,
	-18, -4, 21, 24, 27, 23, 9, -11,
	-19, -3, 11, 21, 23, 16, 7, -9,
	-27, -11, 4, 13, 14, 4, -5, -17,
	-53, -34, -21, -11, -28, -14, -24, -43,
}

func init() {
	// Bake material values into PST tables
	type pstPair struct {
		mg    *[64]int
		eg    *[64]int
		piece int // index into mgPieceValue/egPieceValue
	}
	pairs := []pstPair{
		{&mgPawnPST, &egPawnPST, 0},
		{&mgKnightPST, &egKnightPST, 1},
		{&mgBishopPST, &egBishopPST, 2},
		{&mgRookPST, &egRookPST, 3},
		{&mgQueenPST, &egQueenPST, 4},
		{&mgKingPST, &egKingPST, 5},
	}
	for _, p := range pairs {
		for sq := 0; sq < 64; sq++ {
			p.mg[sq] += mgPieceValue[p.piece]
			p.eg[sq] += egPieceValue[p.piece]
		}
	}
}
