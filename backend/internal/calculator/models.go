package calculator

type Operation string

const (
	Add        Operation = "add"
	Subtract   Operation = "subtract"
	Multiply   Operation = "multiply"
	Divide     Operation = "divide"
	Power      Operation = "power"
	SquareRoot Operation = "sqrt"
	Percentage Operation = "percentage"
)

type Calculation struct {
	Operation Operation
	Operands  []float64
}
