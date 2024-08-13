package interpreter

import (
	"fmt"
	loxerror "golox/error"
	"golox/expr"
	expression "golox/expr"
	"golox/stmt"
	tkn "golox/token"
	loxvalue "golox/value"
	visitor "golox/visitor"
)

type Interpreter struct {
	Results []*visitor.VisitorResult
}

func NewInterpreter() *Interpreter {
	return &Interpreter{Results: []*visitor.VisitorResult{},}
}

func (i *Interpreter) evaluate(expr expression.Expr) *visitor.VisitorResult {
	return expr.Accept(i)
}

func (i *Interpreter) Interpret(statements []stmt.Stmt) {
	for _, statement := range statements {
		result := statement.Accept(i)
		i.Results = append(i.Results, result)
	}
}

func (i *Interpreter) VisitPrintStatement(printStmt stmt.PrintStmt) *visitor.VisitorResult {
	lv := getLoxVaule(i.evaluate(printStmt.E))
	fmt.Println(lv.ToString())
	return visitor.NewVisitorResult(nil, nil)
}

func (i *Interpreter) VisitExpressionStatement(exprStmt stmt.ExprStmt) *visitor.VisitorResult {
	return i.evaluate(exprStmt.E)
}

func (i *Interpreter) VisitLiteral(expr expression.LiteralExpr) *visitor.VisitorResult {
	return visitor.NewVisitorResult(expr.Value, nil)
}

func (i *Interpreter) VisitUnary(expr expr.UnaryExpr) *visitor.VisitorResult {
	right := getLoxVaule(i.evaluate(expr.Right))

	var result loxvalue.LoxValue

	switch (expr.Operator.Type) {
	case tkn.MINUS:
		number, err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		result = number.Minus()
	case tkn.BANG:
		result := loxvalue.NewBoolean(!loxvalue.IsTruthy(right))
		return visitor.NewVisitorResult(result, nil)
	}
	//todo: error
	return visitor.NewVisitorResult(result, nil)
}

func getLoxVaule(r *visitor.VisitorResult) loxvalue.LoxValue {
	return r.Result.(loxvalue.LoxValue)
}

func (i *Interpreter) VisitBinary(expr expr.BinaryExpr) *visitor.VisitorResult {
	left := getLoxVaule(i.evaluate(expr.Left))
	right := getLoxVaule(i.evaluate(expr.Right))
	//todo: error checking

	switch (expr.Operator.Type) {
	case tkn.MINUS:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(left.Subtract(right), nil)

	case tkn.SLASH:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(left.Divide(right), nil)

	case tkn.STAR:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(left.Multiply(right), nil)

	case tkn.PLUS:
		
		result, err := binaryPlus(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(result, nil)

	case tkn.GREATER:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(left.Greater(right), nil)

	case tkn.GREATER_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(left.GreaterEqual(right), nil)

	case tkn.LESS:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(left.Less(right), nil)

	case tkn.LESS_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return visitor.NewVisitorResult(nil, err)
		}
		return visitor.NewVisitorResult(left.LessEqual(right), nil)
	
	case tkn.BANG_EQUAL:
		
		result := loxvalue.NewBoolean(!loxvalue.IsEqual(left, right))
		return visitor.NewVisitorResult(result, nil)

	case tkn.EQUAL_EQUAL:

		result := loxvalue.NewBoolean(loxvalue.IsEqual(left, right))
		return visitor.NewVisitorResult(result, nil)

	}

	//todo: error
	return nil
}

func (i *Interpreter) VisitGrouping(expr expr.GroupingExpr) *visitor.VisitorResult {
	return i.evaluate(expr.Expr)
}

func binaryPlus(operator tkn.Token, left loxvalue.LoxValue, right loxvalue.LoxValue) (loxvalue.LoxValue, error) {
	
	if left.Type() == loxvalue.NUMBER && right.Type() == loxvalue.NUMBER {
		leftNum := left.(*loxvalue.Number)
		rightNum := right.(*loxvalue.Number)
		return leftNum.Add(rightNum), nil
	}

	if left.Type() == loxvalue.STRING && right.Type() == loxvalue.STRING {
		leftStr := left.(*loxvalue.String)
		rightStr := right.(*loxvalue.String)
		return leftStr.Concat(rightStr), nil
	}
	
	return nil, loxerror.NewErrorFromToken(operator, "Operands must be two numbers or two strings.")

}

func checkNumberOperand(operator tkn.Token, v loxvalue.LoxValue) (*loxvalue.Number, error) {
	if v.Type() == loxvalue.NUMBER {
		return v.(*loxvalue.Number), nil
	}
	return nil, loxerror.NewErrorFromToken(operator, "Operand must be a number.")
}

func checkNumberOperands(operator tkn.Token, left loxvalue.LoxValue, right loxvalue.LoxValue) (*loxvalue.Number, *loxvalue.Number, error) {
	if left.Type() == loxvalue.NUMBER && right.Type() == loxvalue.NUMBER {
		return left.(*loxvalue.Number), right.(*loxvalue.Number), nil
	}
	return nil, nil, loxerror.NewErrorFromToken(operator, "Operands must be a numbers.")
}