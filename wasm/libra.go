//go:build js && wasm

package main

import (
	"syscall/js"

	. "github.com/eugenioenko/libra-chess/pkg"
)

var board *Board

func jsNewBoard(this js.Value, args []js.Value) interface{} {
	board = NewBoard()
	return js.ValueOf(true)
}

func jsLoadInitial(this js.Value, args []js.Value) interface{} {
	if board == nil {
		board = NewBoard()
	}
	board.LoadInitial()
	return js.ValueOf(true)
}

func jsFromFEN(this js.Value, args []js.Value) interface{} {
	if board == nil {
		board = NewBoard()
	}
	fen := args[0].String()
	ok, err := board.FromFEN(fen)
	if !ok || err != nil {
		return js.ValueOf(err.Error())
	}
	return js.ValueOf(true)
}

func jsIterativeDeepeningSearch(this js.Value, args []js.Value) interface{} {
	if board == nil {
		return js.ValueOf("")
	}
	ms := 1500
	if len(args) > 0 {
		ms = args[0].Int()
	}
	move := board.IterativeDeepeningSearch(SearchOptions{TimeDepthLimitInMs: ms})
	if move == nil {
		return js.ValueOf("")
	}
	return js.ValueOf(move.ToUCI())
}

func jsToFEN(this js.Value, args []js.Value) interface{} {
	if board == nil {
		return js.ValueOf("")
	}
	return js.ValueOf(board.ToFEN())
}

func jsMove(this js.Value, args []js.Value) interface{} {
	if board == nil {
		return js.ValueOf(false)
	}
	uci := args[0].String()
	move := board.ParseUCIMove(uci)
	if move == nil {
		return js.ValueOf(false)
	}
	board.Move(*move)
	return js.ValueOf(true)
}

func jsPerftParallel(this js.Value, args []js.Value) interface{} {
	if board == nil {
		return js.ValueOf(-1)
	}
	if len(args) == 0 {
		return js.ValueOf(-1)
	}
	depth := args[0].Int()
	// Use PerftParallel for better performance if available
	nodes := board.PerftParallel(depth)
	return js.ValueOf(nodes)
}

func jsPerft(this js.Value, args []js.Value) interface{} {
	if board == nil {
		return js.ValueOf(-1)
	}
	if len(args) == 0 {
		return js.ValueOf(-1)
	}
	depth := args[0].Int()
	// Use PerftParallel for better performance if available
	nodes := board.Perft(depth)
	return js.ValueOf(nodes)
}

func registerCallbacks() {
	js.Global().Set("libraNewBoard", js.FuncOf(jsNewBoard))
	js.Global().Set("libraLoadInitial", js.FuncOf(jsLoadInitial))
	js.Global().Set("libraFromFEN", js.FuncOf(jsFromFEN))
	js.Global().Set("libraIterativeDeepeningSearch", js.FuncOf(jsIterativeDeepeningSearch))
	js.Global().Set("libraToFEN", js.FuncOf(jsToFEN))
	js.Global().Set("libraMove", js.FuncOf(jsMove))
	js.Global().Set("libraPerft", js.FuncOf(jsPerft))
	js.Global().Set("libraPerftParallel", js.FuncOf(jsPerftParallel))
}

func main() {
	board = NewBoard()
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
