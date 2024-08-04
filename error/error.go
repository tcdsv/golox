package loxerror

import (
	"fmt"
)

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