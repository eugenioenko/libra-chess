package libra_test

import (
	"fmt"
	. "libra/pkg"
	"testing"
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
	hashA := ZobristHash(board)
	board.FromFEN("rnbqkb1r/ppp1p1pp/8/3p1pP1/3Pn3/2N5/PPP1PP1P/R1BQKBNR w KQkq - 1 5")
	hashB := ZobristHash(board)
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
	board.GeneratePawnMoves(board.WhiteToMove)
	if len(board.Moves) != 16 {
		t.Fail()
	}
	if board.CountMoves().Capture != 0 {
		t.Fail()
	}
	board.FromFEN("rnbqkbnr/4p3/pp1p1p1p/2p3p1/1P1PPPP1/P1P5/7P/RNBQKBNR w KQkq - 0 8")
	board.GeneratePawnMoves(board.WhiteToMove)
	if len(board.Moves) != 11 {
		t.Fail()
	}
	if board.CountMoves().Capture != 3 {
		t.Fail()
	}

	/*
		board.FromFEN("rnbqkbnr/ppppppp1/7p/8/P7/8/1PPPPPPP/RNBQKBNR w KQkq - 0 2")
		board.GenerateLegalMoves()
		if len(board.Moves) != 0 {
			t.Fail()
		}


		/*
			board.FromFEN("8/6k1/p7/P5K1/8/8/8/8 b - - 0 1")
			board.GenerateLegalMoves()
			if len(board.Moves) != 0 {
				t.Fail()
			}
	*/

}

func TestShouldGenerateOnPassantPawnMoves(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbqkbnr/8/pp1p1p1p/2p1pPp1/1P1PP1P1/P1P5/7P/RNBQKBNR w KQkq e6 0 9")
	board.GeneratePawnMoves(board.WhiteToMove)
	if board.CountMoves().Capture != 4 {
		t.Fail()
	}
}

func TestShouldGeneratePromotionMoves(t *testing.T) {
	board := NewBoard()
	board.FromFEN("3n4/4P3/8/8/2K2k2/8/8/8 w - - 0 1")
	board.GeneratePawnMoves(board.WhiteToMove)
	if board.CountMoves().Promotion != 8 {
		t.Fail()
	}

	board.FromFEN("8/2P5/8/4k3/8/4K3/7p/8 b - - 0 1")
	board.GeneratePawnMoves(board.WhiteToMove)
	if board.CountMoves().Promotion != 4 {
		t.Fail()
	}

}

func TestShouldGenerateRookMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("1k4r1/8/2R4p/8/8/8/8/7K")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves) != 14 {
		t.Fail()
	}
	if board.CountMoves().Capture != 1 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/R7")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves) != 14 {
		t.Fail()
	}

	board.FromFEN("R7/8/8/8/8/8/8/8")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves) != 14 {
		t.Fail()
	}

	board.FromFEN("7R/8/8/8/8/8/8/8")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves) != 14 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7R")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves) != 14 {
		t.Fail()
	}

	board.FromFEN("8/8/4R3/8/8/8/8/8")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves) != 14 {
		t.Fail()
	}

}

func TestShouldGenerateBishopMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4B3/8/8/8/8")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves) != 13 {
		t.Fail()
	}

	board.FromFEN("B7/8/8/8/8/8/8/8")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves) != 7 {
		t.Fail()
	}

	board.FromFEN("7B/8/8/8/8/8/8/8")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves) != 7 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7B")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves) != 7 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/B7")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves) != 7 {
		t.Fail()
	}
}

func TestShouldGenerateQueenMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4Q3/8/8/8/8")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves) != 27 {
		t.Fail()
	}

	board.FromFEN("Q7/8/8/8/8/8/8/8")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves) != 21 {
		t.Fail()
	}

	board.FromFEN("7Q/8/8/8/8/8/8/8")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves) != 21 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7Q")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves) != 21 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/Q7")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves) != 21 {
		t.Fail()
	}
}

func TestShouldGenerateKingMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4K3/8/8/8/8")
	board.GenerateKingMoves(board.WhiteToMove)
	t.Logf("King moves center: %d", len(board.Moves))
	if len(board.Moves) != 8 {
		t.Fail()
	}

	board.FromFEN("K7/8/8/8/8/8/8/8")
	board.GenerateKingMoves(board.WhiteToMove)
	t.Logf("King moves corner: %d", len(board.Moves))
	if len(board.Moves) != 3 {
		t.Fail()
	}

	// King vs queen (king on a8, queen on b8)
	board.FromFEN("KQ6/8/8/8/8/8/8/8")
	board.GenerateKingMoves(board.WhiteToMove)
	t.Logf("King moves vs queen: %d", len(board.Moves))
	if len(board.Moves) != 2 {
		t.Fail()
	}

	// King vs queen 2 (king on h8, queen on g8)
	board.FromFEN("6QK/8/8/8/8/8/8/8")
	board.GenerateKingMoves(board.WhiteToMove)
	t.Logf("King moves vs queen 2: %d", len(board.Moves))
	if len(board.Moves) != 2 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/K7")
	board.GenerateKingMoves(board.WhiteToMove)
	t.Logf("King moves corner 2: %d", len(board.Moves))
	if len(board.Moves) != 3 {
		t.Fail()
	}
}

func TestShouldGenerateKnightMoves(t *testing.T) {
	board := NewBoard()

	board.FromFEN("8/8/8/4N3/8/8/8/8")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves) != 8 {
		t.Fail()
	}

	board.FromFEN("N7/8/8/8/8/8/8/8")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves) != 2 {
		t.Fail()
	}

	board.FromFEN("7N/8/8/8/8/8/8/8")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves) != 2 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/7N")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves) != 2 {
		t.Fail()
	}

	board.FromFEN("8/8/8/8/8/8/8/N7")

	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves) != 2 {
		t.Fail()
	}
}

func TestShouldGenerateOnlyLegalMoves(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/8/5k2/8/q5q1/q5q1/2P1P3/3K4 w - - 0 1")
	board.GenerateLegalMoves()

	if len(board.Moves) != 1 {
		t.Fail()
	}
}

// This test checks that en passant is not allowed if it exposes the king to check.
func TestEnPassantExposesKing(t *testing.T) {
	// FEN: White pawn on e5, black pawn on d5, white king on e1, black rook on e4
	// After d5-d4, e5xd6 en passant would expose the white king to check from the rook
	board := NewBoard()
	board.FromFEN("4k3/8/8/3pP3/4r3/8/8/4K3 w - - 0 1")
	board.GenerateLegalMoves()
	for _, move := range board.Moves {
		if move.MoveType == MoveEnPassant {
			t.Errorf("En passant should not be legal if it exposes the king to check")
		}
	}
}

