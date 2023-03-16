package testing

import (
	"fmt"
	libra "libra/pkg"
	"testing"
)

func TestShouldInsertValues(t *testing.T) {
	board := libra.NewBoard()
	ok, error := board.LoadFromFEN(libra.BoardInitialFEN)

	if error != nil {
		fmt.Println(error)
	}

	if ok != true {
		t.Fail()
	}
}
