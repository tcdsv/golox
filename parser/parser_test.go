package parser_test

import (
	"golox/expr"
	"golox/parser"
	"golox/scanner"
	"golox/stmt"
	tkn "golox/token"
	loxvalue "golox/value"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_Unary(t *testing.T) {
	file, err := loadFile("unary.lox")
	require.NoError(t, err)
	statements, errors := parse(file)
	require.Empty(t, errors)
	unaryExpr, ok := statements[0].(stmt.ExprStmt).E.(expr.UnaryExpr)
	require.True(t, ok)
	require.Equal(t, tkn.BANG, unaryExpr.Operator.Type)
	literalExpr, ok := unaryExpr.Right.(expr.LiteralExpr)
	require.True(t, ok)
	literalBool, ok := literalExpr.Value.(*loxvalue.Boolean)
	require.True(t, ok)
	require.Equal(t, true, literalBool.Value)
}

func TestParser_String(t *testing.T) {
	file, err := loadFile("string.lox")
	require.NoError(t, err)
	statements, errors := parse(file)
	require.Empty(t, errors)
	literalExpr, ok := statements[0].(stmt.ExprStmt).E.(expr.LiteralExpr)
	require.True(t, ok)
	literalValue, ok := literalExpr.Value.(*loxvalue.String)
	require.True(t, ok)
	require.Equal(t, "foo", literalValue.Value)
}

func TestParser_GroupingMissingParen(t *testing.T) {
	file, err := loadFile("grouping_error.lox")
	require.NoError(t, err)
	_, errors := parse(file)
	require.Len(t, errors, 1)
}

func TestParser_MissingExpression(t *testing.T) {
	file, err := loadFile("missing_expression_error.lox")
	require.NoError(t, err)
	_, errors := parse(file)
	require.Len(t, errors, 1)
}

func parse(source string) ([]stmt.Stmt, []error) {
	s := scanner.NewScanner(source)
	tokens, errors := s.Scan()
	if len(errors) > 0 {
		return nil, errors
	}
	p := parser.NewParser(tokens)
	statements, err := p.Parse()
	if err != nil {
		errors := []error{err}
		return nil, errors
	}
	return statements, nil
}

func loadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

