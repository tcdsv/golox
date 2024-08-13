package stmt

import (
	"golox/expr"
	"golox/token"
	"golox/visitor"
)

type StmtVisitor interface {
	VisitExpressionStatement(exprStmt ExprStmt) *visitor.VisitorResult
	VisitPrintStatement(printStmt PrintStmt) *visitor.VisitorResult
	VisitVariableStatement(variableStmt VarStmt) *visitor.VisitorResult
}

type Stmt interface {
	Accept(visitor StmtVisitor) *visitor.VisitorResult
}

type ExprStmt struct {
	E expr.Expr
}

func (es ExprStmt) Accept(visitor StmtVisitor) *visitor.VisitorResult {
	return visitor.VisitExpressionStatement(es)
}

type PrintStmt struct {
	E expr.Expr
}

func (ps PrintStmt) Accept(visitor StmtVisitor) *visitor.VisitorResult {
	return visitor.VisitPrintStatement(ps)
}

type VarStmt struct {
	Name token.Token
	Initializer expr.Expr
}

func (s VarStmt) Accept(visitor StmtVisitor) *visitor.VisitorResult {
	return visitor.VisitVariableStatement(s)
}