package expr

import (
	tkn "golox/token"
	loxvalue "golox/value"
)

type ExprVisitor interface {
	VisitLiteral(element LiteralExpr) (interface{}, error)
	VisitUnary(element UnaryExpr) (interface{}, error)
	VisitBinary(element BinaryExpr) (interface{}, error)
	VisitGrouping(element GroupingExpr) (interface{}, error)
	// VisitVariable(element VariableExpr) (interface{}, error)
	// VisitAssing(element AssignExpr) (interface{}, error)
	VisitLogical(element LogicalExpr) (interface{}, error)
}

type Expr interface {
	Evaluate(visitor ExprVisitor) (interface{}, error)
}

type LiteralExpr struct {
	Value loxvalue.LoxValue
}

func (e LiteralExpr) Evaluate(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator tkn.Token
	Right     Expr
}

func (e UnaryExpr) Evaluate(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnary(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) Evaluate(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGrouping(e)
}

type BinaryExpr struct {
	Operator tkn.Token
	Left     Expr
	Right    Expr
}

func (e BinaryExpr) Evaluate(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinary(e)
}

type VariableExpr struct {
	Name tkn.Token
}

func (e VariableExpr) Evaluate(visitor ExprVisitor) (interface{}, error) {
	return nil, nil
	// return visitor.VisitVariable(e)
}

type AssignExpr struct {
	Name tkn.Token
	Right Expr
}

func (e AssignExpr) Evaluate(visitor ExprVisitor) (interface{}, error) {
	return nil, nil
	// return visitor.VisitAssing(e)
}

type LogicalExpr struct {
	Operator tkn.Token
	Left Expr
	Right Expr
}

func (e LogicalExpr) Evaluate(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLogical(e)
}