package expr

import "golox/scanner"

type Visitor interface {
	VisitLiteral(element LiteralExpr)
	VisitUnary(element UnaryExpr)
	VisitBinary(element BinaryExpr)
	VisitGrouping(element GroupingExpr)
}

type Expr interface {
	Accept(visitor Visitor) 
}

type LiteralExpr struct {
	Value interface{}
}

func (e LiteralExpr) Accept(visitor Visitor) {
	visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator scanner.Token
	Expr     Expr
}

func (e UnaryExpr) Accept(visitor Visitor) {
	visitor.VisitUnary(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) Accept(visitor Visitor) {
	visitor.VisitGrouping(e)
}

type BinaryExpr struct {
	Operator scanner.Token
	Left     Expr
	Right    Expr
}

func (e BinaryExpr) Accept(visitor Visitor) {
	visitor.VisitBinary(e)
}
