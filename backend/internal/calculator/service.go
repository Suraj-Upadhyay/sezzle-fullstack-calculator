package calculator

import (
	"fmt"
	"math"
)

type operation struct {
	arity int
	apply func(x []float64) (float64, error)
}

var operations = map[Operation]operation{
	Add:        {2, add},
	Subtract:   {2, subtract},
	Multiply:   {2, multiply},
	Divide:     {2, divide},
	Power:      {2, power},
	SquareRoot: {1, squareRoot},
	Percentage: {2, percentage},
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Calculate(c Calculation) (float64, error) {
	op, ok := operations[c.Operation]
	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrUnsupportedOperation, c.Operation)
	}
	if len(c.Operands) != op.arity {
		return 0, fmt.Errorf("%w: %s takes %d, got %d",
			ErrOperandCount, c.Operation, op.arity, len(c.Operands))
	}
	for _, x := range c.Operands {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return 0, fmt.Errorf("%w: %v", ErrInvalidOperand, x)
		}
	}

	result, err := op.apply(c.Operands)
	switch {
	case err != nil:
		return 0, err
	case math.IsNaN(result):
		return 0, ErrUndefinedResult
	case math.IsInf(result, 0):
		return 0, ErrResultOutOfRange
	case result == 0:
		return 0, nil // normalises -0
	default:
		return result, nil
	}
}

func add(x []float64) (float64, error) {
	return x[0] + x[1], nil
}

func subtract(x []float64) (float64, error) {
	return x[0] - x[1], nil
}

func multiply(x []float64) (float64, error) {
	return x[0] * x[1], nil
}

func divide(x []float64) (float64, error) {
	if x[1] == 0 {
		return 0, ErrDivisionByZero
	}
	return x[0] / x[1], nil
}

func power(x []float64) (float64, error) {
	if x[0] == 0 && x[1] < 0 {
		return 0, ErrDivisionByZero
	}
	return math.Pow(x[0], x[1]), nil
}

func squareRoot(x []float64) (float64, error) {
	if x[0] < 0 {
		return 0, fmt.Errorf("%w: square root of %v", ErrUndefinedResult, x[0])
	}
	return math.Sqrt(x[0]), nil
}

// percentage returns x[0] percent of x[1].
func percentage(x []float64) (float64, error) {
	return x[0] * x[1] / 100, nil
}
