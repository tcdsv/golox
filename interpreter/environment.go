package interpreter

import (
	loxerror "golox/error"
	"golox/token"
	loxvalue "golox/value"
)

type Environment struct {
	enclosing *Environment
	values map[string]loxvalue.LoxValue
}

func NewGlobalEnv() *Environment {
	return &Environment{
		enclosing: nil,
		values: make(map[string]loxvalue.LoxValue),
	}
}

func NewLocalEnv(env *Environment) *Environment {
	return &Environment{
		enclosing: env,
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
	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}
	return nil, loxerror.NewErrorFromToken(name, "Undefined variable '" + name.Lexeme + "'.")
}

func (env *Environment) Assing(name token.Token, value loxvalue.LoxValue) error {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return nil
	}
	if env.enclosing != nil {
		return env.enclosing.Assing(name, value)
	}
	return loxerror.NewErrorFromToken(name, "Undefined variable '" + name.Lexeme + "'.")
}