package libra

import (
	"fmt"
	"os"
	"strings"
)

// Piece Types
const (
	SquareA2 = 48
	SquareA7 = 8
)

const BoardInitialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type KingLocations struct {
	White byte
	Black byte
}

type CastlingAvailability struct {
	BlackKingSide  bool
	BlackQueenSide bool
	WhiteKingSide  bool
	WhiteQueenSide bool
}

type Board struct {
	// 0: h1, 63: a8
	Position             [64]byte
	AttackedSquares      []bool
	CastlingAvailability CastlingAvailability
	Pieces               PieceColorLocation
	WhiteToMove          bool
	Moves                []Move
	OnPassant            byte
	HalfMoveClock        int
	FullMoveCounter      int
	Played               string
}

func NewBoard() *Board {
	board := &Board{}
	return board
}

func (board *Board) Initialize(position [64]byte) {
	board.Position = position
	board.CastlingAvailability = CastlingAvailability{
		BlackKingSide:  true,
		BlackQueenSide: true,
		WhiteKingSide:  true,
		WhiteQueenSide: true,
	}
	board.Pieces = NewPieceColorLocation()
	board.WhiteToMove = true
	board.Moves = []Move{}
	board.OnPassant = 0
	board.HalfMoveClock = 0
	board.FullMoveCounter = 0
}

func (board *Board) Clone() *Board {
	clone := NewBoard()
	clone.Position = board.Position
	clone.CastlingAvailability = CastlingAvailability{
		BlackKingSide:  true,
		BlackQueenSide: true,
		WhiteKingSide:  true,
		WhiteQueenSide: true,
	}
	clone.Pieces = NewPieceColorLocation()
	clone.WhiteToMove = board.WhiteToMove
	clone.Moves = []Move{}
	clone.OnPassant = board.OnPassant
	clone.HalfMoveClock = board.HalfMoveClock
	clone.FullMoveCounter = board.FullMoveCounter
	clone.Played = board.Played
	return clone
}

// https://en.wikipedia.org/wiki/Forsyth–Edwards_Notation
func (board *Board) LoadInitial() {
	board.Initialize([64]byte{})
	board.LoadFromFEN(BoardInitialFEN)
}

func (board *Board) LoadFromFEN(fen string) (bool, error) {
	board.Initialize([64]byte{})
	parts := strings.Split(fen, " ")
	if len(parts) == 0 {
		return false, fmt.Errorf("invalid FEN, missing blocks, at least piece list block required")
	}

	// 1. Piece placement data
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return false, fmt.Errorf("invalid FEN, missing ranks")
	}

	// TODO validate characters in FEN code are valid

	index := 0
	for _, pieces := range ranks {
		for _, piece := range pieces {
			if CharIsNumber(piece) {
				emptyCount := int(piece - '0')
				board.removePieces(index, emptyCount)
				index += emptyCount
			} else {
				board.Position[index] = byte(piece)
				index += 1
			}
		}
	}
	if index != 64 {
		return false, fmt.Errorf("invalid FEN, missing pieces")
	}

	// 2. Active Color
	if len(parts) > 1 {
		board.WhiteToMove = parts[1] == "w"
	}

	// 3. Castling
	if len(parts) > 2 {
		if parts[2] == "-" {
			// TODO castling
			board.CastlingAvailability = CastlingAvailability{
				BlackKingSide:  false,
				BlackQueenSide: false,
				WhiteKingSide:  false,
				WhiteQueenSide: false,
			}
		} else {
			board.CastlingAvailability = CastlingAvailability{
				BlackKingSide:  true,
				BlackQueenSide: true,
				WhiteKingSide:  true,
				WhiteQueenSide: true,
			}
		}
	}

	// 4. On-passant
	if len(parts) > 3 && parts[3] != "-" {
		onPassant, ok := SquareNameToIndex(parts[3])
		if ok {
			board.OnPassant = onPassant
		}
	}

	// generate piece list table
	board.UpdatePiecesLocation()
	return true, nil
}

func (board *Board) removePieces(start int, count int) (bool, error) {
	if start+count > 64 {
		return false, fmt.Errorf("invalid remove pieces range, out of range")
	}
	for index := 0; index < count; index++ {
		board.Position[start+index] = 0
	}
	return true, nil
}

