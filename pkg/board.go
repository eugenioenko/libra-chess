package libra

import (
	"fmt"
	"strconv"
	"strings"
)

type CastlingAvailability struct {
	BlackKingSide  bool
	BlackQueenSide bool
	WhiteKingSide  bool
	WhiteQueenSide bool
}

// Board represents the state of a chess game, including piece positions, castling rights, en passant, move clocks, and move history.
type Board struct {
	// Bitboards for each piece type and color
	WhitePawns   uint64
	WhiteKnights uint64
	WhiteBishops uint64
	WhiteRooks   uint64
	WhiteQueens  uint64
	WhiteKing    uint64
	BlackPawns   uint64
	BlackKnights uint64
	BlackBishops uint64
	BlackRooks   uint64
	BlackQueens  uint64
	BlackKing    uint64
	// AttackedSquares is a bitboard marking which squares are currently attacked by the opponent
	AttackedSquares uint64
	// CastlingAvailability tracks which castling rights are still available
	CastlingAvailability CastlingAvailability
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
	board := &Board{
		WhitePawns:      0,
		WhiteKnights:    0,
		WhiteBishops:    0,
		WhiteRooks:      0,
		WhiteQueens:     0,
		WhiteKing:       0,
		BlackPawns:      0,
		BlackKnights:    0,
		BlackBishops:    0,
		BlackRooks:      0,
		BlackQueens:     0,
		BlackKing:       0,
		AttackedSquares: 0,
		CastlingAvailability: CastlingAvailability{
			BlackKingSide:  true,
			BlackQueenSide: true,
			WhiteKingSide:  true,
			WhiteQueenSide: true,
		},
		WhiteToMove:     true,
		OnPassant:       0,
		HalfMoveClock:   0,
		FullMoveCounter: 1,
	}
	return board
}

func (board *Board) Reset() {
	board.WhitePawns = 0
	board.WhiteKnights = 0
	board.WhiteBishops = 0
	board.WhiteRooks = 0
	board.WhiteQueens = 0
	board.WhiteKing = 0
	board.BlackPawns = 0
	board.BlackKnights = 0
	board.BlackBishops = 0
	board.BlackRooks = 0
	board.BlackQueens = 0
	board.BlackKing = 0
	board.AttackedSquares = 0
	board.CastlingAvailability = CastlingAvailability{
		BlackKingSide:  false,
		BlackQueenSide: false,
		WhiteKingSide:  false,
		WhiteQueenSide: false,
	}
	board.WhiteToMove = true
	board.OnPassant = 0
	board.HalfMoveClock = 0
	board.FullMoveCounter = 1
}

// LoadInitial sets up the board to the standard chess starting position.
func (board *Board) LoadInitial() {
	board.Reset()
	board.FromFEN(BoardInitialFEN)
}

// FromFEN loads a position from a FEN string. Returns false and error if the FEN is invalid.
func (board *Board) FromFEN(fen string) (bool, error) {
	// Reset all bitboards
	board.WhitePawns, board.WhiteKnights, board.WhiteBishops, board.WhiteRooks, board.WhiteQueens, board.WhiteKing = 0, 0, 0, 0, 0, 0
	board.BlackPawns, board.BlackKnights, board.BlackBishops, board.BlackRooks, board.BlackQueens, board.BlackKing = 0, 0, 0, 0, 0, 0

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
				index += emptyCount
			} else {
				bb := uint64(1) << index
				switch piece {
				case 'P':
					board.WhitePawns |= bb
				case 'N':
					board.WhiteKnights |= bb
				case 'B':
					board.WhiteBishops |= bb
				case 'R':
					board.WhiteRooks |= bb
				case 'Q':
					board.WhiteQueens |= bb
				case 'K':
					board.WhiteKing |= bb
				case 'p':
					board.BlackPawns |= bb
				case 'n':
					board.BlackKnights |= bb
				case 'b':
					board.BlackBishops |= bb
				case 'r':
					board.BlackRooks |= bb
				case 'q':
					board.BlackQueens |= bb
				case 'k':
					board.BlackKing |= bb
				}
				index++
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
					board.Move(*move)
				}
			}
			break
		}
	}
}

