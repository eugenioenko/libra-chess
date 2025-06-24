package libra

import (
	"fmt"
	"sync/atomic"
	"time"
)

// SearchStats holds only statistics for a search
// (thread-safe via atomic operations)
type SearchStats struct {
	NodesSearched        uint64 // Total nodes visited
	NodesPruned          uint64 // Nodes cut off by alpha-beta pruning
	TTHits               uint64 // Transposition table hits
	TTStores             uint64 // Transposition table stores
	BetaCutoffs          uint64 // Beta cutoffs (prunes)
	NullMovePrunes       uint64 // Null move pruning occurrences
	MoveGenerations      uint64 // Number of times legal moves were generated
	MaxSearchDepth       int32  // Maximum depth reached in the search
	TimeSpentNanoseconds int64  // Total time taken for the search (nanoseconds)
	BestScore            int    // Best score found in the search

	startTime time.Time // unexported field for tracking time
}

// Increment functions for thread-safe updates
func (s *SearchStats) IncNodesSearched() {
	atomic.AddUint64(&s.NodesSearched, 1)
}

func (s *SearchStats) IncNodesPruned() {
	atomic.AddUint64(&s.NodesPruned, 1)
}

func (s *SearchStats) IncTTHit() {
	atomic.AddUint64(&s.TTHits, 1)
}

func (s *SearchStats) IncTTStore() {
	atomic.AddUint64(&s.TTStores, 1)
}

func (s *SearchStats) IncBetaCutoff() {
	atomic.AddUint64(&s.BetaCutoffs, 1)
}

func (s *SearchStats) IncNullMovePrune() {
	atomic.AddUint64(&s.NullMovePrunes, 1)
}

func (s *SearchStats) IncMoveGeneration() {
	atomic.AddUint64(&s.MoveGenerations, 1)
}

// SetMaxSearchDepth sets the maximum search depth reached (thread-safe)
func (s *SearchStats) SetMaxSearchDepth(depth int32) {
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

// AddTimeSpent adds duration to TimeSpentNanoseconds (thread-safe)
func (s *SearchStats) AddTimeSpent(d time.Duration) {
	atomic.AddInt64(&s.TimeSpentNanoseconds, d.Nanoseconds())
}

// For timing
func (s *SearchStats) StartTimer() {
	s.startTime = time.Now()
}

func (s *SearchStats) StopTimer() {
	if !s.startTime.IsZero() {
		d := time.Since(s.startTime)
		s.AddTimeSpent(d)
		s.startTime = time.Time{} // reset
	}
}

func (s *SearchStats) String() string {
	dur := time.Duration(s.TimeSpentNanoseconds)
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
Time Spent:            %s
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
		dur.String(),
	)
}

func (s *SearchStats) Print() {
	fmt.Print(s.String())
}

// PrintUCI prints the search stats in UCI info format
func (s *SearchStats) PrintUCI() {
	dur := time.Duration(s.TimeSpentNanoseconds)
	nps := int64(0)
	if dur > 0 {
		nps = int64(float64(s.NodesSearched) / dur.Seconds())
	}
	nodesTotal := s.NodesSearched + s.NodesPruned
	prunedPercent := 0.0
	if nodesTotal > 0 {
		prunedPercent = (float64(s.NodesPruned) / float64(nodesTotal)) * 100.0
	}
	fmt.Printf("info depth %d score cp %d nodes %d nps %d prun %.0f%% time %d\n",
		s.MaxSearchDepth,
		s.BestScore,
		s.NodesSearched,
		nps,
		prunedPercent,
		dur.Milliseconds(),
	)
}
