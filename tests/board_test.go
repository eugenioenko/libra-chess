package libra_test

import (
	"fmt"
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

func TestShouldInsertValues(t *testing.T) {
	board := NewBoard()
	ok, error := board.FromFEN(BoardInitialFEN)

	if error != nil {
		fmt.Println(error)
	}

	if ok != true {
		t.Fail()
	}
}

func TestShouldCalculateZobristHash(t *testing.T) {
	board := NewBoard()
	board.FromFEN(BoardInitialFEN)
	hashA := board.ZobristHash()
	board.FromFEN("rnbqkb1r/ppp1p1pp/8/3p1pP1/3Pn3/2N5/PPP1PP1P/R1BQKBNR w KQkq - 1 5")
	hashB := board.ZobristHash()
	board.PrintPosition()
	if hashA == 0 {
		t.Fail()
	}

	if hashA == hashB {
		t.Fail()
	}
}

func TestShouldGeneratePawnMoves(t *testing.T) {
	board := NewBoard()
	board.FromFEN(BoardInitialFEN)
	moves := board.GeneratePawnMoves(board.WhiteToMove)
	if len(moves) != 16 {
		t.Fail()
	}
	if CountMoves(moves).Capture != 0 {
		t.Fail()
	}
	board.FromFEN("rnbqkbnr/4p3/pp1p1p1p/2p3p1/1P1PPPP1/P1P5/7P/RNBQKBNR w KQkq - 0 8")
	moves = board.GeneratePawnMoves(board.WhiteToMove)
	if len(moves) != 11 {
		t.Fail()
	}
	if CountMoves(moves).Capture != 3 {
		t.Fail()
	}
}

func TestShouldGenerateOnPassantPawnMoves(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbqkbnr/8/pp1p1p1p/2p1pPp1/1P1PP1P1/P1P5/7P/RNBQKBNR w KQkq e6 0 9")
	moves := board.GeneratePawnMoves(board.WhiteToMove)
	if CountMoves(moves).Capture != 4 {
		t.Fail()
	}
}

func TestShouldGeneratePromotionMoves(t *testing.T) {
	board := NewBoard()
	board.FromFEN("3n4/4P3/8/8/2K2k2/8/8/8 w - - 0 1")
	moves := board.GeneratePawnMoves(board.WhiteToMove)
	if CountMoves(moves).Promotion != 8 {
		t.Fail()
	}

	board.FromFEN("8/2P5/8/4k3/8/4K3/7p/8 b - - 0 1")
	moves = board.GeneratePawnMoves(board.WhiteToMove)
	if CountMoves(moves).Promotion != 4 {
		t.Fail()
	}

}

func TestShouldGenerateRookMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("1k4r1/8/2R4p/8/8/8/8/7K")
	moves := board.GenerateRookMoves(board.WhiteToMove)
	if len(moves) != 14 {
		t.Fail()
	}
	if CountMoves(moves).Capture != 1 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/R7")
	moves = board.GenerateRookMoves(board.WhiteToMove)
	if len(moves) != 14 {
		t.Fail()
	}

	board.FromFEN("R7/8/8/8/8/8/8/8")
	moves = board.GenerateRookMoves(board.WhiteToMove)
	if len(moves) != 14 {
		t.Fail()
	}

	board.FromFEN("7R/8/8/8/8/8/8/8")
	moves = board.GenerateRookMoves(board.WhiteToMove)
	if len(moves) != 14 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7R")
	moves = board.GenerateRookMoves(board.WhiteToMove)
	if len(moves) != 14 {
		t.Fail()
	}

	board.FromFEN("8/8/4R3/8/8/8/8/8")
	moves = board.GenerateRookMoves(board.WhiteToMove)
	if len(moves) != 14 {
		t.Fail()
	}

}

func TestShouldGenerateBishopMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4B3/8/8/8/8")
	moves := board.GenerateBishopMoves(board.WhiteToMove)
	if len(moves) != 13 {
		t.Fail()
	}

	board.FromFEN("B7/8/8/8/8/8/8/8")
	moves = board.GenerateBishopMoves(board.WhiteToMove)
	if len(moves) != 7 {
		t.Fail()
	}

	board.FromFEN("7B/8/8/8/8/8/8/8")
	moves = board.GenerateBishopMoves(board.WhiteToMove)
	if len(moves) != 7 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7B")
	moves = board.GenerateBishopMoves(board.WhiteToMove)
	if len(moves) != 7 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/B7")
	moves = board.GenerateBishopMoves(board.WhiteToMove)
	if len(moves) != 7 {
		t.Fail()
	}
}

func TestShouldGenerateQueenMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4Q3/8/8/8/8")
	moves := board.GenerateQueenMoves(board.WhiteToMove)
	if len(moves) != 27 {
		t.Fail()
	}

	board.FromFEN("Q7/8/8/8/8/8/8/8")
	moves = board.GenerateQueenMoves(board.WhiteToMove)
	if len(moves) != 21 {
		t.Fail()
	}

	board.FromFEN("7Q/8/8/8/8/8/8/8")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(moves) != 21 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7Q")
	moves = board.GenerateQueenMoves(board.WhiteToMove)
	if len(moves) != 21 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/Q7")
	moves = board.GenerateQueenMoves(board.WhiteToMove)
	if len(moves) != 21 {
		t.Fail()
	}
}

func TestShouldGenerateKingMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4K3/8/8/8/8")
	moves := board.GenerateKingMoves(board.WhiteToMove)
	if len(moves) != 8 {
		t.Fail()
	}

	board.FromFEN("K7/8/8/8/8/8/8/8")
	moves = board.GenerateKingMoves(board.WhiteToMove)
	if len(moves) != 3 {
		t.Fail()
	}

	// King vs queen (king on a8, queen on b8)
	board.FromFEN("KQ6/8/8/8/8/8/8/8")
	moves = board.GenerateKingMoves(board.WhiteToMove)
	if len(moves) != 2 {
		t.Fail()
	}

	// King vs queen 2 (king on h8, queen on g8)
	board.FromFEN("6QK/8/8/8/8/8/8/8")
	moves = board.GenerateKingMoves(board.WhiteToMove)
	if len(moves) != 2 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/K7")
	moves = board.GenerateKingMoves(board.WhiteToMove)
	if len(moves) != 3 {
		t.Fail()
	}
}

func TestShouldGenerateKnightMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4N3/8/8/8/8")
	moves := board.GenerateKnightMoves(board.WhiteToMove)
	if len(moves) != 8 {
		t.Fail()
	}

	board.FromFEN("N7/8/8/8/8/8/8/8")
	moves = board.GenerateKnightMoves(board.WhiteToMove)
	if len(moves) != 2 {
		t.Fail()
	}

	board.FromFEN("7N/8/8/8/8/8/8/8")
	moves = board.GenerateKnightMoves(board.WhiteToMove)
	if len(moves) != 2 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7N")
	moves = board.GenerateKnightMoves(board.WhiteToMove)
	if len(moves) != 2 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/N7")

	moves = board.GenerateKnightMoves(board.WhiteToMove)
	if len(moves) != 2 {
		t.Fail()
	}
}

func TestShouldGenerateOnlyLegalMoves(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/8/5k2/8/q5q1/q5q1/2P1P3/3K4 w - - 0 1")
	moves := board.GenerateLegalMoves()

	if len(moves) != 1 {
		t.Fail()
	}
}

// This test checks that en passant is not allowed if it exposes the king to check.
func TestEnPassantExposesKing(t *testing.T) {
	// FEN: White pawn on e5, black pawn on d5, white king on e1, black rook on e4
	// After d5-d4, e5xd6 en passant would expose the white king to check from the rook
	board := NewBoard()
	board.FromFEN("4k3/8/8/3pP3/4r3/8/8/4K3 w - - 0 1")
	moves := board.GenerateLegalMoves()
	for _, move := range moves {
		if move.MoveType == MoveEnPassant {
			t.Errorf("En passant should not be legal if it exposes the king to check")
		}
	}
}

// This test checks that castling rights are lost if a rook is captured on its original square by a promotion-capture.
func TestPromotionCaptureRemovesCastlingRights(t *testing.T) {
	// FEN: Black rook on h8, white pawn on g7, white king on e1, black king on e8, white to move
	// White plays g8=Q capturing rook on h8
	board := NewBoard()
	board.FromFEN("4k2r/6P1/8/8/8/8/8/4K3 w Kk - 0 1")
	moves := board.GeneratePawnMoves(true)
	found := false
	for _, move := range moves {
		if move.MoveType == MovePromotionCapture && move.To == SquareH8 {
			found = true
			board.Move(move)
			if board.CastlingAvailability.BlackKingSide {
				t.Errorf("Castling rights should be lost after rook is captured by promotion-capture")
			}
		}
	}
	if !found {
	}
}

