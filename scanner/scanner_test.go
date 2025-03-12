package scanner_test

import (
	loxerror "golox/error"
	"golox/scanner"
	tkn "golox/token"
	loxvalue "golox/value"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanner_Tokens(t *testing.T) {

	tests := []struct {
		input   string
		token 	tkn.Token
	}{
		{"(", tkn.NewToken(tkn.LEFT_PAREN, "(", nil, 1)},
		{")", tkn.NewToken(tkn.RIGHT_PAREN, ")", nil, 1)},
		{"{", tkn.NewToken(tkn.LEFT_BRACE, "{", nil, 1)},
		{"}", tkn.NewToken(tkn.RIGHT_BRACE, "}", nil, 1)},
		{",", tkn.NewToken(tkn.COMMA, ",", nil, 1)},
		{".", tkn.NewToken(tkn.DOT, ".", nil, 1)},
		{"-", tkn.NewToken(tkn.MINUS, "-", nil, 1)},
		{"+", tkn.NewToken(tkn.PLUS, "+", nil, 1)},
		{";", tkn.NewToken(tkn.SEMICOLON, ";", nil, 1)},
		{"/", tkn.NewToken(tkn.SLASH, "/", nil, 1)},
		{"*", tkn.NewToken(tkn.STAR, "*", nil, 1)},
		{"!", tkn.NewToken(tkn.BANG, "!", nil, 1)},
		{"!=", tkn.NewToken(tkn.BANG_EQUAL, "!=", nil, 1)},
		{"=", tkn.NewToken(tkn.EQUAL, "=", nil, 1)},
		{"==", tkn.NewToken(tkn.EQUAL_EQUAL, "==", nil, 1)},
		{">", tkn.NewToken(tkn.GREATER, ">", nil, 1)},
		{">=", tkn.NewToken(tkn.GREATER_EQUAL, ">=", nil, 1)},
		{"<", tkn.NewToken(tkn.LESS, "<", nil, 1)},
		{"<=", tkn.NewToken(tkn.LESS_EQUAL, "<=", nil, 1)},
		{"and", tkn.NewToken(tkn.AND, "and", nil, 1)},
		{"class", tkn.NewToken(tkn.CLASS, "class", nil, 1)},
		{"else", tkn.NewToken(tkn.ELSE, "else", nil, 1)},
		{"false", tkn.NewToken(tkn.FALSE, "false", &loxvalue.Boolean{Value: false}, 1)},
		{"fun", tkn.NewToken(tkn.FUN, "fun", nil, 1)},
		{"for", tkn.NewToken(tkn.FOR, "for", nil, 1)},
		{"if", tkn.NewToken(tkn.IF, "if", nil, 1)},
		{"nil", tkn.NewToken(tkn.NIL, "nil", &loxvalue.Nil{}, 1)},
		{"or", tkn.NewToken(tkn.OR, "or", nil, 1)},
		{"print", tkn.NewToken(tkn.PRINT, "print", nil, 1)},
		{"return", tkn.NewToken(tkn.RETURN, "return", nil, 1)},
		{"super", tkn.NewToken(tkn.SUPER, "super", nil, 1)},
		{"this", tkn.NewToken(tkn.THIS, "this", nil, 1)},
		{"true", tkn.NewToken(tkn.TRUE, "true", &loxvalue.Boolean{Value: true}, 1)},
		{"var", tkn.NewToken(tkn.VAR, "var", nil, 1)},
		{"while", tkn.NewToken(tkn.WHILE, "while", nil, 1)},
		{"identifier", tkn.NewToken(tkn.IDENTIFIER, "identifier", nil, 1)},
		{"123", tkn.NewToken(tkn.NUMBER, "123", &loxvalue.Number{Value: 123}, 1)},
		{"\"hello\"", tkn.NewToken(tkn.STRING, "\"hello\"", &loxvalue.String{Value: "hello"}, 1)},
	}

	for _, test := range tests {
		testToken(t, test.input, test.token)
	}

}

func TestScanner_TokensError(t *testing.T) {
	
	tests := []struct {
		input   string
		expected 	*loxerror.Error
	}{
		{"&", loxerror.NewError(1, "", loxerror.SCANNER_ERROR_UNEXPECTED_CHARACTER)},
		{"\"foo", loxerror.NewError(1, "", loxerror.SCANNER_ERROR_UNTERMINATED_STRING)},
	}

	for _, test := range tests {
		testTokenError(t, test.input, test.expected)
	}

}

func testToken(t *testing.T, input string, expected tkn.Token) {
	
	s := scanner.NewScanner(input)
	tokens, errors := s.Scan()
	require.Empty(t, errors)
	require.Len(t, tokens, 2)
	require.Equal(t, expected, tokens[0])
	require.Equal(t, tkn.EOF, tokens[len(tokens)-1].Type)

}

func testTokenError(t *testing.T, input string, expected *loxerror.Error) {

	scanner := scanner.NewScanner(input)
	_, errors := scanner.Scan()
	require.NotEmpty(t, errors)
	require.Equal(t, expected, errors[0])

}