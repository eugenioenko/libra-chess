package main

import (
	"bufio"
	"fmt"
	. "libra/pkg"
	"math/rand"
	"os"
	"strings"
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
			fmt.Println("id author YourName")
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "ucinewgame":
			board = NewBoard()
			board.LoadInitial()
		case "position":
			board.ParseAndApplyPosition(fields[1:])
		case "go":
			board.GenerateLegalMoves()
			fmt.Println("info score cp 0")
			if len(board.Moves) > 0 {
				randomMoveIndex := rand.Intn(len(board.Moves))
				move := board.Moves[randomMoveIndex]
				fmt.Printf("bestmove %s\n", move.ToUCI())
			} else {
				fmt.Println("bestmove 0000")
			}
		case "quit":
			return
		}
	}
}
