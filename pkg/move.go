package libra

import (
	"fmt"
)

const (
	MoveQuiet = iota
	MoveCapture
	MoveEnPassant
	MoveCastle
	MovePromotion
	MovePromotionCapture
)

type Move struct {
	Piece     byte
	From      byte
	To        byte
	MoveType  byte
	Promotion byte
	Captured  byte
}

type MovesCount struct {
	All       int
	Quiet     int
	Capture   int
	Promotion int
}

func NewMovesCount() *MovesCount {
	return &MovesCount{All: 0,
		Quiet:     0,
		Capture:   0,
		Promotion: 0,
	}
}

func NewMove(piece, from, to, moveType, promotion, captured byte) Move {
	return Move{
		Piece:     piece,
		From:      from,
		To:        to,
		MoveType:  moveType,
		Captured:  captured,
		Promotion: promotion,
	}
}

type MoveState struct {
	WhitePawns      uint64
	WhiteKnights    uint64
	WhiteBishops    uint64
	WhiteRooks      uint64
	WhiteQueens     uint64
	WhiteKing       uint64
	BlackPawns      uint64
	BlackKnights    uint64
	BlackBishops    uint64
	BlackRooks      uint64
	BlackQueens     uint64
	BlackKing       uint64
	Castling        CastlingState
	OnPassant       byte
	HalfMoveClock   int
	FullMoveCounter int
	WhiteToMove     bool
}

// CountMoves returns a summary of the number of moves by type in the current move list.
func CountMoves(moves []Move) *MovesCount {
	count := NewMovesCount()
	count.All = len(moves)
	for _, move := range moves {
		if move.MoveType == MoveQuiet {
			count.Quiet += 1
		} else if move.MoveType == MovePromotion || move.MoveType == MovePromotionCapture {
			count.Promotion += 1
		} else if move.MoveType == MoveCapture || move.MoveType == MovePromotionCapture || move.MoveType == MoveEnPassant {
			count.Capture += 1
		}
	}
	return count
}

func (move Move) ToUCI() string {
	from, _ := SquareIndexToName(move.From)
	to, _ := SquareIndexToName(move.To)
	uci := from + to
	if move.MoveType == MovePromotion || move.MoveType == MovePromotionCapture {
		promo := ""
		switch move.Promotion {
		case WhiteQueen, BlackQueen:
			promo = "q"
		case WhiteRook, BlackRook:
			promo = "r"
		case WhiteBishop, BlackBishop:
			promo = "b"
		case WhiteKnight, BlackKnight:
			promo = "n"
		}
		uci += promo
	}
	return uci
}

func (move Move) ToMove() string {
	from := BoardSquareNames[move.From]
	to := BoardSquareNames[move.To]
	piece := pieceCodeToFont[move.Piece]
	capture := " "
	if move.MoveType == MoveCapture || move.MoveType == MovePromotionCapture || move.MoveType == MoveEnPassant {
		capture = "x"
	}
	return fmt.Sprintf("%s%s%s%s\n", piece, from, capture, to)
}

