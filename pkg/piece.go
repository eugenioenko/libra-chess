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

var PieceCodeToValue = map[byte]int{
	WhitePawn:   1,
	WhiteKnight: 3,
	WhiteBishop: 3,
	WhiteRook:   5,
	WhiteQueen:  9,
	WhiteKing:   100,
	BlackPawn:   1,
	BlackKnight: 3,
	BlackBishop: 3,
	BlackRook:   5,
	BlackQueen:  9,
	BlackKing:   100,
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

// UpdatePiecesLocation updates the Pieces field to reflect the current board position.
// This is required after any move or FEN load.
func (board *Board) UpdatePiecesLocation() {
	board.Pieces = NewPieceColorLocation()

	for index, piece := range board.Position {
		switch {
		case piece == WhitePawn:
			board.Pieces.White.Pawns = append(board.Pieces.White.Pawns, byte(index))
		case piece == WhiteKnight:
			board.Pieces.White.Knights = append(board.Pieces.White.Knights, byte(index))
		case piece == WhiteBishop:
			board.Pieces.White.Bishops = append(board.Pieces.White.Bishops, byte(index))
		case piece == WhiteRook:
			board.Pieces.White.Rooks = append(board.Pieces.White.Rooks, byte(index))
		case piece == WhiteQueen:
			board.Pieces.White.Queens = append(board.Pieces.White.Queens, byte(index))
		case piece == WhiteKing:
			board.Pieces.White.King = byte(index)
		case piece == BlackPawn:
			board.Pieces.Black.Pawns = append(board.Pieces.Black.Pawns, byte(index))
		case piece == BlackKnight:
			board.Pieces.Black.Knights = append(board.Pieces.Black.Knights, byte(index))
		case piece == BlackBishop:
			board.Pieces.Black.Bishops = append(board.Pieces.Black.Bishops, byte(index))
		case piece == BlackRook:
			board.Pieces.Black.Rooks = append(board.Pieces.Black.Rooks, byte(index))
		case piece == BlackQueen:
			board.Pieces.Black.Queens = append(board.Pieces.Black.Queens, byte(index))
		case piece == BlackKing:
			board.Pieces.Black.King = byte(index)
		}
	}
}

// CountPieces returns the total number of pieces on the board.
func (board *Board) CountPieces() int {
	count := 0
	for _, piece := range board.Position {
		if piece != 0 {
			count++
		}
	}
	return count
}
