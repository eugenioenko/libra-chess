package libra

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
// TODO: this function is currently the slowest part of the board update process.
// It appends the location of each piece to the corresponding slice in the Pieces field which can be inefficient.
// Preallocating slices for each piece type does improve performance but not significantly and the resulting code is less readable.
// Consider using a more efficient data structure or approach if performance becomes a concern.
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

func (pcl *PieceColorLocation) Clone() PieceColorLocation {
	return PieceColorLocation{
		White: pcl.White.Clone(),
		Black: pcl.Black.Clone(),
	}
}
