package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "github.com/eugenioenko/libra-chess/pkg"
)

// Takes a list of FEN position and their best moves and creates a JSON file
// with Zobrist hashes as keys and lists of moves as values.
func main() {
	board := NewBoard()
	file, _ := os.Open("books/book.txt")
	defer file.Close()

	outFile, _ := os.Create("books/book.json")
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)

	writer.WriteString("{\n")

	firstMove := true
	firstPosition := true
	for scanner.Scan() {
		line := scanner.Text()
		output := ""
		if strings.HasPrefix(line, "pos ") {
			if !firstPosition {
				output += "],\n"
			}
			board.FromFEN(line[4:])
			hash := board.ZobristHashWasm()
			output += fmt.Sprintf("\t\"%#x\": [", hash)
			firstMove = true
			firstPosition = false
		} else {
			if !firstMove {
				output += ", "
			}
			parts := strings.Fields(line)
			output += fmt.Sprintf("\"%s\"", parts[0])
			firstMove = false
		}
		writer.WriteString(output)
	}
	writer.WriteString("],\n}\n")
	if err := writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing to Parsed.txt: %v\n", err)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading Book.txt: %v\n", err)
	}
}
