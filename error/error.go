package loxerror

import (
	"fmt"
	tkn "golox/token"
)

const PARSE_ERROR_MISSING_RIGHT_PAREN = "Expect ')' after expression."
const PARSE_ERROR_VARIABLE_EXPR_MISSING_NAME = "Expect variable name."
const PARSE_ERROR_VARIABLE_EXPR_MISSING_SEMICOLON = "Expect ';' after variable declaration."

type Error struct {
	Line    int
	Where   string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] Error %s: %s", e.Line, e.Where, e.Message)
}

func NewError(line int, where, message string) *Error {
	return &Error{
		Line:    line,
		Where:   where,
		Message: message,
	}
}

func NewErrorFromToken(token tkn.Token, message string) *Error {
	if token.Type == tkn.EOF {
		return NewError(token.Line, " at end", message)
	}
	return NewError(token.Line, " at '" + token.Lexeme + "'", message)
}
