package calculator

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		expression string
		want       float64
	}{
		{"42", 42},
		{"1 + 2", 3},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * 4", 20},
		{"10 - 4 - 3", 3},
		{"100 / 10 / 5", 2},
		{"2 ^ 3 ^ 2", 512},
		{"-2 ^ 2", -4},
		{"(-2) ^ 2", 4},
		{"2 ^ -1", 0.5},
		{"--3", 3},
		{"+3", 3},
		{"√16", 4},
		{"√(9 + 16)", 5},
		{"√16 * 2", 8},
		{"sqrt(81)", 9},
		{"50%", 0.5},
		{"200 × 15%", 30},
		{"(1 + 1)%", 0.02},
		{"2 × −3", -6},
		{"7 ÷ 2", 3.5},
		{"((((1))))", 1},
		{".5 + 5.", 5.5},
		{"1e3 + 1E-3", 1000.001},
		{"1e+21 / 1e21", 1},
		{" 1+\t2 ", 3},
		{"-0 * 5", 0},
	}

	svc := NewService()
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			got, err := svc.Evaluate(tt.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-9 ||
				math.Signbit(got) != math.Signbit(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateSyntaxErrors(t *testing.T) {
	tests := []struct {
		expression string
		position   int
		reason     string
	}{
		{"", 1, "expression is empty"},
		{"   ", 1, "expression is empty"},
		{"1 +", 4, "unexpected end of expression"},
		{"(1 + 2", 7, `missing ")"`},
		{"(1 2)", 4, `expected ")" but found "2"`},
		{"1 + 2)", 6, `unexpected ")"`},
		{"()", 2, `unexpected ")"`},
		{"2(3)", 2, `unexpected "("`},
		{"* 2", 1, `unexpected "*"`},
		{"1.2.3", 1, `invalid number "1.2.3"`},
		{".", 1, `invalid number "."`},
		{"1e400", 1, `number "1e400" is out of range`},
		{"2 & 3", 3, `unexpected character '&'`},
		{"1e", 2, `unexpected character 'e'`},
		{"%5", 1, `unexpected "%"`},
	}

	svc := NewService()
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			_, err := svc.Evaluate(tt.expression)

			var syntaxErr *SyntaxError
			if !errors.As(err, &syntaxErr) {
				t.Fatalf("err = %v, want a SyntaxError", err)
			}
			if !errors.Is(err, ErrInvalidExpression) {
				t.Error("SyntaxError does not unwrap to ErrInvalidExpression")
			}
			if syntaxErr.Position != tt.position ||
				syntaxErr.Reason != tt.reason {
				t.Errorf("got %d %q, want %d %q", syntaxErr.Position,
					syntaxErr.Reason, tt.position, tt.reason)
			}
		})
	}
}

func TestEvaluateRejectsDeepNesting(t *testing.T) {
	deep := strings.Repeat("(", maxDepth+1) + "1" +
		strings.Repeat(")", maxDepth+1)

	var syntaxErr *SyntaxError
	_, err := NewService().Evaluate(deep)
	if !errors.As(err, &syntaxErr) ||
		syntaxErr.Reason != "expression is nested too deeply" {
		t.Errorf("err = %v, want nesting error", err)
	}
}

func TestEvaluateArithmeticErrors(t *testing.T) {
	tests := []struct {
		expression string
		want       error
	}{
		{"1 / 0", ErrDivisionByZero},
		{"5 / (3 - 3)", ErrDivisionByZero},
		{"0 ^ -1", ErrDivisionByZero},
		{"√-4", ErrUndefinedResult},
		{"(-8) ^ (1/3)", ErrUndefinedResult},
		{"10 ^ 400", ErrResultOutOfRange},
		{"1e308 * 10", ErrResultOutOfRange},
	}

	svc := NewService()
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			if _, err := svc.Evaluate(tt.expression); !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}
