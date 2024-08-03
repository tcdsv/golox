package scanner_test

import (
	loxerror "golox/error"
	"golox/scanner"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScan_LeftParen(t *testing.T) {
	s := scanner.NewScanner("(")
	tokens, errors := s.Scan()
	require.Empty(t, errors)
	require.Len(t, tokens, 2)
	leftParenToken := tokens[0]
	require.Equal(t, scanner.LEFT_PAREN, leftParenToken.Type)
	require.Equal(t, "(", leftParenToken.Lexeme)
	require.Equal(t, nil, leftParenToken.Literal)
	eofToken := tokens[1]
	require.Equal(t, scanner.EOF, eofToken.Type)
}

func TestScan_UnterminatedString(t *testing.T) {
	s := scanner.NewScanner("\"foo")
	tokens, errors := s.Scan()
	require.Len(t, errors, 1)
	require.Len(t, tokens, 1)
	err, _ := errors[0].(*loxerror.Error)
	require.Equal(t, "unterminated string", err.Message)
	require.Equal(t, scanner.EOF, tokens[0].Type)
}

func TestScan_UnexpectedCharacter(t *testing.T) {
	s := scanner.NewScanner("&")
	tokens, errors := s.Scan()
	require.Len(t, errors, 1)
	require.Len(t, tokens, 1)
	err, _ := errors[0].(*loxerror.Error)
	require.Equal(t, "unexpected character", err.Message)
	require.Equal(t, scanner.EOF, tokens[0].Type)
}