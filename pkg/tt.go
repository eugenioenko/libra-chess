package libra

import "sync"

type TTKey struct {
	Hash  uint64
	Depth int
}

type TranspositionTable struct {
	table map[TTKey]int
	mu    sync.RWMutex
}

func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		table: make(map[TTKey]int),
	}
}

func (tt *TranspositionTable) Get(hash uint64, depth int) (int, bool) {
	key := TTKey{Hash: hash, Depth: depth}
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	val, ok := tt.table[key]
	return val, ok
}

func (tt *TranspositionTable) Set(hash uint64, depth int, value int) {
	key := TTKey{Hash: hash, Depth: depth}
	tt.mu.Lock()
	defer tt.mu.Unlock()
	tt.table[key] = value
}

func (tt *TranspositionTable) Clear() {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	tt.table = make(map[TTKey]int)
}

func (tt *TranspositionTable) Size() int {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	return len(tt.table)
}
