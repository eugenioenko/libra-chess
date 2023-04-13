package libra

const (
	MoveQuiet = iota
	MoveCapture
	MoveOnPassant
	MovePromotion
	MovePromotionCapture
	MoveCastle
)

type Move struct {
	From     byte
	To       byte
	MoveType byte
	Data     [2]byte
}

func NewMove(from byte, to byte, moveType byte, data [2]byte) *Move {
	return &Move{
		From:     from,
		To:       to,
		MoveType: moveType,
		Data:     data,
	}
}

type BoardMoves struct {
	All        []*Move
	Quite      []*Move
	Captures   []*Move
	Promotions []*Move
}

func NewBoardMoves() *BoardMoves {
	return &BoardMoves{
		All:        []*Move{},
		Quite:      []*Move{},
		Captures:   []*Move{},
		Promotions: []*Move{},
	}
}

// N, S, E, W, NE, SE, SW, SE

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

var SquaresToEdge [64][8]byte = generateSquaresToEdge()
var BoardDirOffsets [8]int8 = [8]int8{-8, 1, 8, -1, -7, 9, 7, -9}
var SquareKnightJumps [64][8]byte = generateKnightJumps()

func (board *Board) GeneratePawnMoves() {
	squares := board.Pieces.White.Pawns
	if !board.WhiteToMove {
		squares = board.Pieces.Black.Pawns
	}

	for _, square := range squares {
		// two squares forward
		var leftRankSquare byte = SquareA2
		if !board.WhiteToMove {
			leftRankSquare = SquareA7
		}
		if square >= leftRankSquare && square <= leftRankSquare+8 {
			var amountToMove int8 = 16
			if !board.WhiteToMove {
				amountToMove = -16
			}
			squareToMove := square - byte(amountToMove)
			if board.IsSquareEmpty(squareToMove) {
				board.AddMove(NewMove(square, squareToMove, MoveQuiet, [2]byte{0, 0}))
			}
		}

		// one move forward and promotion
		var amountToMove int8 = 8
		var dirToMove int8 = 1
		if !board.WhiteToMove {
			dirToMove = -1
		}
		squareToMove := square - byte(amountToMove*dirToMove)
		if board.IsSquareEmpty(squareToMove) {
			if (board.WhiteToMove && board.IsSquareAt8thRank(squareToMove)) || (!board.WhiteToMove && board.IsSquareAt1stRank(squareToMove)) {
				board.AddMove(NewMove(square, squareToMove, MovePromotion, [2]byte{0, 0}))
			} else {
				board.AddMove(NewMove(square, squareToMove, MoveQuiet, [2]byte{0, 0}))
			}
		}

		// captures
		if board.WhiteToMove {
			leftSquare := square - 8 - 1
			if !board.IsSquareAtHFile(square) && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareBlack(leftSquare) {
				if board.IsSquareAt8thRank(leftSquare) {
					board.AddPromotion(NewMove(square, leftSquare, MovePromotionCapture, [2]byte{0, 0}))
				} else {
					board.AddCapture(NewMove(square, leftSquare, MoveCapture, [2]byte{0, 0}))
				}
			}
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(leftSquare) && board.IsPieceAtSquareBlack(leftSquare+8) {
				board.AddCapture(NewMove(square, leftSquare, MoveOnPassant, [2]byte{0, 0}))
			}
			rightSquare := square - 8 + 1
			if !board.IsSquareAtAFile(square) && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareBlack(rightSquare) {
				if board.IsSquareAt1stRank(leftSquare) {
					board.AddPromotion(NewMove(square, rightSquare, MovePromotionCapture, [2]byte{0, 0}))
				} else {
					board.AddCapture(NewMove(square, rightSquare, MoveCapture, [2]byte{0, 0}))
				}
			}
			if !board.IsSquareAtAFile(square) && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare+8) {
				board.AddCapture(NewMove(square, rightSquare, MoveOnPassant, [2]byte{0, 0}))
			}
		} else {
			rightSquare := square + 8 - 1
			if !board.IsSquareAtHFile(square) && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareBlack(rightSquare) {
				if board.IsSquareAt1stRank(rightSquare) {
					board.AddCapture(NewMove(square, rightSquare, MovePromotionCapture, [2]byte{0, 0}))
				} else {
					board.AddCapture(NewMove(square, rightSquare, MoveCapture, [2]byte{0, 0}))
				}
			}
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare-8) {
				board.AddCapture(NewMove(square, rightSquare, MoveOnPassant, [2]byte{0, 0}))
			}
			leftSquare := square + 8 + 1
			if !board.IsSquareAtAFile(square) && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareBlack(leftSquare) {
				if board.IsSquareAt1stRank(rightSquare) {
					board.AddCapture(NewMove(square, leftSquare, MovePromotionCapture, [2]byte{0, 0}))
				} else {
					board.AddCapture(NewMove(square, leftSquare, MoveCapture, [2]byte{0, 0}))
				}
			}
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare-8) {
				board.AddCapture(NewMove(square, rightSquare, MoveOnPassant, [2]byte{0, 0}))
			}
		}
	}
}

func (board *Board) GenerateSlidingMoves(pieces []byte, startDir byte, endDir byte) {
	for _, square := range pieces {
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			amountToMove := int8(SquaresToEdge[square][dirOffset])
			for moveIndex := int8(1); moveIndex <= amountToMove; moveIndex++ {
				squareTo := int8(square) + (offset * moveIndex)
				isQuiteMove := board.AddQuiteOrCapture(square, byte(squareTo))
				if !isQuiteMove {
					break
				}
			}
		}
	}
}

func (board *Board) GenerateKingMoves() {
	square := board.Pieces.White.King
	if !board.WhiteToMove {
		square = board.Pieces.Black.King
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[square][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(square) + offset
			board.AddQuiteOrCapture(square, byte(squareTo))
		}
	}
}

func (board *Board) GenerateRookMoves() {
	rooks := board.Pieces.White.Rooks
	if !board.WhiteToMove {
		rooks = board.Pieces.Black.Rooks
	}
	board.GenerateSlidingMoves(rooks, 0, 4)
}

func (board *Board) GenerateBishopMoves() {
	bishops := board.Pieces.White.Bishops
	if !board.WhiteToMove {
		bishops = board.Pieces.Black.Bishops
	}
	board.GenerateSlidingMoves(bishops, 4, 8)
}

func (board *Board) GenerateQueenMoves() {
	queues := board.Pieces.White.Queens
	if !board.WhiteToMove {
		queues = board.Pieces.Black.Queens
	}
	board.GenerateSlidingMoves(queues, 0, 8)
}

func (board *Board) GenerateKnightMoves() {
	knights := board.Pieces.White.Knights
	if !board.WhiteToMove {
		knights = board.Pieces.Black.Knights
	}
	for _, square := range knights {
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				board.AddQuiteOrCapture(square, squareTo)
			}
		}
	}
}
