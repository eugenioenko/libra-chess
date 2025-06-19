package libra

func CharIsNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

func SquareNameToIndex(name string) (byte, bool) {
	val, ok := squareNameToIndex[name]
	return val, ok
}

func MathMinByte(a byte, b byte) byte {
	if a > b {
		return b
	}
	return a
}

// SquareIndexToName converts a square index (0-63) to its algebraic name (e.g., 0 -> "a8").
func SquareIndexToName(idx byte) (string, bool) {
	if idx >= 64 {
		return "", false
	}
	return BoardSquareNames[idx], true
}

func PieceCodeToFont(piece byte) string {
	return pieceCodeToFont[piece]
}
