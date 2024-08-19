package expr

import (
	tkn "golox/token"
	loxvalue "golox/value"
)

type ExprVisitor[T any] interface {
	VisitLiteral(element LiteralExpr) T
	VisitUnary(element UnaryExpr) T
	VisitBinary(element BinaryExpr) T
	VisitGrouping(element GroupingExpr) T
	VisitVariable(element VariableExpr) T
	VisitAssing(element AssignExpr) T
}

type LoxValueResult struct {
	Value loxvalue.LoxValue
	Error error
}

type Expr interface {
	Evaluate(visitor ExprVisitor[LoxValueResult]) LoxValueResult
}

type LiteralExpr struct {
	Value loxvalue.LoxValue
}

func (e LiteralExpr) Evaluate(visitor ExprVisitor[LoxValueResult]) LoxValueResult {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator tkn.Token
	Right     Expr
}

func (e UnaryExpr) Evaluate(visitor ExprVisitor[LoxValueResult]) LoxValueResult {
	return visitor.VisitUnary(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) Evaluate(visitor ExprVisitor[LoxValueResult]) LoxValueResult {
	return visitor.VisitGrouping(e)
}

type BinaryExpr struct {
	Operator tkn.Token
	Left     Expr
	Right    Expr
}

func (e BinaryExpr) Evaluate(visitor ExprVisitor[LoxValueResult]) LoxValueResult {
	return visitor.VisitBinary(e)
}

type VariableExpr struct {
	Name tkn.Token
}

func (e VariableExpr) Evaluate(visitor ExprVisitor[LoxValueResult]) LoxValueResult {
	return visitor.VisitVariable(e)
}

type AssignExpr struct {
	Name tkn.Token
	Right Expr
}

func (e AssignExpr) Evaluate(visitor ExprVisitor[LoxValueResult]) LoxValueResult {
	return visitor.VisitAssing(e)
}