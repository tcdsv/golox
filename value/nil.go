package loxvalue

type Nil struct{}

func (n Nil) Type() int {
	return NIL
}

func (n Nil) toString() string {
	return "nil"
}