package interpreter

import (
	loxerror "golox/error"
	"golox/expr"
	expression "golox/expr"
	tkn "golox/token"
	loxvalue "golox/value"
)

type Interpreter struct {
	Value loxvalue.LoxValue
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) evaluate(expr expression.Expr) *expr.VisitorResult {
	return expr.Accept(i)
}

func (i *Interpreter) Interpret(expr expression.Expr) {
	result := expr.Accept(i)
	i.Value = result.Result.(loxvalue.LoxValue)
}

func (i *Interpreter) VisitLiteral(expr expression.LiteralExpr) *expression.VisitorResult {
	return expression.NewVisitorResult(expr.Value, nil)
}

func (i *Interpreter) VisitUnary(expr expr.UnaryExpr) *expression.VisitorResult {
	right := getLoxVaule(i.evaluate(expr.Right))

	var result loxvalue.LoxValue

	switch (expr.Operator.Type) {
	case tkn.MINUS:
		number, err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		result = number.Minus()
	case tkn.BANG:
		result := loxvalue.NewBoolean(!loxvalue.IsTruthy(right))
		return expression.NewVisitorResult(result, nil)
	}
	//todo: error
	return expression.NewVisitorResult(result, nil)
}

func getLoxVaule(r *expression.VisitorResult) loxvalue.LoxValue {
	return r.Result.(loxvalue.LoxValue)
}

func (i *Interpreter) VisitBinary(expr expr.BinaryExpr) *expression.VisitorResult {
	left := getLoxVaule(i.evaluate(expr.Left))
	right := getLoxVaule(i.evaluate(expr.Right))
	//todo: error checking

	switch (expr.Operator.Type) {
	case tkn.MINUS:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(left.Subtract(right), nil)

	case tkn.SLASH:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(left.Divide(right), nil)

	case tkn.STAR:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(left.Multiply(right), nil)

	case tkn.PLUS:
		
		result, err := binaryPlus(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(result, nil)

	case tkn.GREATER:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(left.Greater(right), nil)

	case tkn.GREATER_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(left.GreaterEqual(right), nil)

	case tkn.LESS:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(left.Less(right), nil)

	case tkn.LESS_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return expression.NewVisitorResult(nil, err)
		}
		return expression.NewVisitorResult(left.LessEqual(right), nil)
	
	case tkn.BANG_EQUAL:
		
		result := loxvalue.NewBoolean(!loxvalue.IsEqual(left, right))
		return expression.NewVisitorResult(result, nil)

	case tkn.EQUAL_EQUAL:

		result := loxvalue.NewBoolean(loxvalue.IsEqual(left, right))
		return expression.NewVisitorResult(result, nil)

	}

	//todo: error
	return nil
}

func (i *Interpreter) VisitGrouping(expr expr.GroupingExpr) *expression.VisitorResult {
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