package main

import (
	"fmt"
	"golox/expr"
	loxprinter "golox/printer"
	"golox/scanner"
)

type Error struct {
	Line    int
	Where   string
	Message int
}

/*
   new Expr.Unary(
       new Token(TokenType.MINUS, "-", null, 1),
       new Expr.Literal(123)),
*/

func main() {

	t := scanner.NewToken(scanner.MINUS, nil, "-", 1)

	le := expr.LiteralExpr{Value: 123}

	unary := expr.UnaryExpr{
		Token: t,
		Expr:  le,
	}

	p := loxprinter.AstPrinter{}
	fmt.Println("output:" + p.Print(unary))

	/*if len(os.Args) > 1 {
		fmt.Println("Useage: jlox [script]")
		os.Exit(64)
	} else if len(os.Args) == 1 {
		//runFile
	}
	runPrompt()*/
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

func report(err Error) {
	//print error
}
