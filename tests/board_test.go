package testing

import (
	"fmt"
	libra "libra/pkg"
	"testing"
)

func TestShouldInsertValues(t *testing.T) {
	board := libra.NewBoard()
	ok, error := board.LoadFromFEN(libra.BoardInitialFEN)

	if error != nil {
		fmt.Println(error)
	}

	if ok != true {
		t.Fail()
	}
}

func TestShouldCalculateZobristHash(t *testing.T) {
	board := libra.NewBoard()
	board.LoadFromFEN(libra.BoardInitialFEN)
	hashA := libra.ZobristHash(board)
	board.LoadFromFEN("rnbqkb1r/ppp1p1pp/8/3p1pP1/3Pn3/2N5/PPP1PP1P/R1BQKBNR w KQkq - 1 5")
	hashB := libra.ZobristHash(board)
	board.Print()
	if hashA == 0 {
		t.Fail()
	}

	if hashA == hashB {
		t.Fail()
	}
}

func TestShouldGeneratePawnMoves(t *testing.T) {
	board := libra.NewBoard()
	board.LoadFromFEN(libra.BoardInitialFEN)
	board.GeneratePawnMoves(board.WhiteToMove)
	if len(board.Moves.All) != 16 {
		t.Fail()
	}
	if len(board.Moves.Captures) != 0 {
		t.Fail()
	}
	board.LoadFromFEN("rnbqkbnr/4p3/pp1p1p1p/2p3p1/1P1PPPP1/P1P5/7P/RNBQKBNR w KQkq - 0 8")
	board.GeneratePawnMoves(board.WhiteToMove)
	if len(board.Moves.All) != 11 {
		t.Fail()
	}
	if len(board.Moves.Captures) != 3 {
		t.Fail()
	}
}

func TestShouldGenerateOnPassantPawnMoves(t *testing.T) {
	board := libra.NewBoard()
	board.LoadFromFEN("rnbqkbnr/8/pp1p1p1p/2p1pPp1/1P1PP1P1/P1P5/7P/RNBQKBNR w KQkq e6 0 9")
	board.GeneratePawnMoves(board.WhiteToMove)
	if len(board.Moves.Captures) != 4 {
		t.Fail()
	}
}

func TestShouldGeneratePromotionMoves(t *testing.T) {
	board := libra.NewBoard()
	board.LoadFromFEN("3n4/4P3/8/8/2K2k2/8/8/8 w - - 0 1")
	board.GeneratePawnMoves(board.WhiteToMove)
	if len(board.Moves.Promotions) != 8 {
		t.Fail()
	}

	board.LoadFromFEN("8/2P5/8/4k3/8/4K3/7p/8 b - - 0 1")
	board.GeneratePawnMoves(board.WhiteToMove)
	if len(board.Moves.Promotions) != 4 {
		t.Fail()
	}

}

func TestShouldGenerateRookMoves(t *testing.T) {
	board := libra.NewBoard()

	board.LoadFromFEN("1k4r1/8/2R4p/8/8/8/8/7K")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves.All) != 14 {
		t.Fail()
	}
	if len(board.Moves.Captures) != 1 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/R7")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves.All) != 14 {
		t.Fail()
	}

	board.LoadFromFEN("R7/8/8/8/8/8/8/8")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves.All) != 14 {
		t.Fail()
	}

	board.LoadFromFEN("7R/8/8/8/8/8/8/8")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves.All) != 14 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/7R")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves.All) != 14 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/4R3/8/8/8/8/8")
	board.GenerateRookMoves(board.WhiteToMove)
	if len(board.Moves.All) != 14 {
		t.Fail()
	}

}

func TestShouldGenerateBishopMoves(t *testing.T) {
	board := libra.NewBoard()

	board.LoadFromFEN("8/8/8/4B3/8/8/8/8")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves.All) != 13 {
		t.Fail()
	}

	board.LoadFromFEN("B7/8/8/8/8/8/8/8")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves.All) != 7 {
		t.Fail()
	}

	board.LoadFromFEN("7B/8/8/8/8/8/8/8")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves.All) != 7 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/7B")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves.All) != 7 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/B7")
	board.GenerateBishopMoves(board.WhiteToMove)
	if len(board.Moves.All) != 7 {
		t.Fail()
	}
}

func TestShouldGenerateQueenMoves(t *testing.T) {
	board := libra.NewBoard()

	board.LoadFromFEN("8/8/8/4Q3/8/8/8/8")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves.All) != 27 {
		t.Fail()
	}

	board.LoadFromFEN("Q7/8/8/8/8/8/8/8")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves.All) != 21 {
		t.Fail()
	}

	board.LoadFromFEN("7Q/8/8/8/8/8/8/8")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves.All) != 21 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/7Q")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves.All) != 21 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/Q7")
	board.GenerateQueenMoves(board.WhiteToMove)
	if len(board.Moves.All) != 21 {
		t.Fail()
	}
}

func TestShouldGenerateKingMoves(t *testing.T) {
	board := libra.NewBoard()

	board.LoadFromFEN("8/8/8/4K3/8/8/8/8")
	board.GenerateKingMoves(board.WhiteToMove)
	if len(board.Moves.All) != 8 {
		t.Fail()
	}

	board.LoadFromFEN("K7/8/8/8/8/8/8/8")
	board.GenerateKingMoves(board.WhiteToMove)
	if len(board.Moves.All) != 3 {
		t.Fail()
	}

	board.LoadFromFEN("7Q/8/8/8/8/8/8/8")
	board.GenerateKingMoves(board.WhiteToMove)
	if len(board.Moves.All) != 3 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/7Q")
	board.GenerateKingMoves(board.WhiteToMove)
	if len(board.Moves.All) != 3 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/K7")
	board.GenerateKingMoves(board.WhiteToMove)
	if len(board.Moves.All) != 3 {
		t.Fail()
	}
}

func TestShouldGenerateKnightMoves(t *testing.T) {
	board := libra.NewBoard()

	board.LoadFromFEN("8/8/8/4N3/8/8/8/8")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves.All) != 8 {
		t.Fail()
	}

	board.LoadFromFEN("N7/8/8/8/8/8/8/8")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves.All) != 2 {
		t.Fail()
	}

	board.LoadFromFEN("7N/8/8/8/8/8/8/8")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves.All) != 2 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/7N")
	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves.All) != 2 {
		t.Fail()
	}

	board.LoadFromFEN("8/8/8/8/8/8/8/N7")

	board.GenerateKnightMoves(board.WhiteToMove)
	if len(board.Moves.All) != 2 {
		t.Fail()
	}
}

func TestShouldGenerateAttackVector(t *testing.T) {
	board := libra.NewBoard()

	board.LoadFromFEN("1k6/ppp5/8/3b2R1/3B4/8/5PPP/6K1 b - - 0 1")
	board.GenerateMoves()
	if len(board.Moves.All) != 8 {
		t.Fail()
	}

}
