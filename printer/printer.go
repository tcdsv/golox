package loxprinter

import (
	"golox/expr"
	"strconv"
	"strings"
)

type AstPrinter struct{}

func (a AstPrinter) Print(expr expr.Expr) string {
	return expr.AcceptPrinter(a)
}

func (a AstPrinter) VisitLiteral(expr expr.LiteralExpr) string {
	v, _ := expr.Value.(int)
	return strconv.Itoa(v)
}

func (a AstPrinter) VisitUnary(expr expr.UnaryExpr) string {
	return a.parenthesize(expr.Token.Text, expr.Expr)
}

func (a AstPrinter) parenthesize(name string, exprs ...expr.Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(a.Print(expr))
	}

	builder.WriteString(")")

	return builder.String()
}
