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
	hashA := board.ZobristHash()
	board.LoadFromFEN("rnbqkb1r/ppp1p1pp/8/3p1pP1/3Pn3/2N5/PPP1PP1P/R1BQKBNR w KQkq - 1 5")
	hashB := board.ZobristHash()
	board.Print()
	if hashA == 0 {
		t.Fail()
	}

	if hashA == hashB {
		t.Fail()
	}
}