func (board *Board) Print() {
	fmt.Println()
	for index, piece := range board.Position {

		if index%8 == 0 {
			fmt.Print(8 - index/8)
			fmt.Print(" | ")
		}
		if piece != 0 {
			fmt.Print(PieceCodeToFont(piece))
		} else {
			fmt.Print(" ")
		}
		fmt.Print(" ")
		if index > 0 && ((index+1)%8) == 0 {
			fmt.Print("\n")
		}
	}
	fmt.Print("   ----------------\n    A B C D E F G H\n\n")
}

func (board *Board) IsSquareValid(square byte) bool {
	return square <= 63
}

func (board *Board) IsSquareEmpty(square byte) bool {
	return board.IsSquareValid(square) && board.Position[square] == 0
}

func (board *Board) IsSquareEmptyAndNotAttacked(square byte) bool {
	return board.IsSquareEmpty(square) && !board.AttackedSquares[square]
}

func (board *Board) IsSquareOccupied(square byte) bool {
	return board.IsSquareValid(square) && board.Position[square] > 0
}

func (board *Board) IsSquarePawn(square byte) bool {
	return board.Position[square] == WhitePawn || board.Position[square] == BlackPawn
}

func (board *Board) IsSquareOnPassant(square byte) bool {
	return board.OnPassant == square
}

func (board *Board) IsPieceAtSquareBlack(square byte) bool {
	return board.Position[square] >= 98
}

func (board *Board) IsPieceAtSquareWhite(square byte) bool {
	return board.Position[square] > 0 && board.Position[square] < 98
}

func (board *Board) SquareToRank(square byte) byte {
	return 8 - square/8
}

func (board *Board) IsSquareAtAFile(square byte) bool {
	return square%8 == 0
}

func (board *Board) IsSquareAtHFile(square byte) bool {
	return (square+1)%8 == 0
}

func (board *Board) AddQuietOrCapture(from, to byte, whiteToMove bool) bool {
	if board.IsSquareEmpty(to) {
		board.AddQuietMove(from, to)
		return true
	} else {
		if (whiteToMove && board.IsPieceAtSquareBlack(to)) || (!whiteToMove && board.IsPieceAtSquareWhite(to)) {
			board.AddCapture(from, to, MoveCapture, whiteToMove)
		}
		return false
	}
}

func (board *Board) AddMove(move Move) {
	board.Moves = append(board.Moves, move)
}

func (board *Board) AddQuietMove(from, to byte) {
	move := NewMove(from, to, MoveQuiet, [2]byte{0, 0})
	board.Moves = append(board.Moves, move)
}

func (board *Board) AddCastleMove(from, to byte) {
	move := NewMove(from, to, MoveCastle, [2]byte{0, 0})
	board.Moves = append(board.Moves, move)
}

func (board *Board) AddCapture(from, to, moveType byte, whiteToMove bool) {
	captured := board.Position[to]
	if moveType == MoveEnPassant {
		if whiteToMove {
			captured = BlackPawn
		} else {
			captured = WhitePawn
		}
	}
	move := NewMove(from, to, moveType, [2]byte{captured, 0})
	board.AddMove(move)
}

func (board *Board) AddPromotion(from, to, captured byte, whiteToMove bool) {
	var moveType byte = MovePromotion
	if captured != 0 {
		moveType = MovePromotionCapture
	}
	var queenPiece byte = WhiteQueen
	var rookPiece byte = WhiteRook
	var bishopPiece byte = WhiteBishop
	var knightPiece byte = WhiteKnight

	if !whiteToMove {
		queenPiece = BlackQueen
		rookPiece = BlackRook
		bishopPiece = BlackBishop
		knightPiece = BlackKnight
	}

	move := NewMove(from, to, moveType, [2]byte{queenPiece, captured})
	board.AddMove(move)

	move = NewMove(from, to, moveType, [2]byte{rookPiece, captured})
	board.AddMove(move)

	move = NewMove(from, to, moveType, [2]byte{bishopPiece, captured})
	board.AddMove(move)

	move = NewMove(from, to, moveType, [2]byte{knightPiece, captured})
	board.AddMove(move)
}

