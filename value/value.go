package loxvalue

const (
	NIL = iota
	BOOLEAN
	NUMBER
	STRING
)

type LoxValue interface {
	Type() int
	toString() string
}

func IsEqual(a LoxValue, b LoxValue) bool {
	if a.Type() != b.Type() {
		return false
	}
	if a.Type() == NIL {
		return true
	}

	var firstValue interface{}
	var secondValue interface{}

	switch a.Type() {
	case BOOLEAN:
		firstValue = a.(Boolean).Value
		secondValue = b.(Boolean).Value
	case NUMBER:
		firstValue = a.(Number).Value
		secondValue = b.(Number).Value
	case STRING:
		firstValue = a.(String).Value
		secondValue = b.(String).Value
	}

	return firstValue == secondValue
}

func IsTruthy(value LoxValue) bool {
	if value.Type() == NIL {
		return false
	}
	if value.Type() == BOOLEAN {
		lB := value.(*Boolean)
		return lB.Value
	}
	return true
}