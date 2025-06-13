package libra_test

import (
	. "libra/pkg"
	"testing"
)

var tt = NewTranspositionTable()

func TestSearch5(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	score, move := board.Search(4, tt)

	if score < 200 {
		t.Errorf("Expected score > 200, got %d", score)
	}
	if move == nil || move.ToUCI() != "d7c8q" {
		t.Errorf("Expected move d7c8q, got %s", move.ToUCI())
	}
}

func TestSearch4(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	score, move := board.Search(4, tt)

	if score > -200 {
		t.Errorf("Expected score > -200, got %d", score)
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

/*
func TestSearchMate(t *testing.T) {
	board := NewBoard()
	board.FromFEN("2Q1k3/8/7Q/8/8/8/8/4K3 w - - 0 1")
	score, move := board.Search(7, tt)

	fmt.Printf("Score: %d, Move: %s\n", score, move.ToUCI())

}
*/
