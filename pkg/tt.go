package libra

type TTKey struct {
	Hash  uint64
	Depth int
}

type TranspositionTable struct {
	table map[TTKey]int
}

func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		table: make(map[TTKey]int),
	}
}

func (tt *TranspositionTable) Get(hash uint64, depth int) (int, bool) {
	key := TTKey{Hash: hash, Depth: depth}
	val, ok := tt.table[key]
	return val, ok
}

func (tt *TranspositionTable) Set(hash uint64, depth int, value int) {
	key := TTKey{Hash: hash, Depth: depth}
	tt.table[key] = value
}

func (tt *TranspositionTable) Clear() {
	tt.table = make(map[TTKey]int)
}

func (tt *TranspositionTable) Size() int {
	return len(tt.table)
}
