package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		operands []float64
		want     float64
	}{
		{"add", Add, []float64{2, 3}, 5},
		{"add negatives", Add, []float64{-2.5, -0.5}, -3},
		{"subtract", Subtract, []float64{2, 5}, -3},
		{"multiply", Multiply, []float64{-4, 2.5}, -10},
		{"multiply by zero is +0", Multiply, []float64{-1, 0}, 0},
		{"divide", Divide, []float64{10, 4}, 2.5},
		{"divide zero", Divide, []float64{0, 7}, 0},
		{"power", Power, []float64{2, 10}, 1024},
		{"negative exponent", Power, []float64{2, -2}, 0.25},
		{"fractional exponent", Power, []float64{9, 0.5}, 3},
		{"zero to the zero", Power, []float64{0, 0}, 1},
		{"sqrt", SquareRoot, []float64{16}, 4},
		{"sqrt zero", SquareRoot, []float64{0}, 0},
		{"percentage", Percentage, []float64{15, 200}, 30},
		{"percentage over 100", Percentage, []float64{150, 40}, 60},
	}

	svc := NewService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Calculate(Calculation{tt.op, tt.operands})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want || math.Signbit(got) != math.Signbit(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	const huge = math.MaxFloat64

	tests := []struct {
		name     string
		op       Operation
		operands []float64
		want     error
	}{
		{"unknown op", "modulo", []float64{1, 2}, ErrUnsupportedOperation},
		{"empty operation", "", []float64{1, 2}, ErrUnsupportedOperation},
		{"too few operands", Add, []float64{1}, ErrOperandCount},
		{"too many operands", Add, []float64{1, 2, 3}, ErrOperandCount},
		{"no operands", SquareRoot, nil, ErrOperandCount},
		{"binary sqrt", SquareRoot, []float64{4, 2}, ErrOperandCount},
		{"NaN operand", Add, []float64{math.NaN(), 1}, ErrInvalidOperand},
		{"infinite operand", Add, []float64{math.Inf(1), 1}, ErrInvalidOperand},
		{"divide by zero", Divide, []float64{1, 0}, ErrDivisionByZero},
		{"zero by zero", Divide, []float64{0, 0}, ErrDivisionByZero},
		{"zero to negative power", Power, []float64{0, -1}, ErrDivisionByZero},
		{"sqrt of negative", SquareRoot, []float64{-1}, ErrUndefinedResult},
		{"root of negative", Power, []float64{-8, 1.0 / 3}, ErrUndefinedResult},
		{"add overflow", Add, []float64{huge, huge}, ErrResultOutOfRange},
		{"mul overflow", Multiply, []float64{huge, 2}, ErrResultOutOfRange},
		{"pow overflow", Power, []float64{10, 400}, ErrResultOutOfRange},
		{"div overflow", Divide, []float64{huge, 0.5}, ErrResultOutOfRange},
	}

	svc := NewService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Calculate(Calculation{tt.op, tt.operands})
			if !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}
