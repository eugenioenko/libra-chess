package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/eugenioenko/libra-chess/pkg"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	board := NewBoard()
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "uci":
			fmt.Println("id name LibraChess")
			fmt.Println("id author eugenioenko")
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "ucinewgame":
			board = NewBoard()
			board.LoadInitial()
		case "position":
			board.ParseAndApplyPosition(fields[1:])
		case "go":
			tt := NewTranspositionTable()
			depth := 4
			material := board.CountPieces()
			if material < 20 {
				depth = 4
			}
			if board.OnlyKingLeft() {
				depth = 4
			}
			score, move := board.Search(depth, tt)
			fmt.Printf("info score cp %d\n", score)
			if move != nil {
				fmt.Printf("bestmove %s\n", move.ToUCI())
			} else {
				fmt.Println("bestmove 0000")
			}
		case "quit":
			return
		}
	}
}
