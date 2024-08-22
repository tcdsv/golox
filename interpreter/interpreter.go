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
	env     *Environment
}


func (i *Interpreter) VisitLogical(LogicalExpr expr.LogicalExpr) expr.LoxValueResult {

	left := i.Evaluate(LogicalExpr.Left) 
	if left.Error != nil {
		return left
	}

	operator := LogicalExpr.Operator.Type
	switch (operator) {
		case tkn.OR:
	
			if loxvalue.IsTruthy(left.Value) {
				return left
			}
	
		case tkn.AND:
		
			if !loxvalue.IsTruthy(left.Value) {
				return left
			}
	}
	return i.Evaluate(LogicalExpr.Right)

}

// VisitIfStatement implements stmt.StmtVisitor.
func (i *Interpreter) VisitIfStatement(IfStmt stmt.IfStmt) *visitor.VisitorResult {
	
	conditionValue := i.Evaluate(IfStmt.Condition)
	if conditionValue.Error != nil {
		return visitor.NewVisitorResult(nil, conditionValue.Error)
	}

	isConditionTrue := loxvalue.IsTruthy(conditionValue.Value) 
	if isConditionTrue {
		return i.execute(IfStmt.ThenBrnach)
	} else if IfStmt.ElseBranch != nil {
		return i.execute(IfStmt.ElseBranch)
	}

	return visitor.NewVisitorResult(nil, nil)

}

// VisitBlockStatement implements stmt.StmtVisitor.
func (i *Interpreter) VisitBlockStatement(BlockStmt stmt.BlockStmt) *visitor.VisitorResult {
	i.executeBlock(BlockStmt.Statements, NewLocalEnv(i.env))
	return visitor.NewVisitorResult(nil, nil)
}

func (i *Interpreter) executeBlock(statements []stmt.Stmt, localEnv *Environment) *visitor.VisitorResult {

	prevEnv := localEnv.enclosing
	i.env = localEnv

	defer func() {
		i.env = prevEnv
	}()

	for _, statement := range statements {
		res := i.execute(statement)
		if res.Err != nil {
			return visitor.NewVisitorResult(nil, res.Err)
		}
	}

	return visitor.NewVisitorResult(nil, nil)

}

// VisitAssing implements expr.ExprVisitor.
func (i *Interpreter) VisitAssing(assignExpr expression.AssignExpr) expr.LoxValueResult {

	res := i.Evaluate(assignExpr.Right)
	if res.Error != nil {
		return res
	}

	//todo: design a solution to get rid of this casting thing.
	// if res.Err is nil, then casting should be done at an earlier stage.
	err := i.env.Assing(assignExpr.Name, res.Value)
	if err != nil {
		return expression.LoxValueResult{Error: err}
	}
	return res
}

// VisitVariableStatement implements stmt.StmtVisitor.
func (i *Interpreter) VisitVariableStatement(variableStmt stmt.VarStmt) *visitor.VisitorResult {

	var varExpr loxvalue.LoxValue
	if variableStmt.Initializer != nil {
		initializer := i.Evaluate(variableStmt.Initializer)
		if initializer.Error != nil {
			return visitor.NewVisitorResult(nil, initializer.Error)
		}
		varExpr = initializer.Value
	} else {
		varExpr = loxvalue.Nil{}
	}
	i.env.Define(variableStmt.Name.Lexeme, varExpr)
	return visitor.NewVisitorResult(nil, nil)

}

// VisitVariable implements expr.ExprVisitor.
func (i *Interpreter) VisitVariable(element expression.VariableExpr) expr.LoxValueResult {
	variable, err := i.env.Get(element.Name)
	return expression.LoxValueResult{Value: variable, Error: err}
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Results: []*visitor.VisitorResult{},
		env:     NewGlobalEnv(),
	}
}

func (i *Interpreter) Evaluate(expr expression.Expr) expr.LoxValueResult {
	return expr.Evaluate(i)
}

func (i *Interpreter) execute(stmt stmt.Stmt) *visitor.VisitorResult {
	return stmt.Accept(i)
}

func (i *Interpreter) Interpret(statements []stmt.Stmt) {
	for _, statement := range statements {
		result := i.execute(statement)
		i.Results = append(i.Results, result)
	}
}

func (i *Interpreter) VisitPrintStatement(printStmt stmt.PrintStmt) *visitor.VisitorResult {

	lv := i.Evaluate(printStmt.E) //todo: handle evaluation errors.
	fmt.Println(lv.Value.ToString())
	return visitor.NewVisitorResult(nil, nil)

}

func (i *Interpreter) VisitExpressionStatement(exprStmt stmt.ExprStmt) *visitor.VisitorResult {
	res := i.Evaluate(exprStmt.E)
	return visitor.NewVisitorResult(res, res.Error)
}

func (i *Interpreter) VisitLiteral(expr expression.LiteralExpr) expr.LoxValueResult {
	return expression.LoxValueResult{Value: expr.Value}
}

func (i *Interpreter) VisitUnary(expr expr.UnaryExpr) expr.LoxValueResult {

	right := i.Evaluate(expr.Right)
	if right.Error != nil {
		return right
	}

	var result loxvalue.LoxValue

	switch expr.Operator.Type {
	case tkn.MINUS:
		number, err := checkNumberOperand(expr.Operator, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		result = number.Minus()
	case tkn.BANG:
		result = loxvalue.NewBoolean(!loxvalue.IsTruthy(right.Value))
	}
	//todo: error
	return expression.LoxValueResult{Value: result}

}

func (i *Interpreter) VisitBinary(expr expr.BinaryExpr) expr.LoxValueResult {

	left := i.Evaluate(expr.Left)
	right := i.Evaluate(expr.Right)

	if left.Error != nil {
		return left
	}
	if right.Error != nil {
		return right
	}

	switch expr.Operator.Type {
	case tkn.MINUS:

		left, right, err := checkNumberOperands(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: left.Subtract(right)}

	case tkn.SLASH:

		left, right, err := checkNumberOperands(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: left.Divide(right)}

	case tkn.STAR:

		left, right, err := checkNumberOperands(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: left.Multiply(right)}

	case tkn.PLUS:

		result, err := binaryPlus(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: result}

	case tkn.GREATER:

		left, right, err := checkNumberOperands(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: left.Greater(right)}

	case tkn.GREATER_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: left.GreaterEqual(right)}

	case tkn.LESS:

		left, right, err := checkNumberOperands(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: left.Less(right)}

	case tkn.LESS_EQUAL:

		left, right, err := checkNumberOperands(expr.Operator, left.Value, right.Value)
		if err != nil {
			return expression.LoxValueResult{Error: err}
		}
		return expression.LoxValueResult{Value: left.LessEqual(right)}

	case tkn.BANG_EQUAL:

		result := loxvalue.NewBoolean(!loxvalue.IsEqual(left.Value, right.Value))
		return expression.LoxValueResult{Value: result}

	case tkn.EQUAL_EQUAL:

		result := loxvalue.NewBoolean(loxvalue.IsEqual(left.Value, right.Value))
		return expression.LoxValueResult{Value: result}

	}

	//todo: error
	return expression.LoxValueResult{Error: nil}
}

func (i *Interpreter) VisitGrouping(expr expr.GroupingExpr) expr.LoxValueResult {
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
