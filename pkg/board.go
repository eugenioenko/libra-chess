package libra

import (
	"fmt"
	"strconv"
	"strings"
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
func (board *Board) ParseAndApplyPosition(positionArgs []string) {
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
	board.FromFEN(fen)
	// Play moves if any
	for i := movesStart; i < len(positionArgs); i++ {
		if positionArgs[i] == "moves" {
			for _, moveStr := range positionArgs[i+1:] {
				move := board.ParseUCIMove(moveStr)
				if move != nil {
					board.MakeMove(*move)
				}
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

func (board *Board) PrintMove(move Move) {
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
	fmt.Printf("%s %s -> %s [%s] data: %v\n", pieceStr, fromName, toName, moveType, move.Data)
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

// IsSquareAttacked returns true if the square is under attack
func (board *Board) IsSquareAttacked(square byte) bool {
	return board.AttackedSquares[square]
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

// OnlyKingLeft returns true if the given color has only the king left on the board.
// Pass true for white, false for black.
func (board *Board) OnlyKingLeft() bool {
	if board.WhiteToMove {
		return len(board.Pieces.Black.Pawns) == 0 &&
			len(board.Pieces.Black.Knights) == 0 &&
			len(board.Pieces.Black.Bishops) == 0 &&
			len(board.Pieces.Black.Rooks) == 0 &&
			len(board.Pieces.Black.Queens) == 0 &&
			board.Pieces.Black.King != 0
	} else {
		return len(board.Pieces.White.Pawns) == 0 &&
			len(board.Pieces.White.Knights) == 0 &&
			len(board.Pieces.White.Bishops) == 0 &&
			len(board.Pieces.White.Rooks) == 0 &&
			len(board.Pieces.White.Queens) == 0 &&
			board.Pieces.White.King != 0
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

func (board *Board) Clone() *Board {
	clone := NewBoard()
	// Deep copy Position
	copy(clone.Position[:], board.Position[:])
	clone.CastlingAvailability = board.CastlingAvailability
	clone.Pieces = board.Pieces
	clone.WhiteToMove = board.WhiteToMove
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

// ParseUCIMove parses a move in UCI format (e.g., "e2e4", "e7e8q") and returns a Move struct.
func (board *Board) ParseUCIMove(moveStr string) *Move {
	if len(moveStr) < 4 {
		return nil
	}
	from, ok1 := SquareNameToIndex(moveStr[0:2])
	to, ok2 := SquareNameToIndex(moveStr[2:4])
	if !ok1 || !ok2 {
		return nil
	}
	// Generate all legal moves and find the one matching from/to (and promotion if present)
	moves := board.GenerateLegalMoves()
	for _, move := range moves {
		if move.From == from && move.To == to {
			// Handle promotion
			if len(moveStr) == 5 {
				promo := moveStr[4]
				var promoPiece byte
				if board.WhiteToMove {
					promoPiece = WhitePromotionMap[promo]
				} else {
					promoPiece = BlackPromotionMap[promo]
				}
				if move.Data[0] == promoPiece {
					return &move
				}
			} else if move.MoveType != MovePromotion && move.MoveType != MovePromotionCapture {
				return &move
			}
		}
	}
	return nil // Not found
}
