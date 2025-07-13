package libra

import (
	"fmt"
	"sync/atomic"
	"time"
)

type SearchResult struct {
	NodesSearched   uint64 // Total nodes visited
	NodesPruned     uint64 // Nodes cut off by alpha-beta pruning
	TTHits          uint64 // Transposition table hits
	TTStores        uint64 // Transposition table stores
	BetaCutoffs     uint64 // Beta cutoffs (prunes)
	NullMovePrunes  uint64 // Null move pruning occurrences
	MoveGenerations uint64 // Number of times legal moves were generated
	MaxSearchDepth  int32  // Maximum depth reached in the search
	TimeSpentInMs   int64  // Total time taken for the search (milliseconds)
	BestScore       int    // Best score found in the search
	PVMove          string // Best move line in UCI format
	BestMove        *Move  // Best	move in UCI format
	IsInterrupted   bool   // Whether the search was interrupted

	startTime time.Time // unexported field for tracking time
}

// Increment functions for thread-safe updates
func (s *SearchResult) IncNodesSearched() {
	atomic.AddUint64(&s.NodesSearched, 1)
}

func (s *SearchResult) IncNodesPruned() {
	atomic.AddUint64(&s.NodesPruned, 1)
}

func (s *SearchResult) IncTTHit() {
	atomic.AddUint64(&s.TTHits, 1)
}

func (s *SearchResult) IncTTStore() {
	atomic.AddUint64(&s.TTStores, 1)
}

func (s *SearchResult) IncBetaCutoff() {
	atomic.AddUint64(&s.BetaCutoffs, 1)
}

func (s *SearchResult) IncNullMovePrune() {
	atomic.AddUint64(&s.NullMovePrunes, 1)
}

func (s *SearchResult) IncMoveGeneration() {
	atomic.AddUint64(&s.MoveGenerations, 1)
}

// SetMaxSearchDepth sets the maximum search depth reached (thread-safe)
func (s *SearchResult) SetMaxSearchDepth(depth int32) {
	for {
		current := atomic.LoadInt32(&s.MaxSearchDepth)
		if depth > current {
			if atomic.CompareAndSwapInt32(&s.MaxSearchDepth, current, depth) {
				break
			}
		} else {
			break
		}
	}
}

// AddTimeSpent adds duration to TimeSpentInMs (thread-safe)
func (s *SearchResult) AddTimeSpent(d time.Duration) {
	atomic.AddInt64(&s.TimeSpentInMs, d.Milliseconds())
}

// For timing
func (s *SearchResult) StartTimer() {
	s.startTime = time.Now()
}

func (s *SearchResult) StopTimer() {
	if !s.startTime.IsZero() {
		d := time.Since(s.startTime)
		s.AddTimeSpent(d)
		s.startTime = time.Time{} // reset
	}
}

func (s *SearchResult) String() string {
	nodesTotal := s.NodesSearched + s.NodesPruned
	prunedPercent := 0.0
	if nodesTotal > 0 {
		prunedPercent = (float64(s.NodesPruned) / float64(nodesTotal)) * 100.0
	}
	return fmt.Sprintf(`Search Stats:
---------------------
Nodes Searched:        %d
Nodes Pruned:          %d
Nodes Pruned %%:        %.2f%%
Nodes Total:           %d
TT Hits:               %d
TT Stores:             %d
Beta Cutoffs:          %d
Null Move Prunes:      %d
Move Generations:      %d
Max Search Depth:      %d
Best Score:            %d
Time Spent:            %dms
`,
		s.NodesSearched,
		s.NodesPruned,
		prunedPercent,
		nodesTotal,
		s.TTHits,
		s.TTStores,
		s.BetaCutoffs,
		s.NullMovePrunes,
		s.MoveGenerations,
		s.MaxSearchDepth,
		s.BestScore,
		s.TimeSpentInMs,
	)
}

func (s *SearchResult) Print() {
	fmt.Print(s.String())
}

// PrintUCI prints the search stats in UCI info format
func (s *SearchResult) PrintUCI() {
	nps := int64(0)
	if s.TimeSpentInMs > 0 {
		nps = int64(s.NodesSearched) * 1000 / s.TimeSpentInMs
	}
	nodesTotal := s.NodesSearched + s.NodesPruned
	prunedPercent := 0.0
	if nodesTotal > 0 {
		prunedPercent = (float64(s.NodesPruned) / float64(nodesTotal)) * 100.0
	}
	bestMove := "0000"
	if s.BestMove != nil {
		bestMove = s.BestMove.ToUCI()
	}
	fmt.Printf("info depth %d score cp %d nodes %d nps %d prun %.0f%% pv %s time %d\n",
		s.MaxSearchDepth,
		s.BestScore,
		s.NodesSearched,
		nps,
		prunedPercent,
		bestMove,
		s.TimeSpentInMs,
	)
}
