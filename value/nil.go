package loxvalue

type Nil struct{}

func (n Nil) Type() int {
	return NIL
}