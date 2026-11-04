package calculator

type CalculateResponse struct {
	Operation Operation `json:"operation"`
	Operands  []float64 `json:"operands"`
	Result    float64   `json:"result"`
}

func newCalculateResponse(c Calculation, result float64) CalculateResponse {
	return CalculateResponse{
		Operation: c.Operation,
		Operands:  c.Operands,
		Result:    result,
	}
}

type EvaluateResponse struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}
