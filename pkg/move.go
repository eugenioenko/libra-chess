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
	From     byte
	To       byte
	MoveType byte
	Data     [2]byte
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

func NewMove(from byte, to byte, moveType byte, data [2]byte) Move {
	return Move{
		From:     from,
		To:       to,
		MoveType: moveType,
		Data:     data,
	}
}

func generateSquaresToEdge() [64][8]byte {
	squares := [64][8]byte{}
	for i := range squares {
		index := byte(i)
		y := index / 8
		x := index - y*8
		south := 7 - y
		north := y
		west := x
		east := 7 - x
		squares[index][0] = north
		squares[index][1] = east
		squares[index][2] = south
		squares[index][3] = west
		squares[index][4] = MathMinByte(north, east)
		squares[index][5] = MathMinByte(south, east)
		squares[index][6] = MathMinByte(south, west)
		squares[index][7] = MathMinByte(north, west)
	}
	return squares
}

func generateKnightJumps() [64][8]byte {
	squares := [64][8]byte{}
	jumpOffsets := [8][2]int8{{1, 2}, {-1, 2}, {2, -1}, {-2, -1}, {-1, -2}, {1, -2}, {-2, 1}, {2, 1}}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			squareFrom := y*8 + x
			for offsetIndex, offset := range jumpOffsets {
				x2 := int8(x) + offset[0]
				y2 := int8(y) + offset[1]
				if x2 >= 0 && y2 >= 0 && x2 < 8 && y2 < 8 {
					squares[squareFrom][offsetIndex] = byte(y2*8 + x2)
				} else {
					squares[squareFrom][offsetIndex] = 255
				}
			}
		}
	}
	return squares
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
		switch move.Data[0] {
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

var SquaresToEdge [64][8]byte = generateSquaresToEdge()
var SquareKnightJumps [64][8]byte = generateKnightJumps()
var BoardDirOffsets [8]int8 = [8]int8{-8, 1, 8, -1, -7, 9, 7, -9}

type MoveState struct {
	Position             [64]byte
	CastlingAvailability CastlingAvailability
	OnPassant            byte
	HalfMoveClock        int
	FullMoveCounter      int
	WhiteToMove          bool
}

// MakeMove now returns a MoveState for undoing the move
func (board *Board) MakeMove(move Move) MoveState {
	// Save current state
	prev := MoveState{
		Position:             board.Position,
		CastlingAvailability: board.CastlingAvailability,
		OnPassant:            board.OnPassant,
		HalfMoveClock:        board.HalfMoveClock,
		FullMoveCounter:      board.FullMoveCounter,
		WhiteToMove:          board.WhiteToMove,
	}

	if !board.WhiteToMove {
		board.FullMoveCounter += 1
	}
	if move.MoveType == MoveCapture || board.IsSquarePawn(move.From) {
		board.HalfMoveClock = 0
	} else {
		board.HalfMoveClock += 1
	}

	piece := board.Position[move.From]

	if move.MoveType == MoveCapture || move.MoveType == MovePromotionCapture || move.MoveType == MoveEnPassant {
		captured := board.getCapturedPiece(move.MoveType, move.To, board.WhiteToMove)
		if captured == WhiteKing || captured == BlackKing {
			panic("Attempted to capture a king!")
		}
	}
	board.Position[move.To] = piece
	board.Position[move.From] = 0

	if move.MoveType == MovePromotion || move.MoveType == MovePromotionCapture {
		board.Position[move.To] = move.Data[0]
	}

	if move.MoveType == MoveEnPassant {
		if board.WhiteToMove {
			board.Position[move.To+8] = 0
		} else {
			board.Position[move.To-8] = 0
		}
	}

	if move.MoveType == MoveCastle {
		if piece == WhiteKing {
			if move.To == SquareG1 {
				board.Position[SquareH1] = 0
				board.Position[SquareF1] = WhiteRook
			} else if move.To == SquareC1 {
				board.Position[SquareA1] = 0
				board.Position[SquareD1] = WhiteRook
			}
			board.CastlingAvailability.WhiteKingSide = false
			board.CastlingAvailability.WhiteQueenSide = false
		} else if piece == BlackKing {
			if move.To == SquareG8 {
				board.Position[SquareH8] = 0
				board.Position[SquareF8] = BlackRook
			} else if move.To == SquareC8 {
				board.Position[SquareA8] = 0
				board.Position[SquareD8] = BlackRook
			}
			board.CastlingAvailability.BlackKingSide = false
			board.CastlingAvailability.BlackQueenSide = false
		}
	}

	board.OnPassant = 0
	if piece == WhitePawn && move.From/8 == 6 && move.To/8 == 4 {
		board.OnPassant = move.From - 8
	} else if piece == BlackPawn && move.From/8 == 1 && move.To/8 == 3 {
		board.OnPassant = move.From + 8
	}

	if piece == WhiteKing {
		board.CastlingAvailability.WhiteKingSide = false
		board.CastlingAvailability.WhiteQueenSide = false
	}
	if piece == BlackKing {
		board.CastlingAvailability.BlackKingSide = false
		board.CastlingAvailability.BlackQueenSide = false
	}
	if piece == WhiteRook {
		if move.From == SquareH1 {
			board.CastlingAvailability.WhiteKingSide = false
		}
		if move.From == SquareA1 {
			board.CastlingAvailability.WhiteQueenSide = false
		}
	}
	if piece == BlackRook {
		if move.From == SquareH8 {
			board.CastlingAvailability.BlackKingSide = false
		}
		if move.From == SquareA8 {
			board.CastlingAvailability.BlackQueenSide = false
		}
	}

	if move.MoveType == MoveCapture || move.MoveType == MovePromotionCapture {
		capturedPieceSquare := move.To
		var capturedPiece byte
		if move.MoveType == MovePromotionCapture {
			capturedPiece = move.Data[1]
		} else {
			capturedPiece = move.Data[0]
		}

		if capturedPiece == WhiteRook {
			if capturedPieceSquare == SquareH1 {
				board.CastlingAvailability.WhiteKingSide = false
			}
			if capturedPieceSquare == SquareA1 {
				board.CastlingAvailability.WhiteQueenSide = false
			}
		}
		if capturedPiece == BlackRook {
			if capturedPieceSquare == SquareH8 {
				board.CastlingAvailability.BlackKingSide = false
			}
			if capturedPieceSquare == SquareA8 {
				board.CastlingAvailability.BlackQueenSide = false
			}
		}
	}

	board.UpdatePiecesLocation()

	board.WhiteToMove = !board.WhiteToMove

	return prev
}

// UndoMove restores the board to a previous MoveState
func (board *Board) UndoMove(state MoveState) {
	board.Position = state.Position
	board.CastlingAvailability = state.CastlingAvailability
	board.OnPassant = state.OnPassant
	board.HalfMoveClock = state.HalfMoveClock
	board.FullMoveCounter = state.FullMoveCounter
	board.WhiteToMove = state.WhiteToMove
	board.UpdatePiecesLocation()
}
