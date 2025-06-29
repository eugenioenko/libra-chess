package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/eugenioenko/libra-chess/pkg"
)

const (
	SearchMaxDepth = 16
)

func main() {
	fmt.Println("Welcome to LibraChess v1.0.1!")
	fmt.Println("Ready to play? Type 'uci' to begin your chess adventure!")
	fmt.Println("Type 'quit' to exit the CLI at any time.")
	fmt.Println("LibraChess is a UCI chess engine, designed to be used with a chess GUI (like CuteChess, CoreChess, PyChess, etc.)")
	fmt.Println("For more information, visit: https://github.com/eugenioenko/libra-chess")
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
			move, stats := board.Search(BaseSearchDepth, tt)
			stats.PrintUCI()
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
