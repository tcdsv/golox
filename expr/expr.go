package expr

import "golox/scanner"

type Visitor[T any] interface {
	VisitLiteral(element LiteralExpr) T
	VisitUnary(element UnaryExpr) T
	VisitBinary(element BinaryExpr) T
	VisitGrouping(element GroupingExpr) T
}

type Expr interface {
	AcceptPrinter(visitor Visitor[string]) string
}

type LiteralExpr struct {
	Value interface{}
}

func (l LiteralExpr) AcceptPrinter(visitor Visitor[string]) string {
	return visitor.VisitLiteral(l)
}

type UnaryExpr struct {
	Operator scanner.Token
	Expr     Expr
}

func (l UnaryExpr) AcceptPrinter(visitor Visitor[string]) string {
	return visitor.VisitUnary(l)
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) AcceptPrinter(visitor Visitor[string]) string {
	return visitor.VisitGrouping(e)
}

type BinaryExpr struct {
	Operator scanner.Token
	Left     Expr
	Right    Expr
}

func (e BinaryExpr) AcceptPrinter(visitor Visitor[string]) string {
	return visitor.VisitBinary(e)
}