// PrintPosition prints the board to the console using Unicode chess symbols.
func (board *Board) PrintPosition() {
	fmt.Println()
	for index := 0; index < 64; index++ {
		if index%8 == 0 {
			fmt.Print(8 - index/8)
			fmt.Print(" | ")
		}
		piece := board.pieceAtSquare(byte(index))
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

// pieceAtSquare returns the piece code at a given square, or 0 if empty.
func (board *Board) pieceAtSquare(square byte) byte {
	mask := uint64(1) << square
	switch {
	case (board.WhitePawns & mask) != 0:
		return WhitePawn
	case (board.WhiteKnights & mask) != 0:
		return WhiteKnight
	case (board.WhiteBishops & mask) != 0:
		return WhiteBishop
	case (board.WhiteRooks & mask) != 0:
		return WhiteRook
	case (board.WhiteQueens & mask) != 0:
		return WhiteQueen
	case (board.WhiteKing & mask) != 0:
		return WhiteKing
	case (board.BlackPawns & mask) != 0:
		return BlackPawn
	case (board.BlackKnights & mask) != 0:
		return BlackKnight
	case (board.BlackBishops & mask) != 0:
		return BlackBishop
	case (board.BlackRooks & mask) != 0:
		return BlackRook
	case (board.BlackQueens & mask) != 0:
		return BlackQueen
	case (board.BlackKing & mask) != 0:
		return BlackKing
	default:
		return 0
	}
}

// PrintMove prints a move using bitboards.
func (board *Board) PrintMove(move Move) {
	piece := board.pieceAtSquare(move.From)
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
	if !board.IsSquareValid(square) {
		return false
	}
	mask := uint64(1) << square
	occupied := board.WhitePawns | board.WhiteKnights | board.WhiteBishops | board.WhiteRooks | board.WhiteQueens | board.WhiteKing |
		board.BlackPawns | board.BlackKnights | board.BlackBishops | board.BlackRooks | board.BlackQueens | board.BlackKing
	return (occupied & mask) == 0
}

// IsSquareEmptyAndNotAttacked returns true if the square is empty and not attacked by the opponent.
func (board *Board) IsSquareEmptyAndNotAttacked(square byte) bool {
	return board.IsSquareEmpty(square) && (board.AttackedSquares&(uint64(1)<<square) == 0)
}

// IsSquareAttacked returns true if the square is under attack (using bitboard).
func (board *Board) IsSquareAttacked(square byte) bool {
	return (board.AttackedSquares & (uint64(1) << square)) != 0
}

// IsSquareOccupied returns true if the square is valid and contains a piece.
func (board *Board) IsSquareOccupied(square byte) bool {
	if !board.IsSquareValid(square) {
		return false
	}
	mask := uint64(1) << square
	occupied := board.WhitePawns | board.WhiteKnights | board.WhiteBishops | board.WhiteRooks | board.WhiteQueens | board.WhiteKing |
		board.BlackPawns | board.BlackKnights | board.BlackBishops | board.BlackRooks | board.BlackQueens | board.BlackKing
	return (occupied & mask) != 0
}

// IsSquarePawn returns true if the square contains a pawn (of either color).
func (board *Board) IsSquarePawn(square byte) bool {
	mask := uint64(1) << square
	return (board.WhitePawns&mask) != 0 || (board.BlackPawns&mask) != 0
}

// IsSquareOnPassant returns true if the square is the current en passant target square.
func (board *Board) IsSquareOnPassant(square byte) bool {
	return board.OnPassant == square
}

// IsPieceAtSquareBlack returns true if the piece at the square is black.
func (board *Board) IsPieceAtSquareBlack(square byte) bool {
	mask := uint64(1) << square
	return (board.BlackPawns|board.BlackKnights|board.BlackBishops|board.BlackRooks|board.BlackQueens|board.BlackKing)&mask != 0
}

// IsPieceAtSquareWhite returns true if the piece at the square is white.
func (board *Board) IsPieceAtSquareWhite(square byte) bool {
	mask := uint64(1) << square
	return (board.WhitePawns|board.WhiteKnights|board.WhiteBishops|board.WhiteRooks|board.WhiteQueens|board.WhiteKing)&mask != 0
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
		return board.BlackPawns == 0 &&
			board.BlackKnights == 0 &&
			board.BlackBishops == 0 &&
			board.BlackRooks == 0 &&
			board.BlackQueens == 0 &&
			board.BlackKing != 0
	} else {
		return board.WhitePawns == 0 &&
			board.WhiteKnights == 0 &&
			board.WhiteBishops == 0 &&
			board.WhiteRooks == 0 &&
			board.WhiteQueens == 0 &&
			board.WhiteKing != 0
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
			piece := board.pieceAtSquare(byte(sq))
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

// Clone returns a deep copy of the board.
func (board *Board) Clone() *Board {
	clone := NewBoard()
	clone.WhitePawns = board.WhitePawns
	clone.WhiteKnights = board.WhiteKnights
	clone.WhiteBishops = board.WhiteBishops
	clone.WhiteRooks = board.WhiteRooks
	clone.WhiteQueens = board.WhiteQueens
	clone.WhiteKing = board.WhiteKing
	clone.BlackPawns = board.BlackPawns
	clone.BlackKnights = board.BlackKnights
	clone.BlackBishops = board.BlackBishops
	clone.BlackRooks = board.BlackRooks
	clone.BlackQueens = board.BlackQueens
	clone.BlackKing = board.BlackKing
	clone.AttackedSquares = board.AttackedSquares
	clone.CastlingAvailability = board.CastlingAvailability
	clone.WhiteToMove = board.WhiteToMove
	clone.OnPassant = board.OnPassant
	clone.HalfMoveClock = board.HalfMoveClock
	clone.FullMoveCounter = board.FullMoveCounter
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

// SetPiece sets the given piece on the given square, clearing any existing piece at that square.
func (board *Board) SetPiece(square byte, piece byte) {
	// Clear any piece at the square
	mask := ^(uint64(1) << square)
	board.WhitePawns &= mask
	board.WhiteKnights &= mask
	board.WhiteBishops &= mask
	board.WhiteRooks &= mask
	board.WhiteQueens &= mask
	board.WhiteKing &= mask
	board.BlackPawns &= mask
	board.BlackKnights &= mask
	board.BlackBishops &= mask
	board.BlackRooks &= mask
	board.BlackQueens &= mask
	board.BlackKing &= mask
	// Set the new piece
	if piece != 0 {
		bit := uint64(1) << square
		switch piece {
		case WhitePawn:
			board.WhitePawns |= bit
		case WhiteKnight:
			board.WhiteKnights |= bit
		case WhiteBishop:
			board.WhiteBishops |= bit
		case WhiteRook:
			board.WhiteRooks |= bit
		case WhiteQueen:
			board.WhiteQueens |= bit
		case WhiteKing:
			board.WhiteKing |= bit
		case BlackPawn:
			board.BlackPawns |= bit
		case BlackKnight:
			board.BlackKnights |= bit
		case BlackBishop:
			board.BlackBishops |= bit
		case BlackRook:
			board.BlackRooks |= bit
		case BlackQueen:
			board.BlackQueens |= bit
		case BlackKing:
			board.BlackKing |= bit
		}
	}
}
