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
	tt := NewTranspositionTable()
	move := board.IterativeDeepeningSearch(SearchOptions{
		TimeLimitInMs:      ms,
		TranspositionTable: tt,
	})
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
	nodes := board.Perft(depth)
	return js.ValueOf(nodes)
}

func registerCallbacks() {
	libra := js.Global().Get("Object").New()
	libra.Set("version", "1.0.1")
	libra.Set("newBoard", js.FuncOf(jsNewBoard))
	libra.Set("loadInitial", js.FuncOf(jsLoadInitial))
	libra.Set("fromFEN", js.FuncOf(jsFromFEN))
	libra.Set("iterativeDeepeningSearch", js.FuncOf(jsIterativeDeepeningSearch))
	libra.Set("toFEN", js.FuncOf(jsToFEN))
	libra.Set("move", js.FuncOf(jsMove))
	libra.Set("perft", js.FuncOf(jsPerft))
	libra.Set("perftParallel", js.FuncOf(jsPerftParallel))
	js.Global().Set("libra", libra)
}

func main() {
	board = NewBoard()
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
