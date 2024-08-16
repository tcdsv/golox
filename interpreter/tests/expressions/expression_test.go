package expression_test

import (
	"golox/interpreter"
	"golox/parser"
	"golox/scanner"
	"golox/stmt"
	loxvalue "golox/value"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpreter_Assign(t *testing.T) {
	file, err := loadFile("assign.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	require.NotEmpty(t, e)
	_, err = interpret(e)
	require.NoError(t, err)
}

func TestInterpreter_Print(t *testing.T) {
	file, err := loadFile("print.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	require.NotEmpty(t, e)
	_, err = interpret(e)
	require.NoError(t, err)
}

func TestInterpreter_BinaryExprPlusError(t *testing.T) {
	file, err := loadFile("binary_expr_plus_error.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value, err := interpret(e)
	require.Nil(t, value)
	require.Error(t, err)
}

func TestInterpreter_UnaryExprMinusError(t *testing.T) {
	file, err := loadFile("unary_expr_minus_error.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value, err := interpret(e)
	require.Nil(t, value)
	require.Error(t, err)
}

func TestInterpreter_GroupingExpr(t *testing.T) {
	file, err := loadFile("grouping_expr.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value, err := interpret(e)

	require.NoError(t, err)
	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(5), loxNumber.Value)
}

func TestInterpreter_BinaryExprPlusNumbers(t *testing.T) {
	file, err := loadFile("binary_expr_plus_numbers.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value, err := interpret(e)

	require.NoError(t, err)
	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(10), loxNumber.Value)
}

func TestInterpreter_BinaryExprMinus(t *testing.T) {
	file, err := loadFile("binary_expr_minus.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value, err := interpret(e)

	require.NoError(t, err)
	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(2), loxNumber.Value)
}

func TestInterpreter_UnaryExprMinus(t *testing.T) {
	file, err := loadFile("unary_expr_minus.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value, err := interpret(e)

	require.NoError(t, err)
	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(-1), loxNumber.Value)
}

func TestInterpreter_UnaryExprBang(t *testing.T) {

	file, err := loadFile("unary_expr_bang.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value, err := interpret(e)

	require.NoError(t, err)
	loxBoolean, ok := value.(*loxvalue.Boolean)
	require.True(t, ok)
	require.Equal(t, false, loxBoolean.Value)
}

func interpret(statements []stmt.Stmt) (loxvalue.LoxValue, error) {
	i := interpreter.NewInterpreter()
	i.Interpret(statements)
	if i.Results[0].Err != nil {
		return nil, i.Results[0].Err
	}
	loxValue, ok := i.Results[0].Result.(loxvalue.LoxValue)
	if ok {
		return loxValue, nil
	}
	return nil, nil
}

func parse(source string) ([]stmt.Stmt, []error) {
	s := scanner.NewScanner(source)
	tokens, errors := s.Scan()
	if len(errors) > 0 {
		return nil, errors
	}
	p := parser.NewParser(tokens)
	return p.Parse()
}

func loadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}