func (board *Board) CountMoves() *MovesCount {
	count := NewMovesCount()
	count.All = len(board.Moves)
	for _, move := range board.Moves {
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

func (board *Board) GeneratePawnMoves(whiteToMove bool) {
	squares := board.Pieces.White.Pawns
	if !whiteToMove {
		squares = board.Pieces.Black.Pawns
	}

	for _, square := range squares {

		// two space forward, pawns are on initial square
		if whiteToMove {
			if board.SquareToRank(square) == 2 {
				squareToMove := square - 16
				if board.IsSquareEmpty(squareToMove) {
					board.AddQuietMove(square, squareToMove)
				}
			}
		} else {
			if board.SquareToRank(square) == 7 {
				squareToMove := square + 16
				if board.IsSquareEmpty(squareToMove) {
					board.AddQuietMove(square, squareToMove)
				}
			}
		}

		// one move forward and promotion
		if whiteToMove {
			squareToMove := square - 8
			if board.IsSquareEmpty(squareToMove) {
				if board.SquareToRank(squareToMove) == 8 {
					board.AddPromotion(square, squareToMove, 0, whiteToMove)
				} else {
					board.AddQuietMove(square, squareToMove)
				}
			}
		} else {
			squareToMove := square + 8
			if board.IsSquareEmpty(squareToMove) {
				if board.SquareToRank(squareToMove) == 1 {
					board.AddPromotion(square, squareToMove, 0, whiteToMove)
				} else {
					board.AddQuietMove(square, squareToMove)
				}
			}
		}

		// captures and promotion captures
		if whiteToMove {
			// left capture with white
			leftSquare := square - 8 - 1
			if !board.IsSquareAtHFile(square) && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareBlack(leftSquare) {
				if board.SquareToRank(leftSquare) == 8 {
					// promotion capture
					captured := board.Position[leftSquare]
					board.AddPromotion(square, leftSquare, captured, whiteToMove)
				} else {
					// normal capture
					board.AddCapture(square, leftSquare, MoveCapture, whiteToMove)
				}
			}
			// en-passant capture left
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(leftSquare) && board.IsPieceAtSquareBlack(leftSquare+8) {
				board.AddCapture(square, leftSquare, MoveEnPassant, whiteToMove)
			}
			// right capture with white
			rightSquare := square - 8 + 1
			if !board.IsSquareAtAFile(square) && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareBlack(rightSquare) {
				if board.SquareToRank(rightSquare) == 1 {
					// promotion capture
					captured := board.Position[rightSquare]
					board.AddPromotion(square, rightSquare, captured, whiteToMove)
				} else {
					// normal capture
					board.AddCapture(square, rightSquare, MoveCapture, whiteToMove)
				}
			}
			// en-passant capture right
			if !board.IsSquareAtAFile(square) && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareBlack(rightSquare+8) {
				board.AddCapture(square, rightSquare, MoveEnPassant, whiteToMove)
			}
		} else {
			// right capture with black
			rightSquare := square + 8 - 1
			if !board.IsSquareAtHFile(square) && board.IsSquareOccupied(rightSquare) && board.IsPieceAtSquareWhite(rightSquare) {
				if board.SquareToRank(rightSquare) == 1 {
					captured := board.Position[rightSquare]
					board.AddPromotion(square, rightSquare, captured, whiteToMove)
				} else {
					board.AddCapture(square, rightSquare, MoveCapture, whiteToMove)
				}
			}
			// en-passant capture right
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(rightSquare) && board.IsPieceAtSquareWhite(rightSquare-8) {
				board.AddCapture(square, rightSquare, MoveEnPassant, whiteToMove)
			}
			// left capture with black
			leftSquare := square + 8 + 1
			if !board.IsSquareAtAFile(square) && board.IsSquareOccupied(leftSquare) && board.IsPieceAtSquareWhite(leftSquare) {
				if board.SquareToRank(leftSquare) == 1 {
					captured := board.Position[leftSquare]
					board.AddPromotion(square, leftSquare, captured, whiteToMove)
				} else {
					board.AddCapture(square, leftSquare, MoveCapture, whiteToMove)
				}
			}
			// en passant capture right
			if !board.IsSquareAtHFile(square) && board.IsSquareOnPassant(leftSquare) && board.IsPieceAtSquareWhite(leftSquare-8) {
				board.AddCapture(square, leftSquare, MoveEnPassant, whiteToMove)
			}
		}
	}
}

func (board *Board) GeneratePawnAttackingSquares(whiteToMove bool) {
	squares := board.Pieces.White.Pawns
	if !whiteToMove {
		squares = board.Pieces.Black.Pawns
	}
	for _, square := range squares {
		// captures and promotion captures
		if whiteToMove {
			// left attack with white
			leftSquare := square - 8 - 1
			if !board.IsSquareAtHFile(square) {
				board.AddCapture(square, leftSquare, MoveCapture, whiteToMove)
			}
			// right attack with white
			rightSquare := square - 8 + 1
			if !board.IsSquareAtAFile(square) {
				board.AddCapture(square, rightSquare, MoveCapture, whiteToMove)
			}
		} else {
			// right attack with black
			rightSquare := square + 8 - 1
			if !board.IsSquareAtHFile(square) {
				board.AddCapture(square, rightSquare, MoveCapture, whiteToMove)
			}
			// left attack with black
			leftSquare := square + 8 + 1
			if !board.IsSquareAtAFile(square) {
				board.AddCapture(square, leftSquare, MoveCapture, whiteToMove)
			}
		}
	}
}

func (board *Board) GenerateSlidingMoves(pieces []byte, startDir byte, endDir byte, whiteToMove bool) {
	for _, square := range pieces {
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			amountToMove := int8(SquaresToEdge[square][dirOffset])
			for moveIndex := int8(1); moveIndex <= amountToMove; moveIndex++ {
				squareTo := int8(square) + (offset * moveIndex)
				isQuietMove := board.AddQuietOrCapture(square, byte(squareTo), whiteToMove)
				if !isQuietMove {
					break
				}
			}
		}
	}
}

func (board *Board) GenerateKingMoves(whiteToMove bool) {
	square := board.Pieces.White.King
	if !whiteToMove {
		square = board.Pieces.Black.King
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[square][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(square) + offset
			board.AddQuietOrCapture(square, byte(squareTo), whiteToMove)
		}
	}
}

func (board *Board) GenerateCastleMoves(whiteToMove bool) {
	if whiteToMove {
		if (board.CastlingAvailability.WhiteQueenSide) && board.IsSquareEmptyAndNotAttacked(57) && board.IsSquareEmptyAndNotAttacked(58) && board.IsSquareEmptyAndNotAttacked(59) {
			board.AddCastleMove(60, 58)
		} else if (board.CastlingAvailability.WhiteKingSide) && board.IsSquareEmptyAndNotAttacked(61) && board.IsSquareEmptyAndNotAttacked(62) {
			board.AddCastleMove(60, 62)
		}
	} else {
		if (board.CastlingAvailability.BlackQueenSide) && board.IsSquareEmptyAndNotAttacked(1) && board.IsSquareEmptyAndNotAttacked(2) && board.IsSquareEmptyAndNotAttacked(3) {
			board.AddCastleMove(4, 2)
		} else if (board.CastlingAvailability.BlackKingSide) && board.IsSquareEmptyAndNotAttacked(5) && board.IsSquareEmptyAndNotAttacked(6) {
			board.AddCastleMove(4, 6)
		}
	}
}

func (board *Board) GenerateRookMoves(whiteToMove bool) {
	rooks := board.Pieces.White.Rooks
	if !whiteToMove {
		rooks = board.Pieces.Black.Rooks
	}
	board.GenerateSlidingMoves(rooks, 0, 4, whiteToMove)
}

func (board *Board) GenerateBishopMoves(whiteToMove bool) {
	bishops := board.Pieces.White.Bishops
	if !whiteToMove {
		bishops = board.Pieces.Black.Bishops
	}
	board.GenerateSlidingMoves(bishops, 4, 8, whiteToMove)
}

func (board *Board) GenerateQueenMoves(whiteToMove bool) {
	queues := board.Pieces.White.Queens
	if !whiteToMove {
		queues = board.Pieces.Black.Queens
	}
	board.GenerateSlidingMoves(queues, 0, 8, whiteToMove)
}

func (board *Board) GenerateKnightMoves(whiteToMove bool) {
	knights := board.Pieces.White.Knights
	if !whiteToMove {
		knights = board.Pieces.Black.Knights
	}
	for _, square := range knights {
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				board.AddQuietOrCapture(square, squareTo, whiteToMove)
			}
		}
	}
}

