package libra

import "sync"

type TTKey struct {
	Hash  uint64
	Depth int
}

type TTEntry struct {
	Score    int
	BestMove Move
	Depth    int
}

type TranspositionTable struct {
	table map[uint64]TTEntry // hash -> entry
	mu    sync.RWMutex
}

func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		table: make(map[uint64]TTEntry),
	}
}

func (tt *TranspositionTable) Get(hash uint64, depth int) (int, bool) {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	entry, ok := tt.table[hash]
	if !ok || entry.Depth < depth {
		return 0, false
	}
	return entry.Score, true
}

func (tt *TranspositionTable) Set(hash uint64, depth int, value int, bestMove Move) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	entry, ok := tt.table[hash]
	if !ok || depth >= entry.Depth {
		tt.table[hash] = TTEntry{Score: value, BestMove: bestMove, Depth: depth}
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
