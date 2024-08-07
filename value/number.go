package loxvalue

type Number struct {
	Value float64
}

func (n Number) Type() int {
	return NUMBER
}

func (n Number) Minus(v Number) Number {
	return Number{
		Value: n.Value - v.Value,
	}
}

func (n Number) Plus(v Number) Number {
	return Number{
		Value: n.Value + v.Value,
	}
}

func (n Number) Divide(v Number) Number {
	return Number{
		Value: n.Value / v.Value,
	}
}

func (n Number) Multiply(v Number) Number {
	return Number{
		Value: n.Value * v.Value,
	}
}

func (n Number) Greater(v Number) Boolean {
	return Boolean{
		Value: n.Value > v.Value,
	}
}

func (n Number) GreaterEqual(v Number) Boolean {
	return Boolean{
		Value: n.Value >= v.Value,
	}
}

func (n Number) Less(v Number) Boolean {
	return Boolean{
		Value: n.Value < v.Value,
	}
}

func (n Number) LessEqual(v Number) Boolean {
	return Boolean{
		Value: n.Value <= v.Value,
	}
}