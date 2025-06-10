package libra

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const BoardInitialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var BoardSquareNames [64]string = [64]string{
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
}

type CastlingAvailability struct {
	BlackKingSide  bool
	BlackQueenSide bool
	WhiteKingSide  bool
	WhiteQueenSide bool
}

// Board represents the state of a chess game, including piece positions, castling rights, en passant, move clocks, and move history.
type Board struct {
	// Position is a 64-byte array representing the board squares (0 = empty, otherwise piece code)
	Position [64]byte
	// AttackedSquares marks which squares are currently attacked by the opponent
	AttackedSquares []bool
	// CastlingAvailability tracks which castling rights are still available
	CastlingAvailability CastlingAvailability
	// Pieces holds the locations of all pieces for both colors
	Pieces PieceColorLocation
	// WhiteToMove is true if it's White's turn, false for Black
	WhiteToMove bool
	// Moves is a list of generated moves (pseudo-legal or legal, depending on context)
	Moves []Move
	// OnPassant is the square index for en passant capture, or 0 if not available
	OnPassant byte
	// HalfMoveClock counts half-moves since the last capture or pawn move (for 50-move rule)
	HalfMoveClock int
	// FullMoveCounter counts the number of full moves (incremented after Black's move)
	FullMoveCounter int
}

// NewBoard creates a new, empty board. You must call LoadInitial or FromFEN to set up a position.
func NewBoard() *Board {
	board := &Board{}
	board.Initialize([64]byte{})
	return board
}

func (board *Board) Reset() {
	for i := range board.Position {
		board.Position[i] = 0
	}
	board.CastlingAvailability = CastlingAvailability{
		BlackKingSide:  false,
		BlackQueenSide: false,
		WhiteKingSide:  false,
		WhiteQueenSide: false,
	}
	board.Pieces = NewPieceColorLocation()
	board.WhiteToMove = true
	board.Moves = []Move{}
	board.OnPassant = 0
	board.HalfMoveClock = 0
	board.FullMoveCounter = 1
	board.AttackedSquares = make([]bool, 64)
}

// Initialize sets up the board with a given position array and resets all state (castling, clocks, etc).
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
	board.FullMoveCounter = 1
	board.AttackedSquares = make([]bool, 64)
}

// LoadInitial sets up the board to the standard chess starting position.
func (board *Board) LoadInitial() {
	board.Reset()
	board.FromFEN(BoardInitialFEN)
}

// FromFEN loads a position from a FEN string. Returns false and error if the FEN is invalid.
func (board *Board) FromFEN(fen string) (bool, error) {
	board.Reset()
	parts := strings.Split(fen, " ")
	if len(parts) == 0 {
		return false, fmt.Errorf("invalid FEN, missing blocks, at least piece list block required")
	}

	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return false, fmt.Errorf("invalid FEN, missing ranks")
	}

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

	if len(parts) > 1 {
		board.WhiteToMove = parts[1] == "w"
	}

	if len(parts) > 2 {
		if parts[2] == "-" {
			board.CastlingAvailability = CastlingAvailability{
				BlackKingSide:  false,
				BlackQueenSide: false,
				WhiteKingSide:  false,
				WhiteQueenSide: false,
			}
		} else {
			castleStr := parts[2]
			board.CastlingAvailability = CastlingAvailability{
				WhiteKingSide:  strings.Contains(castleStr, "K"),
				WhiteQueenSide: strings.Contains(castleStr, "Q"),
				BlackKingSide:  strings.Contains(castleStr, "k"),
				BlackQueenSide: strings.Contains(castleStr, "q"),
			}
		}
	}

	if len(parts) > 3 && parts[3] != "-" {
		onPassant, ok := SquareNameToIndex(parts[3])
		if ok {
			board.OnPassant = onPassant
		}
	}

	if len(parts) > 4 {
		halfMoveVal, err := strconv.Atoi(parts[4])
		if err == nil {
			board.HalfMoveClock = halfMoveVal
		}
	}

	if len(parts) > 5 {
		fullMoveVal, err := strconv.Atoi(parts[5])
		if err == nil {
			board.FullMoveCounter = fullMoveVal
		}
	}

	board.UpdatePiecesLocation()
	return true, nil
}

