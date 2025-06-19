package libra

import (
	"runtime"
	"sync"
)

// PerftParallel parallelizes the top-level perft search using a worker pool limited by GOMAXPROCS.
func (board *Board) PerftParallel(depth int) int {
	if depth == 0 {
		return 1
	}
	moves := board.GenerateLegalMoves()
	if depth == 1 {
		return len(moves)
	}

	numWorkers := runtime.GOMAXPROCS(0)
	moveChan := make(chan Move, len(moves))
	results := make(chan int, len(moves))
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for m := range moveChan {
				clone := board.Clone()
				clone.Move(m)
				nodes := clone.Perft(depth - 1)
				results <- nodes
			}
		}()
	}

	// Send jobs
	for _, move := range moves {
		moveChan <- move
	}
	close(moveChan)

	wg.Wait()
	close(results)

	nodes := 0
	for n := range results {
		nodes += n
	}
	return nodes
}

func (board *Board) Perft(depth int) int {
	if depth == 0 {
		return 1
	}
	moves := board.GenerateLegalMoves()
	if depth == 1 {
		return len(moves)
	}
	nodes := 0
	for _, move := range moves {
		state := board.Move(move)
		nodes += board.Perft(depth - 1)
		board.UndoMove(state)
	}
	return nodes
}
