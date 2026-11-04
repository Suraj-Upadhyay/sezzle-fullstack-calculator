package calculator

type CalculateRequest struct {
	Operation string     `json:"operation" validate:"required"`
	Operands  []*float64 `json:"operands" validate:"required,dive,required"`
}

func (req *CalculateRequest) toDomain() Calculation {
	operands := make([]float64, len(req.Operands))
	for i, x := range req.Operands {
		operands[i] = *x
	}
	return Calculation{
		Operation: Operation(req.Operation),
		Operands:  operands,
	}
}

type EvaluateRequest struct {
	Expression string `json:"expression" validate:"required,max=500"`
}