// ParseAndApplyPosition sets up the board from a UCI "position" command's arguments.
// It supports "startpos" or "fen" and applies any moves listed after "moves".
func (b *Board) ParseAndApplyPosition(positionArgs []string) {
	fen := BoardInitialFEN
	movesStart := 0
	if len(positionArgs) > 0 && positionArgs[0] == "startpos" {
		fen = BoardInitialFEN
		movesStart = 1
	} else if len(positionArgs) > 0 && positionArgs[0] == "fen" {
		fenParts := []string{}
		for i := 1; i < len(positionArgs) && len(fenParts) < 6; i++ {
			fenParts = append(fenParts, positionArgs[i])
			movesStart = i + 1
		}
		fen = strings.Join(fenParts, " ")
	}
	b.FromFEN(fen)
	// Play moves if any
	for i := movesStart; i < len(positionArgs); i++ {
		if positionArgs[i] == "moves" {
			for _, moveStr := range positionArgs[i+1:] {
				move := b.ParseMove(moveStr)
				b.MakeMove(move)
			}
			break
		}
	}
}

// removePieces sets a range of squares to empty (0), used when parsing FEN for empty squares.
func (board *Board) removePieces(start int, count int) (bool, error) {
	if start+count > 64 {
		return false, fmt.Errorf("invalid remove pieces range, out of range")
	}
	for index := 0; index < count; index++ {
		board.Position[start+index] = 0
	}
	return true, nil
}

