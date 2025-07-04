package libra_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

var tt = NewTranspositionTable()

func TestSearch5(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	board.Search(5, tt, 0, nil)
}

func TestSearch4(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	result := board.Search(5, tt, 0, nil)

	if result.BestScore > -200 {
		t.Errorf("Expected score > -200, got %d", result.BestScore)
	}
	if result.BestMove == nil {
		t.Errorf("Expected a move, got nil")
		return
	}

	uci := result.BestMove.ToUCI()

	if uci != "c4c5" {
		t.Errorf("Expected move c4c5, got %s", uci)
	}
}
func TestCaptureWithLessFirst(t *testing.T) {
	board := NewBoard()
	board.FromFEN("k7/8/4p3/3r4/2B1Q3/8/8/7K w - - 0 1")
	result := board.Search(1, tt, 0, nil)

	if result.BestMove == nil || result.BestMove.ToUCI() != "c4d5" {
		t.Errorf("Expected move c4d5, got %s", result.BestMove.ToUCI())
	}

	board.FromFEN("k7/8/4p3/3r4/2Q1B3/8/8/7K w - - 0 1")
	result = board.Search(2, tt, 0, nil)

	if result.BestMove == nil || result.BestMove.ToUCI() != "e4d5" {
		t.Errorf("Expected move e4d5, got %s", result.BestMove.ToUCI())
	}
}

func TestPreferMateInsteadOfCapture(t *testing.T) {
	board := NewBoard()
	board.FromFEN("k7/8/4p3/3r4/2Q1B3/8/8/7K w - - 0 1")
	result := board.Search(5, tt, 0, nil)

	if result.BestMove == nil || result.BestMove.ToUCI() != "c4c7" {
		t.Errorf("Expected move c4c7, got %s", result.BestMove.ToUCI())
	}
}

func TestSearchPerft1(t *testing.T) {
	board := NewBoard()
	board.LoadInitial()
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft2(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1") // Corrected FEN
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft3(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft4(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft5(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8") // Corrected FEN
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft6(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10 ") // Corrected FEN
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft7(t *testing.T) {
	board := NewBoard()
	board.FromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10 ")
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft8(t *testing.T) {
	board := NewBoard()
	board.FromFEN("4k2r/2b2ppp/5n2/7P/1Q1N4/4P3/5PP1/KR6 w k - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft9(t *testing.T) {
	board := NewBoard()
	board.FromFEN("1r5k/2b2ppp/5n2/NPP4P/PKR2B2/8/8/8 w - - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearchPerft10(t *testing.T) {
	board := NewBoard()
	board.FromFEN("8/8/ppk5/2p5/1P6/PKP5/8/8 w - - 0 1")
	tt := NewTranspositionTable()
	board.Search(5, tt, 0, nil)
}

func TestSearch123(t *testing.T) {
	board := NewBoard()
	board.ParseAndApplyPosition(strings.Fields("fen rnbqkbnr/ppp2ppp/4p3/3p4/3P4/1P5P/P1P1PPP1/RNBQKBNR b KQkq - 0 1 moves b8c6 e2e3 d8f6 f1a6 b7a6 d1g4 f8b4 c1d2 b4d6 g4f4 d6f4 e3f4 c6d4 b1a3 d4f3 e1e2 f6a1 e2f3 g8f6 d2a5 a1a2 a5b4 d5d4 b4e1 a2a3"))
	fmt.Println(board.ToFEN())
	move := board.IterativeDeepeningSearch(SearchOptions{
		TimeLimitInMs: 1_000,
	})
	fmt.Printf("Best move: %s\n", move.ToUCI())
}