func (board *Board) MakeMove(move Move) {
	// update counters
	board.FullMoveCounter += 1
	if move.MoveType == MoveCapture || board.IsSquarePawn(move.From) {
		board.HalfMoveClock = 0
	} else {
		board.HalfMoveClock += 1
	}
	if !board.WhiteToMove {
		board.FullMoveCounter += 1
	}

	// update position
	piece := board.Position[move.From]
	board.Position[move.To] = piece
	board.Position[move.From] = 0

	// Remove EnPassant pawn
	if move.MoveType == MoveEnPassant {
		if board.WhiteToMove {
			board.Position[move.To+8] = 0
		} else {
			board.Position[move.To-8] = 0
		}
	}

	// Update pieces location
	// TODO avoid calling this after every move
	board.UpdatePiecesLocation()

	// TODO Handle castling update rook position

	// TODO Update the en passant square if a pawn moved two squares

	// TODO Update castling rights

	// switch active color
	board.Played += PieceCodeToNotation[piece] + BoardSquareNames[move.To] + " "
	board.WhiteToMove = !board.WhiteToMove
}

func (board *Board) GeneratePseudoLegalMoves() {
	// generate moving color all moves
	board.GenerateAttackedSquares(!board.WhiteToMove)
	board.Moves = []Move{}
	board.GeneratePawnMoves(board.WhiteToMove)
	board.GenerateKnightMoves(board.WhiteToMove)
	board.GenerateBishopMoves(board.WhiteToMove)
	board.GenerateRookMoves(board.WhiteToMove)
	board.GenerateQueenMoves(board.WhiteToMove)
	board.GenerateKingMoves(board.WhiteToMove)
	board.GenerateCastleMoves(board.WhiteToMove)
}

