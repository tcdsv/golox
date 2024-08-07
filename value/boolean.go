package loxvalue

type Boolean struct {
	Value bool
}

func NewBoolean(value bool) *Boolean {
	return &Boolean{Value: value}
}

func (b Boolean) Type() int {
	return BOOLEAN
}