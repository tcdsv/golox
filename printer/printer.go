package loxprinter

import (
	"fmt"
	"golox/expr"
	"strconv"
	"strings"
)

type AstPrinter struct{}

func (a AstPrinter) Print(expr expr.Expr) string {
	return expr.AcceptPrinter(a)
}

func (a AstPrinter) VisitLiteral(expr expr.LiteralExpr) string {
	if expr.Value == nil {
		return "nil"
	}

	switch v := expr.Value.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case float64:
		return fmt.Sprintf("%f", v)

	default:
		return "unknown"
		// panic(fmt.Sprintf("Unsupported literal type: %T", v))
	}
}

func (a AstPrinter) VisitUnary(expr expr.UnaryExpr) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Expr)
}

func (a AstPrinter) VisitBinary(expr expr.BinaryExpr) string {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a AstPrinter) VisitGrouping(expr expr.GroupingExpr) string {
	return a.parenthesize("group", expr.Expr)
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
