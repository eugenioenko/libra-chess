package libra_test

import (
	"fmt"
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

// Test that zobrist hashing works correctly for all board position and
// castling rights independently of the move number.
// This is a fixed test that should always return the same Zobrist key for the same
func TestZobristFixedTest(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -")

	if board.ZobristHash() != 0xea4d292ae31746cb {
		t.Errorf("Expected Zobrist key 0xea4d292ae31746cb, got %x", board.ZobristHash())
	}

	board.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	if board.ZobristHash() != 0xea4d292ae31746cb {
		t.Errorf("Expected Zobrist key 0xea4d292ae31746cb, got %x", board.ZobristHash())
	}

	board.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 1 8")

	if board.ZobristHash() != 0xea4d292ae31746cb {
		t.Errorf("Expected Zobrist key 0xea4d292ae31746cb, got %x", board.ZobristHash())
	}
}

func TestZobristEnPassante(t *testing.T) {
	board := NewBoard()
	board.FromFEN("rnbqkb1r/pp3ppp/3p1n2/2pPp3/2P5/2N5/PP2PPPP/R1BQKBNR w KQkq -")
	fmt.Println(board.ToFEN())
	hash1 := board.ZobristHash()

	board.FromFEN("rnbqkb1r/pp3ppp/3p1n2/2pPp3/2P5/2N5/PP2PPPP/R1BQKBNR w KQkq e6")
	hash2 := board.ZobristHash()
	fmt.Println(board.ToFEN())

	if hash1 == hash2 {
		t.Errorf("Hashes should be different if there is enPassant active pawn")
	}

}

func TestZobristHashInitialWasm(t *testing.T) {
	board := NewBoard()
	board.LoadInitial()
	hash1 := board.ZobristHash()

	board.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -")
	hash2 := board.ZobristHash()

	board.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	hash3 := board.ZobristHash()

	if hash3 != hash1 || hash1 != hash2 || hash1 != 0xea4d292ae31746cb {
		t.Errorf("Hashes should be equal for initial position and FEN string")
	}

}
