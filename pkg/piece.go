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

func (pl *PieceLocation) Clone() PieceLocation {
	return PieceLocation{
		Pawns:   append([]byte(nil), pl.Pawns...),
		Knights: append([]byte(nil), pl.Knights...),
		Bishops: append([]byte(nil), pl.Bishops...),
		Rooks:   append([]byte(nil), pl.Rooks...),
		Queens:  append([]byte(nil), pl.Queens...),
		King:    pl.King,
	}
}

func NewPieceColorLocation() PieceColorLocation {
	return PieceColorLocation{
		White: NewPieceLocation(),
		Black: NewPieceLocation(),
	}
}
