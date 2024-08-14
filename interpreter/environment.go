package interpreter

import (
	loxerror "golox/error"
	"golox/token"
	loxvalue "golox/value"
)

type Environment struct {
	values map[string]loxvalue.LoxValue
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]loxvalue.LoxValue),
	}
}

func (env *Environment) Define(name string, value loxvalue.LoxValue) {
	env.values[name] = value
}

func (env *Environment) Get(name token.Token) (loxvalue.LoxValue, error) {
	if value, ok := env.values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, loxerror.NewErrorFromToken(name, "Undefined variable '" + name.Lexeme + "'.")
}