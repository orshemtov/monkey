package main

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		lexer := NewLexer(line)
		for token := lexer.NextToken(); token.Type != EOF; token = lexer.NextToken() {
			fmt.Printf("%+v\n", token)
		}
	}
}
