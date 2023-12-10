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

/*
  public static void main(String[] args) {
    Expr expression = new Expr.Binary(
        new Expr.Unary(
            new Token(TokenType.MINUS, "-", null, 1),
            new Expr.Literal(123)),
        new Token(TokenType.STAR, "*", null, 1),
        new Expr.Grouping(
            new Expr.Literal(45.67)));

    System.out.println(new AstPrinter().print(expression));
  }

*/

func main() {

	b := expr.BinaryExpr{
		Operator: scanner.NewToken(scanner.STAR, nil, "*", 1),
		Left: expr.UnaryExpr{
			Operator: scanner.NewToken(scanner.MINUS, nil, "-", 1),
			Expr: expr.LiteralExpr{
				Value: 123,
			},
		},
		Right: expr.GroupingExpr{
			Expr: expr.LiteralExpr{
				Value: 45.67,
			},
		},
	}

	p := loxprinter.AstPrinter{}
	fmt.Println("output:" + p.Print(b))

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
