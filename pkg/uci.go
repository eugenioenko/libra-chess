package libra

import "fmt"

func GetUCIRemainingTime(whiteToMove bool, fields []string) int {
	var wtime, btime int
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "wtime":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &wtime)
			}
		case "btime":
			if i+1 < len(fields) {
				fmt.Sscanf(fields[i+1], "%d", &btime)
			}
		}
	}
	if whiteToMove {
		return wtime
	}
	return btime
}
