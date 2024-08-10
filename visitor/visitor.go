package visitor

type VisitorResult struct {
	Result interface{}
	Err    error
}

func NewVisitorResult(result interface{}, err error) *VisitorResult {
	return &VisitorResult{
		Result: result,
		Err:    err,
	}
}