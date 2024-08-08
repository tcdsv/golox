package loxvalue

import "strconv"

type Number struct {
	Value float64
}

func NewNumberFromText(number string) (*Number, error) {
	value, err := strconv.ParseFloat(number, 64)
	return &Number{Value: value}, err
}

func (n Number) Type() int {
	return NUMBER
}

func (n Number) Minus() *Number {
	return &Number{Value: -n.Value,}
}

func (n Number) Subtract(v *Number) *Number {
	return &Number{Value: n.Value - v.Value,}
}

func (n Number) Add(v *Number) *Number {
	return &Number{Value: n.Value + v.Value,}
}

func (n Number) Divide(v *Number) *Number {
	return &Number{Value: n.Value / v.Value,}
}

func (n Number) Multiply(v *Number) *Number {
	return &Number{Value: n.Value * v.Value,}
}

func (n Number) Greater(v *Number) *Boolean {
	return &Boolean{Value: n.Value > v.Value,}
}

func (n Number) GreaterEqual(v *Number) *Boolean {
	return &Boolean{Value: n.Value >= v.Value,}
}

func (n Number) Less(v *Number) *Boolean {
	return &Boolean{Value: n.Value < v.Value,}
}

func (n Number) LessEqual(v *Number) *Boolean {
	return &Boolean{Value: n.Value <= v.Value,}
}