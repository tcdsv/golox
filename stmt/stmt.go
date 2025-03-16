package stmt

import (
	"golox/expr"
	"golox/token"
)

type StmtVisitor interface {
	VisitExpressionStatement(exprStmt ExprStmt) (interface{}, error)
	VisitPrintStatement(printStmt PrintStmt) (interface{}, error)
	VisitVariableStatement(variableStmt VarStmt) (interface{}, error)
	VisitBlockStatement(BlockStmt BlockStmt) (interface{}, error)
	VisitIfStatement(IfStmt IfStmt) (interface{}, error)
	VisitWhileStatement(WhileStmt WhileStmt) (interface{}, error)
}

type Stmt interface {
	Accept(visitor StmtVisitor) (interface{}, error)
}

type ExprStmt struct {
	E expr.Expr
}

func (es ExprStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExpressionStatement(es)
}

type PrintStmt struct {
	E expr.Expr
}

func (ps PrintStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return nil, nil
	// return visitor.VisitPrintStatement(ps)
}

type VarStmt struct {
	Name token.Token
	Initializer expr.Expr
}

func (s VarStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return nil, nil
	// return visitor.VisitVariableStatement(s)
}

type BlockStmt struct {
	Statements []Stmt
}

func (s BlockStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return nil, nil
	// return visitor.VisitBlockStatement(s)
}

type IfStmt struct {
	Condition expr.Expr
	ThenBrnach Stmt
	ElseBranch Stmt
}

func (s IfStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return nil, nil
	// return visitor.VisitIfStatement(s)
}

type WhileStmt struct { 
	Condition expr.Expr
	Body Stmt
}

func (s WhileStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return nil, nil
	// return visitor.VisitWhileStatement(s)
}