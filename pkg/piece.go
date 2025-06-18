package libra

import (
	"math/bits"
)

func PieceCodeToFont(piece byte) string {
	return pieceCodeToFont[piece]
}

type PieceLocation struct {
	Pawns   uint64
	Knights uint64
	Bishops uint64
	Rooks   uint64
	Queens  uint64
	King    uint64
}

type PieceColorLocation struct {
	White PieceLocation
	Black PieceLocation
}

func NewPieceLocation() PieceLocation {
	return PieceLocation{
		Pawns:   0,
		Knights: 0,
		Bishops: 0,
		Rooks:   0,
		Queens:  0,
		King:    0,
	}
}

func (pl *PieceLocation) Clone() PieceLocation {
	return PieceLocation{
		Pawns:   pl.Pawns,
		Knights: pl.Knights,
		Bishops: pl.Bishops,
		Rooks:   pl.Rooks,
		Queens:  pl.Queens,
		King:    pl.King,
	}
}

func NewPieceColorLocation() PieceColorLocation {
	return PieceColorLocation{
		White: NewPieceLocation(),
		Black: NewPieceLocation(),
	}
}

// CountPieces returns the total number of pieces on the board.
func (board *Board) CountPieces() int {
	return bits.OnesCount64(board.WhitePawns) + bits.OnesCount64(board.WhiteKnights) + bits.OnesCount64(board.WhiteBishops) +
		bits.OnesCount64(board.WhiteRooks) + bits.OnesCount64(board.WhiteQueens) + bits.OnesCount64(board.WhiteKing) +
		bits.OnesCount64(board.BlackPawns) + bits.OnesCount64(board.BlackKnights) + bits.OnesCount64(board.BlackBishops) +
		bits.OnesCount64(board.BlackRooks) + bits.OnesCount64(board.BlackQueens) + bits.OnesCount64(board.BlackKing)
}

func (pcl *PieceColorLocation) Clone() PieceColorLocation {
	return PieceColorLocation{
		White: pcl.White.Clone(),
		Black: pcl.Black.Clone(),
	}
}
