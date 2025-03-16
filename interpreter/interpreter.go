package interpreter

import (
	loxerror "golox/error"
	"golox/expr"
	expression "golox/expr"
	"golox/stmt"
	tkn "golox/token"
	loxvalue "golox/value"
)

type Interpreter struct {
	env     *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env:     NewGlobalEnv(),
	}
}

// func (i *Interpreter) Interpret(statements []stmt.Stmt) {
// 	for _, statement := range statements {
// 		_, err := i.execute(statement)
// 		if err != nil {
// 			return 
// 		}
// 	}
// }

func (i *Interpreter) VisitLiteral(expr expression.LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitLogical(LogicalExpr expr.LogicalExpr) (interface{}, error) {

	left, err := i.Evaluate(LogicalExpr.Left)
	if err != nil {
		return nil, err
	}

	leftValue, ok := left.(loxvalue.LoxValue)
	if !ok {
		return nil, err //todo: define error message
	}

	operator := LogicalExpr.Operator.Type
	switch operator {
	case tkn.OR:

		if loxvalue.IsTruthy(leftValue) {
			return leftValue, nil
		}

	case tkn.AND:

		if !loxvalue.IsTruthy(leftValue) {
			return leftValue, nil
		}
	}
	return i.Evaluate(LogicalExpr.Right)

}

func (i *Interpreter) Evaluate(expr expression.Expr) (loxvalue.LoxValue, error) {
	
	value, err := expr.Evaluate(i)
	if err != nil {
		return nil, err
	}

	loxValue, ok := value.(loxvalue.LoxValue)
	if !ok {
		return nil, nil //todo casting error
	}

	return loxValue, nil
}

func (i *Interpreter) VisitExpressionStatement(exprStmt stmt.ExprStmt) (interface{}, error) {
	return i.Evaluate(exprStmt.E)
}

func (i *Interpreter) VisitUnary(expr expr.UnaryExpr) (interface{}, error) {

	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	var result loxvalue.LoxValue

	switch expr.Operator.Type {
	case tkn.MINUS:
		number, err := checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		result = number.Minus()
	case tkn.BANG:
		result = loxvalue.NewBoolean(!loxvalue.IsTruthy(right))
	}
	//todo: error
	return result, nil

}

func (i *Interpreter) VisitBinary(expr expr.BinaryExpr) (interface{}, error) {

	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case tkn.MINUS:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.Subtract(right), nil

	case tkn.SLASH:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.Divide(right), nil

	case tkn.STAR:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.Multiply(right), nil

	case tkn.PLUS:

		result, err := binaryPlus(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return result, err

	case tkn.GREATER:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.Greater(right), nil

	case tkn.GREATER_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.GreaterEqual(right), nil

	case tkn.LESS:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.Less(right), nil

	case tkn.LESS_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.LessEqual(right), nil

	case tkn.BANG_EQUAL:

		result := loxvalue.NewBoolean(!loxvalue.IsEqual(left, right))
		return result, nil

	case tkn.EQUAL_EQUAL:

		result := loxvalue.NewBoolean(loxvalue.IsEqual(left, right))
		return result, nil

	}

	//todo: error
	return nil, nil
}

func (i *Interpreter) VisitGrouping(expr expr.GroupingExpr) (interface{}, error) {
	return i.Evaluate(expr.Expr)
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
