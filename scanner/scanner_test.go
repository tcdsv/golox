package scanner_test

import (
	loxerror "golox/error"
	"golox/scanner"
	tkn "golox/token"
	loxvalue "golox/value"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokens(t *testing.T) {

	tests := []struct {
		input   string
		token 	tkn.Token
	}{
		{"!", tkn.NewToken(tkn.BANG, "!", nil, 1)},
		{"1", tkn.NewToken(tkn.NUMBER, "1", &loxvalue.Number{Value: 1}, 1)},
		{"(", tkn.NewToken(tkn.LEFT_PAREN, "(", nil, 1)},
	}

	for _, test := range tests {
		testToken(t, test.input, test.token)
	}

}

func testToken(t *testing.T, input string, expected tkn.Token) {
	
	s := scanner.NewScanner(input)
	tokens, errors := s.Scan()
	require.Empty(t, errors)
	require.Len(t, tokens, 2)
	require.Equal(t, tokens[0], expected)
	require.Equal(t, tokens[len(tokens)-1].Type, tkn.EOF)

}

func TestScan_UnterminatedString(t *testing.T) {
	s := scanner.NewScanner("\"foo")
	tokens, errors := s.Scan()
	require.Len(t, errors, 1)
	require.Len(t, tokens, 1)
	err, _ := errors[0].(*loxerror.Error)
	require.Equal(t, "unterminated string", err.Message)
	require.Equal(t, tkn.EOF, tokens[0].Type)
}

func TestScan_UnexpectedCharacter(t *testing.T) {
	s := scanner.NewScanner("&")
	tokens, errors := s.Scan()
	require.Len(t, errors, 1)
	require.Len(t, tokens, 1)
	err, _ := errors[0].(*loxerror.Error)
	require.Equal(t, "unexpected character", err.Message)
	require.Equal(t, tkn.EOF, tokens[0].Type)
}