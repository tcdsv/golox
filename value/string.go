package loxvalue

type String struct {
	Value string
}

func (s String) Type() int {
	return STRING
}

func (s String) Concat(v String) String {
	return String{
		Value: s.Value + v.Value,
	}
}