// PrintPosition prints the board to the console using Unicode chess symbols.
func (board *Board) PrintPosition() {
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

func (board *Board) PrintMoves() {
	for i, move := range board.Moves {
		piece := board.Position[move.From]
		pieceStr := pieceCodeToFont[piece]
		fromName, _ := SquareIndexToName(move.From)
		toName, _ := SquareIndexToName(move.To)
		moveType := ""
		switch move.MoveType {
		case MoveQuiet:
			moveType = "quiet"
		case MoveCapture:
			moveType = "capture"
		case MoveEnPassant:
			moveType = "en passant"
		case MovePromotion:
			moveType = "promotion"
		case MovePromotionCapture:
			moveType = "promotion-capture"
		case MoveCastle:
			moveType = "castle"
		}
		fmt.Printf("%2d: %s %s -> %s [%s] data: %v\n", i+1, pieceStr, fromName, toName, moveType, move.Data)
	}
}

func (board *Board) PrintMove(move Move) {
	for i, move := range board.Moves {
		piece := board.Position[move.From]
		pieceStr := pieceCodeToFont[piece]
		fromName, _ := SquareIndexToName(move.From)
		toName, _ := SquareIndexToName(move.To)
		moveType := ""
		switch move.MoveType {
		case MoveQuiet:
			moveType = "quiet"
		case MoveCapture:
			moveType = "capture"
		case MoveEnPassant:
			moveType = "en passant"
		case MovePromotion:
			moveType = "promotion"
		case MovePromotionCapture:
			moveType = "promotion-capture"
		case MoveCastle:
			moveType = "castle"
		}
		fmt.Printf("%2d: %s %s -> %s [%s] data: %v\n", i+1, pieceStr, fromName, toName, moveType, move.Data)
	}
}

// IsSquareValid returns true if the square index is within 0-63.
func (board *Board) IsSquareValid(square byte) bool {
	return square <= 63
}

// IsSquareEmpty returns true if the square is valid and contains no piece.
func (board *Board) IsSquareEmpty(square byte) bool {
	return board.IsSquareValid(square) && board.Position[square] == 0
}

// IsSquareEmptyAndNotAttacked returns true if the square is empty and not attacked by the opponent.
func (board *Board) IsSquareEmptyAndNotAttacked(square byte) bool {
	return board.IsSquareEmpty(square) && !board.AttackedSquares[square]
}

// IsSquareOccupied returns true if the square is valid and contains a piece.
func (board *Board) IsSquareOccupied(square byte) bool {
	return board.IsSquareValid(square) && board.Position[square] > 0
}

// IsSquarePawn returns true if the square contains a pawn (of either color).
func (board *Board) IsSquarePawn(square byte) bool {
	return board.Position[square] == WhitePawn || board.Position[square] == BlackPawn
}

// IsSquareOnPassant returns true if the square is the current en passant target square.
func (board *Board) IsSquareOnPassant(square byte) bool {
	return board.OnPassant == square
}

// IsPieceAtSquareBlack returns true if the piece at the square is black.
func (board *Board) IsPieceAtSquareBlack(square byte) bool {
	return board.Position[square] >= 98
}

// IsPieceAtSquareWhite returns true if the piece at the square is white.
func (board *Board) IsPieceAtSquareWhite(square byte) bool {
	return board.Position[square] > 0 && board.Position[square] < 98
}

// SquareToRank returns the rank (0-7) of a square index (0 = rank 0, 7 = rank 7)
func (board *Board) SquareToRank(square byte) byte {
	return square / 8
}

// SquareToFile returns the file (0-7) of a square index (0 = file 0, 7 = file 7)
func (board *Board) SquareToFile(square byte) byte {
	return square % 8
}

// AddQuietOrCapture adds a quiet move if the destination is empty, or a capture if occupied by an opponent's piece.
// Returns true if a quiet move was added, false if a capture or blocked.
func (board *Board) AddQuietOrCapture(from, to byte, whiteToMove bool) bool {

	if board.Position[to] == WhiteKing || board.Position[to] == BlackKing {
		return false
	}
	if board.IsSquareEmpty(to) {
		board.AddQuietMove(from, to)
		return true
	} else {

		if (whiteToMove && board.IsPieceAtSquareBlack(to)) || (!whiteToMove && board.IsPieceAtSquareWhite(to)) {

			if board.Position[to] != WhiteKing && board.Position[to] != BlackKing {
				board.AddCapture(from, to, MoveCapture, whiteToMove)
			}
		}
		return false
	}
}

// AddMove appends a move to the board's move list.
func (board *Board) AddMove(move Move) {
	board.Moves = append(board.Moves, move)
}

// AddQuietMove adds a non-capturing move to the move list.
func (board *Board) AddQuietMove(from, to byte) {
	move := NewMove(from, to, MoveQuiet, [2]byte{0, 0})
	board.Moves = append(board.Moves, move)
}

// AddCastleMove adds a castling move to the move list.
func (board *Board) AddCastleMove(from, to byte) {
	move := NewMove(from, to, MoveCastle, [2]byte{0, 0})
	board.Moves = append(board.Moves, move)
}

// getCapturedPiece returns the captured piece for a given move, handling en passant correctly.
func (board *Board) getCapturedPiece(moveType byte, to byte, whiteToMove bool) byte {
	if moveType == MoveEnPassant {
		if whiteToMove {
			return BlackPawn
		} else {
			return WhitePawn
		}
	}
	return board.Position[to]
}

// AddCapture adds a capturing move to the move list. Handles en passant as a special case.
func (board *Board) AddCapture(from, to, moveType byte, whiteToMove bool) {
	captured := board.getCapturedPiece(moveType, to, whiteToMove)
	move := NewMove(from, to, moveType, [2]byte{captured, 0})
	board.AddMove(move)
}

// AddPromotion adds all possible promotion moves (to queen, rook, bishop, knight) for a pawn reaching the last rank.
// If captured != 0, adds promotion-capture moves.
func (board *Board) AddPromotion(from, to, captured byte, whiteToMove bool) {

	promotionPieces := []byte{WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight}
	if !whiteToMove {
		promotionPieces = []byte{BlackQueen, BlackRook, BlackBishop, BlackKnight}
	}
	for _, promo := range promotionPieces {
		moveType := MovePromotion
		if captured != 0 {
			moveType = MovePromotionCapture
		}
		move := NewMove(from, to, byte(moveType), [2]byte{promo, captured})
		board.AddMove(move)
	}
}

// CountMoves returns a summary of the number of moves by type in the current move list.
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

// GeneratePawnMoves generates all pawn moves (including promotions, captures, en passant) for the current side.
func (board *Board) GeneratePawnMoves(whiteToMove bool) {
	var squares []byte
	var dir int8
	var startRank, promotionRank byte
	if whiteToMove {
		squares = board.Pieces.White.Pawns
		dir = -8
		startRank = 6
		promotionRank = 0
	} else {
		squares = board.Pieces.Black.Pawns
		dir = 8
		startRank = 1
		promotionRank = 7
	}
	for _, square := range squares {
		file := board.SquareToFile(square)
		rank := board.SquareToRank(square)

		to := int8(square) + dir
		if to >= 0 && to < 64 && !board.IsSquareOccupied(byte(to)) {
			if byte(to/8) == promotionRank {
				// Always call AddPromotion, no need to check captured (quiet move)
				board.AddPromotion(square, byte(to), 0, whiteToMove)
			} else {
				board.AddQuietMove(square, byte(to))
				if rank == startRank {
					twoForward := int8(square) + 2*dir
					if twoForward >= 0 && twoForward < 64 && !board.IsSquareOccupied(byte(twoForward)) {
						board.AddQuietMove(square, byte(twoForward))
					}
				}
			}
		}

		for _, df := range []int8{-1, 1} {
			captureFile := int8(file) + df
			if captureFile < 0 || captureFile > 7 {
				continue
			}
			captureTo := int8(square) + dir + df
			if captureTo < 0 || captureTo >= 64 {
				continue
			}
			if board.IsSquareOccupied(byte(captureTo)) && board.IsPieceAtSquareWhite(byte(captureTo)) != whiteToMove {
				if byte(captureTo/8) == promotionRank {
					// Always call AddPromotion, pass captured piece
					board.AddPromotion(square, byte(captureTo), board.Position[byte(captureTo)], whiteToMove)
				} else {
					board.AddCapture(square, byte(captureTo), MoveCapture, whiteToMove)
				}
			}

			if board.IsSquareOnPassant(byte(captureTo)) {
				if (whiteToMove && rank == 3) || (!whiteToMove && rank == 4) {
					board.AddCapture(square, byte(captureTo), MoveEnPassant, whiteToMove)
				}
			}
		}
	}
}

// MarkSlidingAttacks marks all squares attacked by sliding pieces (rooks, bishops, queens) in the given directions.
// Used for attack maps and move generation.
func (board *Board) MarkSlidingAttacks(pieces []byte, startDir byte, endDir byte) {
	for _, square := range pieces {
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			amountToMove := int8(SquaresToEdge[square][dirOffset])
			for moveIndex := int8(1); moveIndex <= amountToMove; moveIndex++ {
				squareTo := int8(square) + offset*moveIndex
				if squareTo < 0 || squareTo >= 64 {
					break
				}
				board.AttackedSquares[byte(squareTo)] = true
				if board.IsSquareOccupied(byte(squareTo)) {
					break
				}
			}
		}
	}
}

