package parser_test

import (
	loxerror "golox/error"
	"golox/expr"
	"golox/parser"
	"golox/scanner"
	"golox/stmt"
	tkn "golox/token"
	loxvalue "golox/value"
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

		{"var x;", stmt.VarStmt{
            Name: tkn.NewToken(tkn.IDENTIFIER, "x", nil, 1),
            Initializer: nil,
        }},

        {"var y = 10;", stmt.VarStmt{
            Name: tkn.NewToken(tkn.IDENTIFIER, "y", nil, 1),
            Initializer: expr.LiteralExpr{Value: &loxvalue.Number{Value: 10}},
        }},
	}

	for _, test := range tests {
		testExpression(t, test.input, test.expected)
	}

}

func TestParser_GroupingExpressionError(t *testing.T) {

	tests := []struct {
		input   	string
		expected	*loxerror.Error
	}{
		{"(3;", &loxerror.Error{Line: 1, Where: " at ';'", Message: "Expect ')' after expression."}},
	}

	for _, test := range tests {
		testExpressionError(t, test.input, test.expected)
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

func testExpressionError(t *testing.T, input string, expected *loxerror.Error) {

	scanner := scanner.NewScanner(input)
	tokens, errors := scanner.Scan()
	require.Empty(t, errors)
	parser := parser.NewParser(tokens)
	_, errors = parser.Parse()
	require.NotEmpty(t, errors)
	require.Equal(t, errors[0], expected)

}