func TestDebugBoardCoordinates(t *testing.T) {
	// Print out all board square indexes and names for reference
	for i := 0; i < 64; i++ {
		t.Logf("Square %d = %s", i, BoardSquareNames[i])
	}

	// Create a simple board to test the coordinates
	board := NewBoard()

	// Put a white king at f1 and black pawn at h2
	board.Position = [64]byte{}
	board.Position[SquareF1] = WhiteKing // f1 = 61
	board.Position[SquareG2] = BlackPawn // g2 = 54

	// Update piece locations
	board.UpdatePiecesLocation()

	// Print the board for verification
	t.Log("Board setup:")
	board.PrintPosition()

	// Calculate black pawn attack squares manually
	pawnPos := byte(SquareG2) // g2
	t.Logf("Black pawn at position %d (%s)", pawnPos, BoardSquareNames[pawnPos])

	file := pawnPos % 8
	rank := pawnPos / 8
	t.Logf("Pawn file = %d, rank = %d", file, rank)

	// Calculate left diagonal attack
	if file > 0 {
		leftAttack := ((rank + 1) * 8) + (file - 1)
		t.Logf("Left attack would be at %d (%s)", leftAttack, BoardSquareNames[leftAttack])
	}

	// Calculate right diagonal attack - should not exist for h file
	if file < 7 {
		rightAttack := ((rank + 1) * 8) + (file + 1)
		t.Logf("Right attack would be at %d (%s)", rightAttack, BoardSquareNames[rightAttack])
	} else {
		t.Logf("No right diagonal attack (pawn already at rightmost file)")
	}

	// Check if white king is at any of these diagonal attack positions
	kingPos := board.Pieces.White.King
	t.Logf("White king at %d (%s)", kingPos, BoardSquareNames[kingPos])

	// Now let's recreate the exact failing test scenario
	t.Logf("\n--- Recreating failing test scenario ---")

	// Start with the initial position from the test
	board2 := NewBoard()
	board2.FromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	t.Logf("Initial board position:")
	board2.PrintPosition()

	// Check positions of key pieces
	t.Logf("White king at: %d (%s)", board2.Pieces.White.King, BoardSquareNames[board2.Pieces.White.King])

	// Find the black pawn at h3
	for _, pawn := range board2.Pieces.Black.Pawns {
		if BoardSquareNames[pawn] == "h3" {
			t.Logf("Black pawn at h3: position %d", pawn)
		}
	}

	// Now make the moves that are causing issues
	// First move: white king from e1 to f1
	rootMove := NewMove(SquareE1, SquareF1, 0, [2]byte{0, 0})
	board2.MakeMove(rootMove)

	t.Logf("\nAfter white king moves to f1:")
	board2.PrintPosition()
	t.Logf("White king now at: %d (%s)",
		board2.Pieces.White.King, BoardSquareNames[board2.Pieces.White.King])

	// Show piece at h3 and g2
	t.Logf("Piece at h3 (47): %v", board2.Position[47])
	t.Logf("Piece at g2 (54): %v", board2.Position[54])

	// Second move: black pawn from h3 to h2 (capture)
	// But the test is defining this move as h3 to g2!
	childMove := NewMove(SquareH3, SquareG2, 1, [2]byte{80, 0})
	board2.MakeMove(childMove)

	t.Logf("\nAfter black pawn supposedly moves from h3 to h2:")
	board2.PrintPosition()

	// Show pieces again
	t.Logf("Piece at h3 (47): %v", board2.Position[47])
	t.Logf("Piece at g2 (54): %v", board2.Position[54])
	t.Logf("Piece at h2 (55): %v", board2.Position[55])

	// Check black pawn positions after move
	t.Logf("\nBlack pawn positions after move:")
	for _, pawn := range board2.Pieces.Black.Pawns {
		t.Logf("Black pawn at: %d (%s)", pawn, BoardSquareNames[pawn])
	}

	// Generate attacked squares by black
	board2.GenerateAttackedSquares(false)

	// Check if white king is attacked
	if board2.AttackedSquares[61] {
		t.Logf("White king at f1 IS attacked after the move")

		// Print all attacked squares to see what's going on
		t.Logf("All attacked squares:")
		for i, attacked := range board2.AttackedSquares {
			if attacked {
				t.Logf("Square %d (%s) is attacked", i, BoardSquareNames[i])
			}
		}
	} else {
		t.Logf("White king at f1 is NOT attacked after the move")
	}
}

// This test checks that castling rights are lost if a rook is captured on its original square by a promotion-capture.
func TestPromotionCaptureRemovesCastlingRights(t *testing.T) {
	// FEN: Black rook on h8, white pawn on g7, white king on e1, black king on e8, white to move
	// White plays g8=Q capturing rook on h8
	board := NewBoard()
	board.FromFEN("4k2r/6P1/8/8/8/8/8/4K3 w Kk - 0 1")
	board.GeneratePawnMoves(true)
	found := false
	for _, move := range board.Moves {
		if move.MoveType == MovePromotionCapture && move.To == SquareH8 {
			found = true
			board.MakeMove(move)
			if board.CastlingAvailability.BlackKingSide {
				t.Errorf("Castling rights should be lost after rook is captured by promotion-capture")
			}
		}
	}
	if !found {
		t.Errorf("No promotion-capture move found to h8")
	}
}