// GenerateSlidingMoves generates all moves for sliding pieces (rooks, bishops, queens) in the given directions.
func (board *Board) GenerateSlidingMoves(pieces []byte, startDir byte, endDir byte, whiteToMove bool) {
	for _, square := range pieces {
		for dirOffset := startDir; dirOffset < endDir; dirOffset++ {
			offset := BoardDirOffsets[dirOffset]
			amountToMove := int8(SquaresToEdge[square][dirOffset])
			for moveIndex := int8(1); moveIndex <= amountToMove; moveIndex++ {
				squareTo := int8(square) + offset*moveIndex
				if squareTo < 0 || squareTo >= 64 {
					break
				}
				isQuietMove := board.AddQuietOrCapture(square, byte(squareTo), whiteToMove)
				if !isQuietMove {
					break
				}
			}
		}
	}
}

// GenerateKingMoves generates all king moves (excluding castling) for the current side.
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

// GenerateCastleMoves generates castling moves if the king and rook have not moved and the path is clear and not attacked.
// Chess rules: King cannot castle out of, through, or into check; squares between must be empty.
func (board *Board) GenerateCastleMoves(whiteToMove bool) {
	if whiteToMove {
		if board.CastlingAvailability.WhiteQueenSide &&
			board.Position[SquareE1] == WhiteKing &&
			board.Position[SquareA1] == WhiteRook &&
			board.IsSquareEmpty(SquareB1) &&
			board.IsSquareEmptyAndNotAttacked(SquareC1) &&
			board.IsSquareEmptyAndNotAttacked(SquareD1) &&
			!board.AttackedSquares[SquareE1] {
			board.AddCastleMove(SquareE1, SquareC1)
		}

		if board.CastlingAvailability.WhiteKingSide &&
			board.Position[SquareE1] == WhiteKing &&
			board.Position[SquareH1] == WhiteRook &&
			board.IsSquareEmptyAndNotAttacked(SquareF1) &&
			board.IsSquareEmptyAndNotAttacked(SquareG1) &&
			!board.AttackedSquares[SquareE1] {
			board.AddCastleMove(SquareE1, SquareG1)
		}
	} else {
		if board.CastlingAvailability.BlackQueenSide &&
			board.Position[SquareE8] == BlackKing &&
			board.Position[SquareA8] == BlackRook &&
			board.IsSquareEmpty(SquareB8) &&
			board.IsSquareEmptyAndNotAttacked(SquareC8) &&
			board.IsSquareEmptyAndNotAttacked(SquareD8) &&
			!board.AttackedSquares[SquareE8] {
			board.AddCastleMove(SquareE8, SquareC8)
		}

		if board.CastlingAvailability.BlackKingSide &&
			board.Position[SquareE8] == BlackKing &&
			board.Position[SquareH8] == BlackRook &&
			board.IsSquareEmptyAndNotAttacked(SquareF8) &&
			board.IsSquareEmptyAndNotAttacked(SquareG8) &&
			!board.AttackedSquares[SquareE8] {
			board.AddCastleMove(SquareE8, SquareG8)
		}
	}
}

