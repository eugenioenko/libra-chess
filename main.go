package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	. "github.com/eugenioenko/libra-chess/pkg"
)

func main() {
	fmt.Println("Welcome to LibraChess v1.0.1!")
	fmt.Println("Ready to play? Type 'uci' to begin your chess adventure!")
	fmt.Println("Type 'quit' to exit the CLI at any time.")
	fmt.Println("LibraChess is a UCI chess engine, designed to be used with a chess GUI (like CuteChess, CoreChess, PyChess, etc.)")
	fmt.Println("For more information, visit: https://github.com/eugenioenko/libra-chess")
	scanner := bufio.NewScanner(os.Stdin)
	board := NewBoard()

	var searchMu sync.Mutex
	var stopChan chan struct{}

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
			goOpts := ParseGoOptions(fields)
			optimalTime, maxTime := goOpts.CalcTimeLimit(board.WhiteToMove)

			searchMu.Lock()
			stopChan = make(chan struct{})
			currentStop := stopChan
			searchMu.Unlock()

			opts := SearchOptions{
				TimeLimitInMs:    optimalTime,
				MaxTimeLimitInMs: maxTime,
				MaxDepth:         goOpts.Depth,
				StopChan:         currentStop,
			}

			go func() {
				bestMove := board.IterativeDeepeningSearch(opts)
				if bestMove != nil {
					fmt.Printf("bestmove %s\n", bestMove.ToUCI())
				} else {
					fmt.Println("bestmove 0000")
				}
			}()
		case "stop":
			searchMu.Lock()
			if stopChan != nil {
				close(stopChan)
				stopChan = nil
			}
			searchMu.Unlock()
		case "quit":
			searchMu.Lock()
			if stopChan != nil {
				close(stopChan)
				stopChan = nil
			}
			searchMu.Unlock()
			return
		}
	}
}
