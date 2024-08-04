package parser_test

import (
	"golox/expr"
	"golox/parser"
	"golox/scanner"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_String(t *testing.T) {
	file, err := loadFile("string.lox")
	require.NoError(t, err)
	e, errors := parse(file)
	require.Empty(t, errors)
	literalExpr, ok := e.(expr.LiteralExpr)
	require.True(t, ok)
	literalValue, ok := literalExpr.Value.(string)
	require.True(t, ok)
	require.Equal(t, "foo", literalValue)
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