// GenerateRookMoves generates all rook moves for the current side.
func (board *Board) GenerateRookMoves(whiteToMove bool) {
	rooks := board.Pieces.White.Rooks
	if !whiteToMove {
		rooks = board.Pieces.Black.Rooks
	}
	board.GenerateSlidingMoves(rooks, 0, 4, whiteToMove)
}

// GenerateBishopMoves generates all bishop moves for the current side.
func (board *Board) GenerateBishopMoves(whiteToMove bool) {
	bishops := board.Pieces.White.Bishops
	if !whiteToMove {
		bishops = board.Pieces.Black.Bishops
	}
	board.GenerateSlidingMoves(bishops, 4, 8, whiteToMove)
}

// GenerateQueenMoves generates all queen moves for the current side.
func (board *Board) GenerateQueenMoves(whiteToMove bool) {
	queens := board.Pieces.White.Queens
	if !whiteToMove {
		queens = board.Pieces.Black.Queens
	}
	board.GenerateSlidingMoves(queens, 0, 8, whiteToMove)
}

// GenerateKnightMoves generates all knight moves for the current side.
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
			panic("BUG: Attempted to capture a king!")
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

func (board *Board) GeneratePseudoLegalMoves() {

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
	prev := board.MakeMove(move)
	// Generate attacked squares for the opponent after the move
	inCheck := board.IsKingInCheck(!board.WhiteToMove)
	board.UndoMove(prev)
	return !inCheck
}

