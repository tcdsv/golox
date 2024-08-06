package expr

import "golox/scanner"

type VisitorResult struct {
	Result interface{}
	Err error
}

func NewVisitorResult(result interface{}, err error) *VisitorResult {
	return &VisitorResult{
		Result: result,
		Err: err,
	}
}

type Visitor interface {
	VisitLiteral(element LiteralExpr) *VisitorResult
	VisitUnary(element UnaryExpr) *VisitorResult
	VisitBinary(element BinaryExpr) *VisitorResult
	VisitGrouping(element GroupingExpr) *VisitorResult
}

type Expr interface {
	Accept(visitor Visitor) *VisitorResult
}

type LiteralExpr struct {
	Value interface{}
}

func (e LiteralExpr) Accept(visitor Visitor) *VisitorResult {
	return visitor.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right     Expr
}

func (e UnaryExpr) Accept(visitor Visitor) *VisitorResult {
	return visitor.VisitUnary(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) Accept(visitor Visitor) *VisitorResult {
	return visitor.VisitGrouping(e)
}

type BinaryExpr struct {
	Operator scanner.Token
	Left     Expr
	Right    Expr
}

func (e BinaryExpr) Accept(visitor Visitor) *VisitorResult {
	return visitor.VisitBinary(e)
}
