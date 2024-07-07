package main

import (
	"fmt"
	"golox/parser"
	"golox/scanner"
)

func main() {

	source := "1"
	scanner := scanner.NewScanner(source)
	tokens, _ := scanner.Scan()
	parser := parser.NewParser(tokens)
	res := parser.Parse()
	fmt.Println(res)
}

func runPrompt() error {
	for {
		fmt.Println("> ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			return err
		}
		if len(input) == 0 {
			break
		}
		run(input)
	}
	return nil
}

func run(source string) {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.Scan()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(tokens)
	}
}
