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

func TestParser_LiteralExpressions(t *testing.T) {
	
	tests := []struct {
		input   	string
		expected	stmt.Stmt
	}{
		{"\"foo\";", stmt.ExprStmt{E: expr.LiteralExpr{Value: loxvalue.NewString("foo")}}},
		
		{"123;", stmt.ExprStmt{
            E: expr.LiteralExpr{Value: &loxvalue.Number{Value: 123}},
        }},

		{"true;", stmt.ExprStmt{
            E: expr.LiteralExpr{Value: loxvalue.NewBoolean(true)},
        }},

		{"nil;", stmt.ExprStmt{
            E: expr.LiteralExpr{Value: &loxvalue.Nil{}},
        }},
	}

	for _, test := range tests {
		testExpression(t, test.input, test.expected)
	}

}

func TestParser_UnaryExpressions(t *testing.T) {

	tests := []struct {
		input   	string
		expected	stmt.Stmt
	}{
		{"!true;", stmt.ExprStmt{
			E: expr.UnaryExpr{
				Operator: tkn.NewToken(tkn.BANG, "!", nil, 1),
				Right: expr.LiteralExpr{Value: loxvalue.NewBoolean(true)},
			},
		}},
	}

	for _, test := range tests {
		testExpression(t, test.input, test.expected)
	}

}

func TestParser_BinaryExpressions(t *testing.T) {

	tests := []struct {
		input   	string
		expected	stmt.Stmt
	}{
        {"1 + 2;", stmt.ExprStmt{
            E: expr.BinaryExpr{
                Left:     expr.LiteralExpr{Value: &loxvalue.Number{Value: 1}},
                Operator: tkn.NewToken(tkn.PLUS, "+", nil, 1),
                Right:    expr.LiteralExpr{Value: &loxvalue.Number{Value: 2}},
            },
        }},
			
	}

	for _, test := range tests {
		testExpression(t, test.input, test.expected)
	}

}

func TestParser_LogicalExpressions(t *testing.T) {

	tests := []struct {
		input   	string
		expected	stmt.Stmt
	}{
        {"true or false;", stmt.ExprStmt{
            E: expr.LogicalExpr{
                Left:     expr.LiteralExpr{Value: loxvalue.NewBoolean(true)},
                Operator: tkn.NewToken(tkn.OR, "or", nil, 1),
                Right:    expr.LiteralExpr{Value: loxvalue.NewBoolean(false)},
            },
        }},
	}

	for _, test := range tests {
		testExpression(t, test.input, test.expected)
	}

}

func TestParser_GroupingExpressions(t *testing.T) {

	tests := []struct {
		input   	string
		expected	stmt.Stmt
	}{
		{"(3);", stmt.ExprStmt{
            E: expr.GroupingExpr{
                Expr: expr.LiteralExpr{Value: &loxvalue.Number{Value: 3}},
            },
        }},
	}

	for _, test := range tests {
		testExpression(t, test.input, test.expected)
	}

}

func TestParser_VariableExpressions(t *testing.T) {
	
	tests := []struct {
		input   	string
		expected	stmt.Stmt
	}{
		{"a;", stmt.ExprStmt{
            E: expr.VariableExpr{
                Name: tkn.NewToken(tkn.IDENTIFIER, "a", nil, 1),
            },
        }},

        {"a = 42;", stmt.ExprStmt{
            E: expr.AssignExpr{
                Name:  tkn.NewToken(tkn.IDENTIFIER, "a", nil, 1),
                Right: expr.LiteralExpr{Value:  &loxvalue.Number{Value: 42}},
            },
        }},
	}

	for _, test := range tests {
		testExpression(t, test.input, test.expected)
	}

}

func testExpression(t *testing.T, input string, expected stmt.Stmt) {

	scanner := scanner.NewScanner(input)
	tokens, errors := scanner.Scan()
	require.Empty(t, errors)
	parser := parser.NewParser(tokens)
	statements, errors := parser.Parse()
	require.Empty(t, errors)
	require.Equal(t, statements[0], expected)

}

func TestParser_VariableDeclarationWithExpression(t *testing.T) {
	file, err := loadFile("variable_declaration_2.lox")
	require.NoError(t, err)
	statements, errors := parse(file)
	require.Empty(t, errors)
	varStmt, ok := statements[0].(stmt.VarStmt)
	require.True(t, ok)
	require.Equal(t, tkn.IDENTIFIER, varStmt.Name.Type)
	require.NotNil(t, varStmt.Initializer)
	e, ok := varStmt.Initializer.(expr.LiteralExpr)
	require.True(t, ok)
	require.Equal(t, loxvalue.STRING, e.Value.Type())
	require.Equal(t, "foo", e.Value.ToString())
}

func TestParser_VariableDeclaration(t *testing.T) {
	file, err := loadFile("variable_declaration_1.lox")
	require.NoError(t, err)
	statements, errors := parse(file)
	require.Empty(t, errors)
	varStmt, ok := statements[0].(stmt.VarStmt)
	require.True(t, ok)
	require.Equal(t, tkn.IDENTIFIER, varStmt.Name.Type)
	require.Nil(t, varStmt.Initializer)
}

func TestParser_GroupingMissingParen(t *testing.T) {
	file, err := loadFile("grouping_error.lox")
	require.NoError(t, err)
	_, errors := parse(file)
	require.Len(t, errors, 1)
}

func TestParser_MissingExpression(t *testing.T) {
	t.Skip()
	//todo: refactor test
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
	return p.Parse()
	/*if err != nil {
		errors := []error{err}
		return nil, errors
	}
	return statements, nil*/
}

func loadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

