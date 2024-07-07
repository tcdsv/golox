package loxprinter

import (
	"fmt"
	"golox/expr"
	"strconv"
	"strings"
)

type AstPrinter struct{
	result string
}

func (a *AstPrinter) Print(expr expr.Expr) {
	expr.Accept(a)
}

func (a *AstPrinter) VisitLiteral(expr expr.LiteralExpr) {
	if expr.Value == nil {
		a.result = "nil"
		return
	}

	switch v := expr.Value.(type) {
	case int:
		a.result = strconv.Itoa(v)
	case string:
		a.result = v
	case float64:
		a.result = fmt.Sprintf("%f", v)
	default:
		a.result = "error - unknown literal value"
		// panic(fmt.Sprintf("Unsupported literal type: %T", v))
	}
}

func (a *AstPrinter) VisitUnary(expr expr.UnaryExpr) {
	a.result = a.parenthesize(expr.Operator.Lexeme, expr.Expr)
}

func (a *AstPrinter) VisitBinary(expr expr.BinaryExpr) {
	a.result = a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGrouping(expr expr.GroupingExpr) {
	a.result = a.parenthesize("group", expr.Expr)
}

func (a AstPrinter) parenthesize(name string, exprs ...expr.Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		a.Print(expr)
		builder.WriteString(a.result)
	}

	builder.WriteString(")")

	return builder.String()
}
