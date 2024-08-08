package expression_test

import (
	"golox/expr"
	"golox/interpreter"
	"golox/parser"
	"golox/scanner"
	loxvalue "golox/value"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpreter_Grouping(t *testing.T) {
	file, err := loadFile("grouping.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value := interpret(e)

	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(5), loxNumber.Value)
}

func TestInterpreter_Addition(t *testing.T) {
	file, err := loadFile("addition.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value := interpret(e)

	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(10), loxNumber.Value)
}

func TestInterpreter_Subtract(t *testing.T) {
	file, err := loadFile("subtract.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value := interpret(e)

	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(2), loxNumber.Value)
}

func TestInterpreter_Minus(t *testing.T) {
	file, err := loadFile("minus.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value := interpret(e)

	loxNumber, ok := value.(*loxvalue.Number)
	require.True(t, ok)
	require.Equal(t, float64(-1), loxNumber.Value)
}

func TestInterpreter_Bang(t *testing.T) {

	file, err := loadFile("bang.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	value := interpret(e)

	loxBoolean, ok := value.(*loxvalue.Boolean)
	require.True(t, ok)
	require.Equal(t, false, loxBoolean.Value)
}

func interpret(e expr.Expr) loxvalue.LoxValue {
	i := interpreter.NewInterpreter()
	i.Interpret(e)
	return i.Value
}

func parse(source string) (expr.Expr, []error) {
	s := scanner.NewScanner(source)
	tokens, errors := s.Scan()
	if len(errors) > 0 {
		return nil, errors
	}
	p := parser.NewParser(tokens)
	expr, err := p.Parse()
	if err != nil {
		errors := []error{err}
		return nil, errors
	}
	return expr, nil
}

func loadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}