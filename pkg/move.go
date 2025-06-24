package libra

const (
	MoveQuiet = iota
	MoveCapture
	MoveEnPassant
	MoveCastle
	MovePromotion
	MovePromotionCapture
)

type Move struct {
	Piece    byte
	From     byte
	To       byte
	MoveType byte
	Promoted byte
	Captured byte
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

func NewMove(piece, from, to, moveType, promotion, capture byte) Move {
	return Move{
		Piece:    piece,
		From:     from,
		To:       to,
		MoveType: moveType,
		Promoted: promotion,
		Captured: capture,
	}
}

type MoveState struct {
	Hash                 uint64
	WhitePawns           uint64
	WhiteKnights         uint64
	WhiteBishops         uint64
	WhiteRooks           uint64
	WhiteQueens          uint64
	WhiteKing            uint64
	BlackPawns           uint64
	BlackKnights         uint64
	BlackBishops         uint64
	BlackRooks           uint64
	BlackQueens          uint64
	BlackKing            uint64
	CastlingAvailability CastlingAvailability
	OnPassant            byte
	HalfMoveClock        int
	FullMoveCounter      int
	WhiteToMove          bool
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
		switch move.Promoted {
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

// Move returns a MoveState for undoing the move
func (board *Board) Move(move Move) MoveState {
	// Save current state
	prev := MoveState{
		WhitePawns:           board.WhitePawns,
		WhiteKnights:         board.WhiteKnights,
		WhiteBishops:         board.WhiteBishops,
		WhiteRooks:           board.WhiteRooks,
		WhiteQueens:          board.WhiteQueens,
		WhiteKing:            board.WhiteKing,
		BlackPawns:           board.BlackPawns,
		BlackKnights:         board.BlackKnights,
		BlackBishops:         board.BlackBishops,
		BlackRooks:           board.BlackRooks,
		BlackQueens:          board.BlackQueens,
		BlackKing:            board.BlackKing,
		CastlingAvailability: board.CastlingAvailability,
		OnPassant:            board.OnPassant,
		HalfMoveClock:        board.HalfMoveClock,
		FullMoveCounter:      board.FullMoveCounter,
		WhiteToMove:          board.WhiteToMove,
		Hash:                 board.Hash,
	}

	if !board.WhiteToMove {
		board.FullMoveCounter += 1
	}
	from := move.From
	to := move.To
	piece := move.Piece

	if move.MoveType == MoveCapture || board.IsSquarePawn(from) {
		board.HalfMoveClock = 0
	} else {
		board.HalfMoveClock += 1
	}

	// Remove piece from 'from' square
	board.clearPieceAtSquare(from, piece)
	// Place piece at 'to' square
	if move.MoveType == MovePromotion || move.MoveType == MovePromotionCapture {
		board.setPieceAtSquare(to, move.Promoted)
	} else {
		board.setPieceAtSquare(to, piece)
	}

	// Incremental Zobrist hash update
	// Remove piece from 'from' square
	board.Hash ^= zobristPieceTable[from][piece]
	// Place piece at 'to' square
	if move.MoveType == MovePromotion || move.MoveType == MovePromotionCapture {
		board.Hash ^= zobristPieceTable[to][move.Promoted]
	} else {
		board.Hash ^= zobristPieceTable[to][piece]
	}

	// Handle en passant
	if move.MoveType == MoveEnPassant {
		if board.WhiteToMove {
			board.clearPieceAtSquare(to+8, BlackPawn)
			board.Hash ^= zobristPieceTable[to+8][BlackPawn]
		} else {
			board.clearPieceAtSquare(to-8, WhitePawn)
			board.Hash ^= zobristPieceTable[to-8][WhitePawn]
		}
	}

	// Handle castling
	if move.MoveType == MoveCastle {
		if piece == WhiteKing {
			if to == SquareG1 {
				board.clearPieceAtSquare(SquareH1, WhiteRook)
				board.setPieceAtSquare(SquareF1, WhiteRook)
				board.Hash ^= zobristPieceTable[SquareH1][WhiteRook]
				board.Hash ^= zobristPieceTable[SquareF1][WhiteRook]
			} else if to == SquareC1 {
				board.clearPieceAtSquare(SquareA1, WhiteRook)
				board.setPieceAtSquare(SquareD1, WhiteRook)
				board.Hash ^= zobristPieceTable[SquareA1][WhiteRook]
				board.Hash ^= zobristPieceTable[SquareD1][WhiteRook]
			}
			board.CastlingAvailability.WhiteKingSide = false
			board.CastlingAvailability.WhiteQueenSide = false
		} else if piece == BlackKing {
			if to == SquareG8 {
				board.clearPieceAtSquare(SquareH8, BlackRook)
				board.setPieceAtSquare(SquareF8, BlackRook)
				board.Hash ^= zobristPieceTable[SquareH8][BlackRook]
				board.Hash ^= zobristPieceTable[SquareF8][BlackRook]
			} else if to == SquareC8 {
				board.clearPieceAtSquare(SquareA8, BlackRook)
				board.setPieceAtSquare(SquareD8, BlackRook)
				board.Hash ^= zobristPieceTable[SquareA8][BlackRook]
				board.Hash ^= zobristPieceTable[SquareD8][BlackRook]
			}
			board.CastlingAvailability.BlackKingSide = false
			board.CastlingAvailability.BlackQueenSide = false
		}
	}
	// Update castling rights for moving king or rook
	if move.MoveType != MoveCastle {
		if piece == WhiteKing && board.CastlingAvailability.WhiteKingSide {
			board.CastlingAvailability.WhiteKingSide = false
			board.CastlingAvailability.WhiteQueenSide = false
		} else if piece == BlackKing && board.CastlingAvailability.BlackKingSide {
			board.CastlingAvailability.BlackKingSide = false
			board.CastlingAvailability.BlackQueenSide = false
		}

		if piece == WhiteRook {
			if from == SquareA1 && board.CastlingAvailability.WhiteQueenSide {
				board.CastlingAvailability.WhiteQueenSide = false
			} else if from == SquareH1 && board.CastlingAvailability.WhiteKingSide {
				board.CastlingAvailability.WhiteKingSide = false
			}
		} else if piece == BlackRook {
			if from == SquareA8 && board.CastlingAvailability.BlackQueenSide {
				board.CastlingAvailability.BlackQueenSide = false
			} else if from == SquareH8 && board.CastlingAvailability.BlackKingSide {
				board.CastlingAvailability.BlackKingSide = false
			}
		}
	}

	// Update castling rights if a rook is captured
	if move.MoveType == MoveCapture || move.MoveType == MovePromotionCapture {
		if move.Captured == WhiteRook {
			if to == SquareA1 && board.CastlingAvailability.WhiteQueenSide {
				board.CastlingAvailability.WhiteQueenSide = false
			} else if to == SquareH1 && board.CastlingAvailability.WhiteKingSide {
				board.CastlingAvailability.WhiteKingSide = false
			}
		} else if move.Captured == BlackRook {
			if to == SquareA8 && board.CastlingAvailability.BlackQueenSide {
				board.CastlingAvailability.BlackQueenSide = false
			} else if to == SquareH8 && board.CastlingAvailability.BlackKingSide {
				board.CastlingAvailability.BlackKingSide = false
			}
		}
	}

	// Update castling rights hash
	if board.CastlingAvailability.WhiteKingSide != prev.CastlingAvailability.WhiteKingSide {
		board.Hash ^= zobristCastlingAvailability.WhiteKingSide
	}
	if board.CastlingAvailability.WhiteQueenSide != prev.CastlingAvailability.WhiteQueenSide {
		board.Hash ^= zobristCastlingAvailability.WhiteQueenSide
	}
	if board.CastlingAvailability.BlackKingSide != prev.CastlingAvailability.BlackKingSide {
		board.Hash ^= zobristCastlingAvailability.BlackKingSide
	}
	if board.CastlingAvailability.BlackQueenSide != prev.CastlingAvailability.BlackQueenSide {
		board.Hash ^= zobristCastlingAvailability.BlackQueenSide
	}

	// Remove captured piece
	if move.MoveType == MoveCapture || move.MoveType == MovePromotionCapture {
		capturedPieceSquare := move.To
		capturedPiece := move.Captured
		board.clearPieceAtSquare(capturedPieceSquare, capturedPiece)
		board.Hash ^= zobristPieceTable[capturedPieceSquare][capturedPiece]
		// Update castling rights if a rook is captured
		if capturedPiece == WhiteRook {
			if capturedPieceSquare == SquareA1 {
				board.CastlingAvailability.WhiteQueenSide = false
			} else if capturedPieceSquare == SquareH1 {
				board.CastlingAvailability.WhiteKingSide = false
			}
		} else if capturedPiece == BlackRook {
			if capturedPieceSquare == SquareA8 {
				board.CastlingAvailability.BlackQueenSide = false
			} else if capturedPieceSquare == SquareH8 {
				board.CastlingAvailability.BlackKingSide = false
			}
		}
	}

	// Update en passant in hash
	if board.OnPassant != 0 {
		board.Hash ^= zobristOnPassantTable[board.OnPassant]
	}
	// Set new en passant
	board.OnPassant = 0
	if piece == WhitePawn && from/8 == 6 && to/8 == 4 {
		board.OnPassant = from - 8
	} else if piece == BlackPawn && from/8 == 1 && to/8 == 3 {
		board.OnPassant = from + 8
	}
	if board.OnPassant != 0 {
		board.Hash ^= zobristOnPassantTable[board.OnPassant]
	}

	// Update side to move
	if board.WhiteToMove {
		board.Hash ^= zobristWhiteToMove
	} else {
		board.Hash ^= zobristBlackToMove
	}
	board.WhiteToMove = !board.WhiteToMove
	if board.WhiteToMove {
		board.Hash ^= zobristWhiteToMove
	} else {
		board.Hash ^= zobristBlackToMove
	}

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
	board.CastlingAvailability = state.CastlingAvailability
	board.OnPassant = state.OnPassant
	board.HalfMoveClock = state.HalfMoveClock
	board.FullMoveCounter = state.FullMoveCounter
	board.WhiteToMove = state.WhiteToMove
	board.Hash = state.Hash
}