// Move returns a MoveState for undoing the move
func (board *Board) Move(move Move) MoveState {
	// Save current state
	prev := MoveState{
		WhitePawns:      board.WhitePawns,
		WhiteKnights:    board.WhiteKnights,
		WhiteBishops:    board.WhiteBishops,
		WhiteRooks:      board.WhiteRooks,
		WhiteQueens:     board.WhiteQueens,
		WhiteKing:       board.WhiteKing,
		BlackPawns:      board.BlackPawns,
		BlackKnights:    board.BlackKnights,
		BlackBishops:    board.BlackBishops,
		BlackRooks:      board.BlackRooks,
		BlackQueens:     board.BlackQueens,
		BlackKing:       board.BlackKing,
		Castling:        board.Castling,
		OnPassant:       board.OnPassant,
		HalfMoveClock:   board.HalfMoveClock,
		FullMoveCounter: board.FullMoveCounter,
		WhiteToMove:     board.WhiteToMove,
	}

	if !board.WhiteToMove {
		board.FullMoveCounter += 1
	}
	from := move.From
	to := move.To
	piece := board.PieceAtSquare(from)

	if move.MoveType == MoveCapture || board.IsSquarePawn(from) {
		board.HalfMoveClock = 0
	} else {
		board.HalfMoveClock += 1
	}

	// Remove piece from 'from' square
	board.clearPieceAtSquare(from, piece)
	// Place piece at 'to' square
	if move.MoveType == MovePromotion || move.MoveType == MovePromotionCapture {
		board.setPieceAtSquare(to, move.Promotion)
	} else {
		board.setPieceAtSquare(to, piece)
	}

	// Handle en passant
	if move.MoveType == MoveEnPassant {
		if board.WhiteToMove {
			board.clearPieceAtSquare(to+8, BlackPawn)
		} else {
			board.clearPieceAtSquare(to-8, WhitePawn)
		}
	}

	// Handle castling
	if move.MoveType == MoveCastle {
		if piece == WhiteKing {
			if to == SquareG1 {
				board.clearPieceAtSquare(SquareH1, WhiteRook)
				board.setPieceAtSquare(SquareF1, WhiteRook)
			} else if to == SquareC1 {
				board.clearPieceAtSquare(SquareA1, WhiteRook)
				board.setPieceAtSquare(SquareD1, WhiteRook)
			}
			board.Castling.WhiteKingSide = false
			board.Castling.WhiteQueenSide = false
		} else if piece == BlackKing {
			if to == SquareG8 {
				board.clearPieceAtSquare(SquareH8, BlackRook)
				board.setPieceAtSquare(SquareF8, BlackRook)
			} else if to == SquareC8 {
				board.clearPieceAtSquare(SquareA8, BlackRook)
				board.setPieceAtSquare(SquareD8, BlackRook)
			}
			board.Castling.BlackKingSide = false
			board.Castling.BlackQueenSide = false
		}
	}

	board.OnPassant = 0
	if piece == WhitePawn && from/8 == 6 && to/8 == 4 {
		board.OnPassant = from - 8
	} else if piece == BlackPawn && from/8 == 1 && to/8 == 3 {
		board.OnPassant = from + 8
	}

	if piece == WhiteKing {
		board.Castling.WhiteKingSide = false
		board.Castling.WhiteQueenSide = false
	}
	if piece == BlackKing {
		board.Castling.BlackKingSide = false
		board.Castling.BlackQueenSide = false
	}
	if piece == WhiteRook {
		if from == SquareH1 {
			board.Castling.WhiteKingSide = false
		}
		if from == SquareA1 {
			board.Castling.WhiteQueenSide = false
		}
	}
	if piece == BlackRook {
		if from == SquareH8 {
			board.Castling.BlackKingSide = false
		}
		if from == SquareA8 {
			board.Castling.BlackQueenSide = false
		}
	}

	if move.MoveType == MoveCapture || move.MoveType == MovePromotionCapture {
		capturedPieceSquare := move.To
		board.clearPieceAtSquare(capturedPieceSquare, move.Captured)
		if move.Captured == WhiteRook {
			if capturedPieceSquare == SquareH1 {
				board.Castling.WhiteKingSide = false
			}
			if capturedPieceSquare == SquareA1 {
				board.Castling.WhiteQueenSide = false
			}
		}
		if move.Captured == BlackRook {
			if capturedPieceSquare == SquareH8 {
				board.Castling.BlackKingSide = false
			}
			if capturedPieceSquare == SquareA8 {
				board.Castling.BlackQueenSide = false
			}
		}
	}

	board.WhiteToMove = !board.WhiteToMove
	return prev
}

// clearPieceAtSquare removes a piece from a square in the bitboards
func (board *Board) clearPieceAtSquare(square byte, piece byte) {
	mask := ^(uint64(1) << square)
	switch piece {
	case WhitePawn:
		board.WhitePawns &= mask
	case WhiteKnight:
		board.WhiteKnights &= mask
	case WhiteBishop:
		board.WhiteBishops &= mask
	case WhiteRook:
		board.WhiteRooks &= mask
	case WhiteQueen:
		board.WhiteQueens &= mask
	case WhiteKing:
		board.WhiteKing &= mask
	case BlackPawn:
		board.BlackPawns &= mask
	case BlackKnight:
		board.BlackKnights &= mask
	case BlackBishop:
		board.BlackBishops &= mask
	case BlackRook:
		board.BlackRooks &= mask
	case BlackQueen:
		board.BlackQueens &= mask
	case BlackKing:
		board.BlackKing &= mask
	}
}

// setPieceAtSquare places a piece on a square in the bitboards
func (board *Board) setPieceAtSquare(square byte, piece byte) {
	mask := uint64(1) << square
	switch piece {
	case WhitePawn:
		board.WhitePawns |= mask
	case WhiteKnight:
		board.WhiteKnights |= mask
	case WhiteBishop:
		board.WhiteBishops |= mask
	case WhiteRook:
		board.WhiteRooks |= mask
	case WhiteQueen:
		board.WhiteQueens |= mask
	case WhiteKing:
		board.WhiteKing |= mask
	case BlackPawn:
		board.BlackPawns |= mask
	case BlackKnight:
		board.BlackKnights |= mask
	case BlackBishop:
		board.BlackBishops |= mask
	case BlackRook:
		board.BlackRooks |= mask
	case BlackQueen:
		board.BlackQueens |= mask
	case BlackKing:
		board.BlackKing |= mask
	}
}

// UndoMove restores the board to a previous MoveState
func (board *Board) UndoMove(state MoveState) {
	board.WhitePawns = state.WhitePawns
	board.WhiteKnights = state.WhiteKnights
	board.WhiteBishops = state.WhiteBishops
	board.WhiteRooks = state.WhiteRooks
	board.WhiteQueens = state.WhiteQueens
	board.WhiteKing = state.WhiteKing
	board.BlackPawns = state.BlackPawns
	board.BlackKnights = state.BlackKnights
	board.BlackBishops = state.BlackBishops
	board.BlackRooks = state.BlackRooks
	board.BlackQueens = state.BlackQueens
	board.BlackKing = state.BlackKing
	board.Castling = state.Castling
	board.OnPassant = state.OnPassant
	board.HalfMoveClock = state.HalfMoveClock
	board.FullMoveCounter = state.FullMoveCounter
	board.WhiteToMove = state.WhiteToMove
}
