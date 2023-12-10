package expr

import "golox/scanner"

type Visitor[T any] interface {
	VisitLiteral(element LiteralExpr) T
	VisitUnary(element UnaryExpr) T
}

type Expr interface {
	AcceptPrinter(visitor Visitor[string]) string
}

// ---------------------

type LiteralExpr struct {
	Value interface{}
}

func (l LiteralExpr) AcceptPrinter(visitor Visitor[string]) string {
	return visitor.VisitLiteral(l)
}

// ---------------------

type UnaryExpr struct {
	Token scanner.Token
	Expr  Expr
}

func (l UnaryExpr) AcceptPrinter(visitor Visitor[string]) string {
	return visitor.VisitUnary(l)
}

// ---------------------
