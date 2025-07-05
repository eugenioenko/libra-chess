package libra_test

import (
	"fmt"
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

// Perft tests for the Libra chess engine.
// These tests are based on the known perft values for various positions.
// The positions are taken from the Chess Programming Wiki.
// https://www.chessprogramming.org/Perft_Results
func TestPerft1(t *testing.T) {
	board := NewBoard()
	board.LoadInitial()
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)
	n5 := board.PerftParallel(5)
	// n6 := board.PerftParallel(6)

	if n1 != 20 || n2 != 400 || n3 != 8902 || n4 != 197281 || n5 != 4865609 { //} || n6 != 119060324 {
		t.Fail()
	}

}

func TestPerft2(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)

	if n1 != 48 || n2 != 2039 || n3 != 97862 || n4 != 4085603 {
		t.Fail()
	}

}

func TestPerft3(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)
	n5 := board.PerftParallel(5)
	n6 := board.PerftParallel(6)
	if n1 != 14 || n2 != 191 || n3 != 2812 || n4 != 43238 || n5 != 674624 || n6 != 11030083 {
		t.Fail()
	}
}

func TestPerft4(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)
	n5 := board.PerftParallel(5)

	if n1 != 6 || n2 != 264 || n3 != 9467 || n4 != 422333 || n5 != 15833292 {
		t.Fail()
	}
}

func TestPerft5(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)
	fmt.Printf("Perft5: %d %d %d %d\n", n1, n2, n3, n4)

	if n1 != 44 || n2 != 1486 || n3 != 62379 || n4 != 2103487 {
		t.Fail()
	}
}

func TestPerft6(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10 ")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)

	if n1 != 46 || n2 != 2079 || n3 != 89890 || n4 != 3894594 {
		t.Fail()
	}
}

func TestPerft7(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10 ")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)

	if n1 != 46 || n2 != 2079 || n3 != 89890 || n4 != 3894594 {
		t.Fail()
	}
}

func TestPerft8(t *testing.T) {
	board := NewBoard()
	board.FromFEN("4k2r/2b2ppp/5n2/7P/1Q1N4/4P3/5PP1/KR6 w k - 0 1")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)

	if n1 != 41 || n2 != 824 || n3 != 33818 || n4 != 684104 {
		t.Fail()
	}
}

func TestPerft9(t *testing.T) {
	board := NewBoard()
	board.FromFEN("1r5k/2b2ppp/5n2/NPP4P/PKR2B2/8/8/8 w - - 0 1")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)

	if n1 != 24 || n2 != 623 || n3 != 14368 || n4 != 354843 {
		t.Fail()
	}
}

func TestPerft10(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/8/ppk5/2p5/1P6/PKP5/8/8 w - - 0 1")
	n1 := board.PerftParallel(1)
	n2 := board.PerftParallel(2)
	n3 := board.PerftParallel(3)
	n4 := board.PerftParallel(4)
	n5 := board.PerftParallel(5)
	n6 := board.PerftParallel(6)
	n7 := board.PerftParallel(7)

	if n1 != 9 || n2 != 78 || n3 != 600 || n4 != 5369 || n5 != 42632 || n6 != 381238 || n7 != 3058563 {
		t.Fail()
	}
}
