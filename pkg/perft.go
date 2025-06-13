package libra

import "sync"

// PerftParallel parallelizes the top-level perft search using goroutines for each root move.
func (board *Board) PerftParallel(depth int) int {
	if depth == 0 {
		return 1
	}
	moves := board.GenerateLegalMoves()
	if depth == 1 {
		return len(moves)
	}

	var wg sync.WaitGroup
	results := make(chan int, len(moves))

	for _, move := range moves {
		wg.Add(1)
		go func(m Move) {
			defer wg.Done()
			child := board.Clone()
			child.MakeMove(m)
			count := child.Perft(depth - 1)
			results <- count
		}(move)
	}

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
		prev := board.MakeMove(move)
		nodes += board.Perft(depth - 1)
		board.UndoMove(prev)
	}
	return nodes
}
