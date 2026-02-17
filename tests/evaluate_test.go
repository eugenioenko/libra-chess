package libra_test

import (
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

func TestEvaluate_InitialPosition(t *testing.T) {
	board := NewBoard()
	board.LoadInitial()
	score := board.Evaluate()
	if score != 0 {
		t.Errorf("Initial position should be 0, got %d", score)
	}
}

func TestEvaluate_TwoKings(t *testing.T) {
	board := NewBoard()
	board.FromFEN("4k3/8/8/8/8/8/8/4K3 w - - 0 1")
	score := board.Evaluate()
	if score != 0 {
		t.Errorf("Symmetric kings should be 0, got %d", score)
	}
}

func TestEvaluate_BlackNoQueen(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	score := board.Evaluate()
	if score != 1011 {
		t.Errorf("Expected 1011, got %d", score)
	}
}

func TestEvaluate_NoKings(t *testing.T) {
	board := NewBoard()
	score := board.Evaluate()
	if score != 0 {
		t.Errorf("No kings, expected 0, got %d", score)
	}
}
