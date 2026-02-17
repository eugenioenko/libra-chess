package libra

import "sync"

const (
	BoundExact = iota
	BoundLower // score is a lower bound (beta cutoff)
	BoundUpper // score is an upper bound (failed low)
)

type TTEntry struct {
	Score    int
	BestMove Move
	Depth    int
	Bound    byte
}

type TranspositionTable struct {
	table map[uint64]TTEntry
	mu    sync.RWMutex
}

func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		table: make(map[uint64]TTEntry),
	}
}

func (tt *TranspositionTable) Get(hash uint64, depth int) (TTEntry, bool) {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	entry, ok := tt.table[hash]
	if !ok || entry.Depth < depth {
		return TTEntry{}, false
	}
	return entry, true
}

func (tt *TranspositionTable) Set(hash uint64, depth int, value int, bestMove Move, bound byte) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	entry, ok := tt.table[hash]
	if !ok || depth >= entry.Depth {
		tt.table[hash] = TTEntry{Score: value, BestMove: bestMove, Depth: depth, Bound: bound}
	}
}

func (tt *TranspositionTable) Clear() {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	tt.table = make(map[uint64]TTEntry)
}

func (tt *TranspositionTable) Size() int {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	return len(tt.table)
}

func (tt *TranspositionTable) BestMoveDeepest(hash uint64) *Move {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	entry, ok := tt.table[hash]
	if !ok {
		return nil
	}
	return &entry.BestMove
}
