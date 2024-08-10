package stmt

import (
	"golox/expr"
	"golox/visitor"
)

type StmtVisitor interface {
	visitExpressionStmt(exprStmt ExprStmt) *visitor.VisitorResult
	visitPrintStmt(printStmt PrintStmt) *visitor.VisitorResult
}

type Stmt interface {
	Accept(visitor StmtVisitor) *visitor.VisitorResult
}

type ExprStmt struct {
	E expr.Expr
}

func (es ExprStmt) Accept(visitor StmtVisitor) *visitor.VisitorResult {
	return visitor.visitExpressionStmt(es)
}

type PrintStmt struct {
	E expr.Expr
}

func (ps PrintStmt) Accept(visitor StmtVisitor) *visitor.VisitorResult {
	return visitor.visitPrintStmt(ps)
}