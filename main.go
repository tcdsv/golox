package main

import (
	"bufio"
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
	var err error
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
		err = runPrompt()
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		run(line)
	}
}

func run(source string) {
	scanner := scanner.NewScanner(source)
	tokens, _ := scanner.Scan()
	fmt.Println(tokens)
}