func TestMoveKingIntoCheck(t *testing.T) {
	board := NewBoard()
	board.Initialize([64]byte{}) // Properly initialize the board
	// Place white king on e1 (SquareE1), black rook on e8 (SquareE8), white to move
	board.Position[SquareE1] = WhiteKing
	board.Position[SquareE8] = BlackRook
	board.WhiteToMove = true
	board.UpdatePiecesLocation()
	board.GenerateLegalMoves()
	// Check that moves where the king would move into the rook's attack line are not generated
	for _, move := range board.Moves {
		if move.From == SquareE1 && (move.To == SquareE2 || move.To == SquareE3 || move.To == SquareE4 || move.To == SquareE5 || move.To == SquareE6 || move.To == SquareE7 || move.To == SquareE8) {
			t.Errorf("King should not be able to move to %d which is in the rook's attack line", move.To)
		}
	}
}

func TestMoveKingAdjacentToEnemyKing(t *testing.T) {
	board := NewBoard()
	// Place white king on e1 (SquareE1), black king on e2 (SquareE2), white to move
	board.Position[SquareE1] = WhiteKing
	board.Position[SquareE2] = BlackKing
	board.WhiteToMove = true
	board.UpdatePiecesLocation()
	board.GenerateLegalMoves()
	for _, move := range board.Moves {
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
	// Blocked pawn at e2 by e3
	board.Position = [64]byte{}
	board.Position[SquareE2] = WhitePawn
	board.Position[SquareE3] = BlackPawn
	board.UpdatePiecesLocation()
	board.WhiteToMove = true
	board.GeneratePawnMoves(true)
	if len(board.Moves) != 0 {
		t.Errorf("Blocked pawn should have no moves")
	}
}

func TestPawnPromotionWithCaptureAndNoCapture(t *testing.T) {
	board := NewBoard()
	// White pawn on g7, black rook on h8
	board.Position = [64]byte{}
	board.Position[SquareG7] = WhitePawn
	board.Position[SquareH8] = BlackRook
	board.UpdatePiecesLocation()
	board.WhiteToMove = true
	board.GeneratePawnMoves(true)
	promotion, capture := false, false
	for _, move := range board.Moves {
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
	// Knight on a1 and h8
	board.Position = [64]byte{}
	board.Position[SquareA1] = WhiteKnight
	board.Position[SquareH8] = WhiteKnight
	board.UpdatePiecesLocation()
	board.WhiteToMove = true
	board.GenerateKnightMoves(true)
	if len(board.Moves) != 4 {
		t.Errorf("Knights on corners should have 2 moves each (total 4 moves), got %d", len(board.Moves))
	}
}

func TestRookBlockedAndEdge(t *testing.T) {
	board := NewBoard()
	// Rook on a1, blocked by pawn on a2
	board.Position = [64]byte{}
	board.Position[SquareA1] = WhiteRook
	board.Position[SquareA2] = WhitePawn
	board.UpdatePiecesLocation()
	board.WhiteToMove = true
	board.GenerateRookMoves(true)
	if len(board.Moves) != 7 {
		t.Errorf("Rook should have 7 moves along the 1st rank, got %d", len(board.Moves))
	}
}

func TestQueenBlockedAndEdge(t *testing.T) {
	board := NewBoard()
	// Queen on a1, blocked by pawn on a2 and b2
	board.Position = [64]byte{}
	board.Position[SquareA1] = WhiteQueen
	board.Position[SquareA2] = WhitePawn
	board.Position[SquareB2] = WhitePawn
	board.UpdatePiecesLocation()
	board.WhiteToMove = true
	board.GenerateQueenMoves(true)
	if len(board.Moves) != 7 {
		t.Errorf("Queen should have 7 moves along the 1st rank, got %d", len(board.Moves))
	}
}

func TestKingCannotMoveIntoCheck(t *testing.T) {
	board := NewBoard()
	// King on e1, enemy rook on e2
	board.Position[SquareE1] = WhiteKing
	board.Position[SquareE3] = BlackRook
	board.UpdatePiecesLocation()
	board.WhiteToMove = true
	board.GenerateLegalMoves()
	for _, move := range board.Moves {
		if move.From == SquareE1 && move.To == SquareE2 {
			t.Errorf("King should not be able to move into check")
		}
	}
}
