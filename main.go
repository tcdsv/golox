package main

import (
	"fmt"
	"golox/scanner"
	"os"
)

// Define exit codes based on sysexits.h
const (
	EX_OK          = 0   // successful termination
	EX_USAGE       = 64  // command line usage error
)

func main() {

	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: golox [script]")
		os.Exit(EX_USAGE)
	}
	if len(args) == 1 {
		fmt.Println("Loading from file is not implemented")
		os.Exit(EX_USAGE)
		// runFile()
	} else {
		runPrompt()
	}
	
}

/*func runFile() {
}*/

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
	tokens, _ := scanner.Scan()
	fmt.Println(tokens)
}
