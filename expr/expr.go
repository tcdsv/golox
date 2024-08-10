package expr

import (
	tkn "golox/token"
	loxvalue "golox/value"
	visitor "golox/visitor"
)

type ExprVisitor interface {
	VisitLiteral(element LiteralExpr) *visitor.VisitorResult
	VisitUnary(element UnaryExpr) *visitor.VisitorResult
	VisitBinary(element BinaryExpr) *visitor.VisitorResult
	VisitGrouping(element GroupingExpr) *visitor.VisitorResult
}

type Expr interface {
	Accept(visitor ExprVisitor) *visitor.VisitorResult
}

type LiteralExpr struct {
	Value loxvalue.LoxValue
}

func (e LiteralExpr) Accept(visitor ExprVisitor) *visitor.VisitorResult {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator tkn.Token
	Right     Expr
}

func (e UnaryExpr) Accept(visitor ExprVisitor) *visitor.VisitorResult {
	return visitor.VisitUnary(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) Accept(visitor ExprVisitor) *visitor.VisitorResult {
	return visitor.VisitGrouping(e)
}

type BinaryExpr struct {
	Operator tkn.Token
	Left     Expr
	Right    Expr
}

func (e BinaryExpr) Accept(visitor ExprVisitor) *visitor.VisitorResult {
	return visitor.VisitBinary(e)
}