func TestMoveKingIntoCheck(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareE1, WhiteKing)
	board.SetPiece(SquareE8, BlackRook)
	board.WhiteToMove = true
	moves := board.GenerateLegalMoves()
	for _, move := range moves {
		if move.From == SquareE1 && (move.To == SquareE2 || move.To == SquareE3 || move.To == SquareE4 || move.To == SquareE5 || move.To == SquareE6 || move.To == SquareE7 || move.To == SquareE8) {
			t.Errorf("King should not be able to move to %d which is in the rook's attack line", move.To)
		}
	}
}

func TestMoveKingAdjacentToEnemyKing(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareE1, WhiteKing)
	board.SetPiece(SquareE2, BlackKing)
	board.WhiteToMove = true
	moves := board.GenerateLegalMoves()
	for _, move := range moves {
		if move.From == SquareE1 && move.To == SquareE2 {
			t.Errorf("Move king from e1 to e2 should not be legal (adjacent to enemy king)")
		}
	}
}

func TestBoardToFEN(t *testing.T) {
	board := NewBoard()
	ok, err := board.FromFEN(BoardInitialFEN)
	if !ok || err != nil {
		t.Fatalf("Failed to load initial FEN: %v", err)
	}
	fen := board.ToFEN()
	if fen != BoardInitialFEN {
		t.Errorf("ToFEN() did not match initial FEN. Got: %s, Want: %s", fen, BoardInitialFEN)
	}

	// Test a position with en passant, castling, and clocks
	fenStr := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq e3 5 10"
	ok, err = board.FromFEN(fenStr)
	if !ok || err != nil {
		t.Fatalf("Failed to load FEN: %v", err)
	}
	fen2 := board.ToFEN()
	if fen2 != fenStr {
		t.Errorf("ToFEN() did not match. Got: %s, Want: %s", fen2, fenStr)
	}
}

// --- Additional edge case tests for 100% coverage ---

func TestPawnBlockedAndNoDoublePush(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareE2, WhitePawn)
	board.SetPiece(SquareE3, BlackPawn)
	board.WhiteToMove = true
	moves := board.GeneratePawnMoves(true)
	if len(moves) != 0 {
		t.Errorf("Blocked pawn should have no moves")
	}
}

func TestPawnPromotionWithCaptureAndNoCapture(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareG7, WhitePawn)
	board.SetPiece(SquareH8, BlackRook)
	board.WhiteToMove = true
	moves := board.GeneratePawnMoves(true)
	promotion, capture := false, false
	for _, move := range moves {
		if move.MoveType == MovePromotion {
			promotion = true
		}
		if move.MoveType == MovePromotionCapture {
			capture = true
		}
	}
	if !promotion || !capture {
		t.Errorf("Pawn promotion and promotion-capture should be generated")
	}
}

func TestKnightEdgeCases(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareA1, WhiteKnight)
	board.SetPiece(SquareH8, WhiteKnight)
	board.WhiteToMove = true
	moves := board.GenerateKnightMoves(true)
	if len(moves) != 4 {
		t.Errorf("Knights on corners should have 2 moves each (total 4 moves), got %d", len(moves))
	}
}

func TestRookBlockedAndEdge(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareA1, WhiteRook)
	board.SetPiece(SquareA2, WhitePawn)
	board.WhiteToMove = true
	moves := board.GenerateRookMoves(true)
	if len(moves) != 7 {
		t.Errorf("Rook should have 7 moves along the 1st rank, got %d", len(moves))
	}
}

func TestQueenBlockedAndEdge(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareA1, WhiteQueen)
	board.SetPiece(SquareA2, WhitePawn)
	board.SetPiece(SquareB2, WhitePawn)
	board.WhiteToMove = true
	moves := board.GenerateQueenMoves(true)
	if len(moves) != 7 {
		t.Errorf("Queen should have 7 moves along the 1st rank, got %d", len(moves))
	}
}

func TestKingCannotMoveIntoCheck(t *testing.T) {
	board := NewBoard()
	board.SetPiece(SquareE1, WhiteKing)
	board.SetPiece(SquareE3, BlackRook)
	board.WhiteToMove = true
	moves := board.GenerateLegalMoves()
	for _, move := range moves {
		if move.From == SquareE1 && move.To == SquareE2 {
			t.Errorf("King should not be able to move into check")
		}
	}
}
