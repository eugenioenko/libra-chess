package libra_test

import (
	"testing"

	. "github.com/eugenioenko/libra-chess/pkg"
)

func BenchmarkPerftParallel1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		board := NewBoard()
		board.LoadInitial()
		board.PerftParallel(1)
	}
}

func BenchmarkPerftParallel2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		board := NewBoard()
		board.LoadInitial()
		board.PerftParallel(2)
	}
}

func BenchmarkPerftParallel3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		board := NewBoard()
		board.LoadInitial()
		board.PerftParallel(3)
	}
}

func BenchmarkPerftParallel4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		board := NewBoard()
		board.LoadInitial()
		board.PerftParallel(4)
	}
}

func BenchmarkPerftParallel5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		board := NewBoard()
		board.LoadInitial()
		board.PerftParallel(5)
	}
}

func BenchmarkPerftParallel6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		board := NewBoard()
		board.LoadInitial()
		board.PerftParallel(6)
	}
}
