package libra_test

import (
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

var tt = NewTranspositionTable()

func TestSearch5(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	board.Search(5, tt)
}

func TestSearch4(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	move, stats := board.Search(4, tt)

	if stats.BestScore > -200 {
		t.Errorf("Expected score > -200, got %d", stats.BestScore)
	}
	if move == nil {
		t.Errorf("Expected a move, got nil")
		return
	}

	uci := move.ToUCI()

	if uci != "g1h1" {
		t.Errorf("Expected move g1h1, got %s", uci)
	}
}
func TestCaptureWithLessFirst(t *testing.T) {
	board := NewBoard()
	board.FromFEN("k7/8/4p3/3r4/2B1Q3/8/8/7K w - - 0 1")
	move, _ := board.Search(1, tt)

	if move == nil || move.ToUCI() != "c4d5" {
		t.Errorf("Expected move c4d5, got %s", move.ToUCI())
	}

	board.FromFEN("k7/8/4p3/3r4/2Q1B3/8/8/7K w - - 0 1")
	move, _ = board.Search(2, tt)

	if move == nil || move.ToUCI() != "e4d5" {
		t.Errorf("Expected move e4d5, got %s", move.ToUCI())
	}
}

func TestPreferMateInsteadOfCapture(t *testing.T) {
	board := NewBoard()
	board.FromFEN("k7/8/4p3/3r4/2Q1B3/8/8/7K w - - 0 1")
	move, _ := board.Search(5, tt)

	if move == nil || move.ToUCI() != "c4c7" {
		t.Errorf("Expected move c4c7, got %s", move.ToUCI())
	}
}

func TestSearchPerft1(t *testing.T) {
	board := NewBoard()
	board.LoadInitial()
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft2(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1") // Corrected FEN
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft3(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft4(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft5(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8") // Corrected FEN
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft6(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10 ") // Corrected FEN
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft7(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10 ")
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft8(t *testing.T) {
	board := NewBoard()
	board.FromFEN("4k2r/2b2ppp/5n2/7P/1Q1N4/4P3/5PP1/KR6 w k - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft9(t *testing.T) {
	board := NewBoard()
	board.FromFEN("1r5k/2b2ppp/5n2/NPP4P/PKR2B2/8/8/8 w - - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt)
}

func TestSearchPerft10(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/8/ppk5/2p5/1P6/PKP5/8/8 w - - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt)
}
