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
	hashA := libra.ZobristHash(board.Position)
	board.LoadFromFEN("rnbqkb1r/ppp1p1pp/8/3p1pP1/3Pn3/2N5/PPP1PP1P/R1BQKBNR w KQkq - 1 5")
	hashB := libra.ZobristHash(board.Position)
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
	board.GeneratePawnMoves()
	if len(board.Moves) != 16 {
		t.Fail()
	}
	if len(board.Captures) != 0 {
		t.Fail()
	}
	board.LoadFromFEN("rnbqkbnr/4p3/pp1p1p1p/2p3p1/1P1PPPP1/P1P5/7P/RNBQKBNR w KQkq - 0 8")
	board.GeneratePawnMoves()
	if len(board.Moves) != 11 {
		t.Fail()
	}
	if len(board.Captures) != 3 {
		t.Fail()
	}
}

func TestShouldGenerateOnPassantPawnMoves(t *testing.T) {
	board := libra.NewBoard()
	board.LoadFromFEN("rnbqkbnr/8/pp1p1p1p/2p1pPp1/1P1PP1P1/P1P5/7P/RNBQKBNR w KQkq e6 0 9")
	board.GeneratePawnMoves()
	if len(board.Captures) != 3 {
		t.Fail()
	}
}
