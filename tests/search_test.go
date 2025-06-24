package libra_test

import (
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

var tt = NewTranspositionTable()

func TestSearch5(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	_, stats := board.Search(5, tt)
	stats.Print()
}
func TestSearch6(t *testing.T) {
	board := NewBoard()
	board.LoadInitial()
	_, stats := board.Search(5, tt)
	stats.Print()
}

func TestSearch4(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	move, stats := board.Search(7, tt)
	if stats.BestScore > -200 {
		t.Errorf("Expected score > -200, got %d", stats.BestScore)
	}
	if move == nil {
		t.Errorf("Expected a move, got nil")
		return
	}

	uci := move.ToUCI()

	if uci != "c4c5" {
		t.Errorf("Expected move c4c5, got %s", uci)
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
	move, stats := board.Search(5, tt)
	stats.Print()

	if move == nil || move.ToUCI() != "c4c7" {
		t.Errorf("Expected move c4c7, got %s", move.ToUCI())
	}
}
