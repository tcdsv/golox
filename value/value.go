package loxvalue

const (
	NIL = iota
	BOOLEAN
	NUMBER
	STRING
)

type LoxValue interface {
	Type() int
	ToString() string
}

func IsEqual(a LoxValue, b LoxValue) bool {
	return a == b
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