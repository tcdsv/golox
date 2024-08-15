package main

import (
	"bufio"
	"fmt"
	"golox/interpreter"
	"golox/parser"
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
	interpreter := interpreter.NewInterpreter()
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		Run(line, interpreter)
	}
}

func Run(source string, interpreter *interpreter.Interpreter) {

	hasError := false
	scanner := scanner.NewScanner(source)
	tokens, errors := scanner.Scan()
	if len(errors) > 0 {
		printErrors(errors)
		hasError = true
	}
	parser := parser.NewParser(tokens)
	statements, errors := parser.Parse()
	if len(errors) > 0 {
		printErrors(errors)
		hasError = true
	}
	if hasError {
		return
	}
	interpreter.Interpret(statements)
	for _, statement := range interpreter.Results {
		if statement.Err != nil {
			fmt.Println(statement.Err.Error())
		}
	}
}

func printErrors(errors []error) {
	for _, err := range errors {
		fmt.Println(err.Error())
	}
}