func (board *Board) IsKingInCheck(whiteToMove bool) bool {

	board.UpdatePiecesLocation()

	board.GenerateAttackedSquares(!whiteToMove)
	king := board.Pieces.White.King
	if !whiteToMove {
		king = board.Pieces.Black.King
	}
	return board.AttackedSquares[king]
}

func (board *Board) GenerateAttackedSquares(whiteToMove bool) {
	board.AttackedSquares = make([]bool, 64)

	if whiteToMove {
		board.MarkSlidingAttacks(board.Pieces.White.Queens, 0, 8)
		board.MarkSlidingAttacks(board.Pieces.White.Bishops, 4, 8)
		board.MarkSlidingAttacks(board.Pieces.White.Rooks, 0, 4)
	} else {
		board.MarkSlidingAttacks(board.Pieces.Black.Queens, 0, 8)
		board.MarkSlidingAttacks(board.Pieces.Black.Bishops, 4, 8)
		board.MarkSlidingAttacks(board.Pieces.Black.Rooks, 0, 4)
	}

	knights := board.Pieces.White.Knights
	if !whiteToMove {
		knights = board.Pieces.Black.Knights
	}
	for _, square := range knights {
		for moveIndex := 0; moveIndex < 8; moveIndex++ {
			squareTo := SquareKnightJumps[square][moveIndex]
			if squareTo < 255 {
				board.AttackedSquares[squareTo] = true
			}
		}
	}

	king := board.Pieces.White.King
	if !whiteToMove {
		king = board.Pieces.Black.King
	}
	for dirOffset := 0; dirOffset < 8; dirOffset++ {
		offset := BoardDirOffsets[dirOffset]
		amountToMove := int8(SquaresToEdge[king][dirOffset])
		if amountToMove > 0 {
			squareTo := int8(king) + offset
			if squareTo >= 0 && squareTo < 64 {
				board.AttackedSquares[byte(squareTo)] = true
			}
		}
	}

	var pawns []byte
	var dir int8
	if whiteToMove {
		pawns = board.Pieces.White.Pawns
		dir = -8
	} else {
		pawns = board.Pieces.Black.Pawns
		dir = 8
	}
	for _, square := range pawns {
		file := board.SquareToFile(square)
		for _, df := range []int8{-1, 1} {
			attackFile := int8(file) + df
			if attackFile < 0 || attackFile > 7 {
				continue
			}
			attack := int8(square) + dir + df
			if attack >= 0 && attack < 64 {
				board.AttackedSquares[byte(attack)] = true
			}
		}
	}
}

// ToFEN returns the FEN (Forsyth-Edwards Notation) string for the current board position.
// FEN encodes the board, turn, castling rights, en passant, halfmove clock, and fullmove number.
func (board *Board) ToFEN() string {
	// Piece placement
	fen := ""
	for rank := 0; rank < 8; rank++ {
		empty := 0
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			piece := board.Position[sq]
			if piece == 0 {
				empty++
			} else {
				if empty > 0 {
					fen += strconv.Itoa(empty)
					empty = 0
				}
				fen += string(piece)
			}
		}
		if empty > 0 {
			fen += strconv.Itoa(empty)
		}
		if rank != 7 {
			fen += "/"
		}
	}

	// Active color
	if board.WhiteToMove {
		fen += " w"
	} else {
		fen += " b"
	}

	// Castling rights
	castle := ""
	if board.CastlingAvailability.WhiteKingSide {
		castle += "K"
	}
	if board.CastlingAvailability.WhiteQueenSide {
		castle += "Q"
	}
	if board.CastlingAvailability.BlackKingSide {
		castle += "k"
	}
	if board.CastlingAvailability.BlackQueenSide {
		castle += "q"
	}
	if castle == "" {
		castle = "-"
	}
	fen += " " + castle

	// En passant
	ep := "-"
	if board.OnPassant > 0 && board.OnPassant < 64 {
		if name, ok := SquareIndexToName(board.OnPassant); ok {
			ep = name
		}
	}
	fen += " " + ep

	// Halfmove clock
	fen += " " + strconv.Itoa(board.HalfMoveClock)
	// Fullmove number
	fen += " " + strconv.Itoa(board.FullMoveCounter)

	return fen
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

