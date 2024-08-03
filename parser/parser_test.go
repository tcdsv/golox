package parser_test

import (
	loxerror "golox/error"
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

func parse(source string) (expr.Expr, []loxerror.Error) {
	s := scanner.NewScanner(source)
	tokens, errors := s.Scan()
	if len(errors) > 0 {
		return nil, errors
	}
	p := parser.NewParser(tokens)
	expr, errors := p.Parse()
	if len(errors) > 0 {
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

