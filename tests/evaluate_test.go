package libra_test

import (
	. "libra/pkg"
	"testing"
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
	board.FromFEN("7K/8/8/8/8/8/8/k7 w - - 0 1")
	// Remove black queen
	score := board.Evaluate()
	if score != 0 {
		t.Errorf("Kings don't have value %d", score)
	}
}

func TestEvaluate_BlackNoQueen(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	score := board.Evaluate()
	if score != 895 {
		t.Errorf("Expected 895, got %d", score)
	}
}

func TestEvaluate_NoKings(t *testing.T) {
	board := NewBoard()
	score := board.Evaluate()
	if score != 0 {
		t.Errorf("No kings, expected 0, got %d", score)
	}
}
