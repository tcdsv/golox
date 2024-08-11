package loxvalue

type String struct {
	Value string
}

func NewString(text string) *String {
	return &String{
		Value: text,
	}
}

func (s String) Type() int {
	return STRING
}

func (s String) toString() string {
	return s.Value
}

func (s String) Concat(v *String) *String {
	return &String{Value: s.Value + v.Value}
}