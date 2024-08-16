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

func (env *Environment) Assing(name token.Token, value loxvalue.LoxValue) error {
	if _, ok := env.values[name.Lexeme]; !ok {
		return loxerror.NewErrorFromToken(name, "Undefined variable '" + name.Lexeme + "'.")
	}
	env.values[name.Lexeme] = value
	return nil
}