func (board *Board) GenerateLegalMoves() {
	legalMoves := []Move{}
	board.GeneratePseudoLegalMoves()
	for _, move := range board.Moves {
		if board.IsMoveLegal(move) {
			legalMoves = append(legalMoves, move)
		}
	}
	board.Moves = legalMoves
}

func (board *Board) IsMoveLegal(move Move) bool {
	tmpBoard := board.Clone()
	tmpBoard.MakeMove(move)
	return !tmpBoard.IsKingInCheck(board.WhiteToMove)
}

func (board *Board) IsKingInCheck(whiteToMove bool) bool {
	board.GenerateAttackedSquares(!whiteToMove)
	king := board.Pieces.White.King
	if !whiteToMove {
		king = board.Pieces.Black.King
	}
	return board.AttackedSquares[king]
}

func (board *Board) Perft(depth int) int {
	if depth == 0 {
		return 1
	}
	board.GenerateLegalMoves()

	if depth == 1 {
		return len(board.Moves)
	}

	var count int
	for _, move := range board.Moves {
		tmpBoard := board.Clone()
		tmpBoard.MakeMove(move)
		count += tmpBoard.Perft(depth - 1)
	}

	return count
}

func (board *Board) PerftMoves(depth int, f *os.File) int {
	if depth == 0 {
		return 1
	}
	board.GenerateLegalMoves()

	if depth == 1 {
		fmt.Fprintln(f, board.Played)
		return len(board.Moves)
	}

	var count int
	for _, move := range board.Moves {
		tmpBoard := board.Clone()
		tmpBoard.MakeMove(move)
		count += tmpBoard.PerftMoves(depth-1, f)
	}

	return count
}

func (board *Board) GenerateAttackedSquares(whiteToMove bool) {
	board.Moves = []Move{}
	board.AttackedSquares = make([]bool, 64)

	// generate opposing color attacked squares
	board.GenerateQueenMoves(whiteToMove)
	board.GenerateKnightMoves(whiteToMove)
	board.GenerateBishopMoves(whiteToMove)
	board.GenerateRookMoves(whiteToMove)
	board.GenerateKingMoves(whiteToMove)
	board.GeneratePawnAttackingSquares(whiteToMove)
	for _, move := range board.Moves {
		if move.MoveType == MoveCapture || move.MoveType == MoveQuiet || move.MoveType == MovePromotionCapture {
			board.AttackedSquares[move.To] = true
		}
	}
}

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
