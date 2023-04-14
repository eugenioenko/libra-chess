package libra

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

func PieceCodeToFont(piece byte) string {
	return pieceCodeToFont[piece]
}

type PieceLocation struct {
	Pawns   []byte
	Knights []byte
	Bishops []byte
	Rooks   []byte
	Queens  []byte
	King    byte
}

type PieceColorLocation struct {
	White PieceLocation
	Black PieceLocation
}

func NewPieceLocation() PieceLocation {
	return PieceLocation{
		Pawns:   []byte{},
		Knights: []byte{},
		Bishops: []byte{},
		Rooks:   []byte{},
		Queens:  []byte{},
		King:    0,
	}
}

func NewPieceColorLocation() PieceColorLocation {
	return PieceColorLocation{
		White: NewPieceLocation(),
		Black: NewPieceLocation(),
	}
}