func (board *Board) Perft(depth int) int {
	if depth == 0 {
		return 1
	}
	board.GenerateLegalMoves()
	if depth == 1 {
		return len(board.Moves)
	}
	nodes := 0
	for _, move := range board.Moves {
		prev := board.MakeMove(move)
		nodes += board.Perft(depth - 1)
		board.UndoMove(prev)
	}
	return nodes
}

func (board *Board) Clone() *Board {
	clone := NewBoard()
	// Deep copy Position
	copy(clone.Position[:], board.Position[:])
	clone.CastlingAvailability = board.CastlingAvailability
	clone.Pieces = board.Pieces
	clone.WhiteToMove = board.WhiteToMove
	clone.Moves = []Move{}
	clone.OnPassant = board.OnPassant
	clone.HalfMoveClock = board.HalfMoveClock
	clone.FullMoveCounter = board.FullMoveCounter
	// Deep copy AttackedSquares
	if board.AttackedSquares != nil {
		clone.AttackedSquares = make([]bool, len(board.AttackedSquares))
		copy(clone.AttackedSquares, board.AttackedSquares)
	} else {
		clone.AttackedSquares = make([]bool, 64) // Ensure it's initialized if source was nil
	}
	return clone
}

// PerftParallel parallelizes the top-level perft search using goroutines for each root move.
func (board *Board) PerftParallel(depth int) int {
	if depth == 0 {
		return 1
	}
	board.GenerateLegalMoves()
	if depth == 1 {
		return len(board.Moves)
	}

	var wg sync.WaitGroup
	results := make(chan int, len(board.Moves))

	for _, move := range board.Moves {
		wg.Add(1)
		go func(m Move) {
			defer wg.Done()
			child := board.Clone()
			child.MakeMove(m)
			count := child.Perft(depth - 1)
			results <- count
		}(move)
	}

	wg.Wait()
	close(results)

	nodes := 0
	for n := range results {
		nodes += n
	}
	return nodes
}

// ParseMove parses a move in UCI format (e.g., "e2e4", "e7e8q") and returns a Move struct.
func (board *Board) ParseMove(moveStr string) Move {
	if len(moveStr) < 4 {
		return Move{}
	}
	from, ok1 := SquareNameToIndex(moveStr[0:2])
	to, ok2 := SquareNameToIndex(moveStr[2:4])
	if !ok1 || !ok2 {
		return Move{}
	}
	// Generate all legal moves and find the one matching from/to (and promotion if present)
	board.GenerateLegalMoves()
	for _, move := range board.Moves {
		if move.From == from && move.To == to {
			// Handle promotion
			if len(moveStr) == 5 {
				promo := moveStr[4]
				promoPiece := byte(0)
				if board.WhiteToMove {
					switch promo {
					case 'q':
						promoPiece = WhiteQueen
					case 'r':
						promoPiece = WhiteRook
					case 'b':
						promoPiece = WhiteBishop
					case 'n':
						promoPiece = WhiteKnight
					}
				} else {
					switch promo {
					case 'q':
						promoPiece = BlackQueen
					case 'r':
						promoPiece = BlackRook
					case 'b':
						promoPiece = BlackBishop
					case 'n':
						promoPiece = BlackKnight
					}
				}
				if move.Data[0] == promoPiece {
					return move
				}
			} else if move.MoveType != MovePromotion && move.MoveType != MovePromotionCapture {
				return move
			}
		}
	}
	return Move{} // Not found
}
