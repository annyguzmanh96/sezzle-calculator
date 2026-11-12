package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculateFloat64Addition(t *testing.T) {
	a := 0.1
	b := 0.2
	got, err := Calculate("add", a, b)
	if err != nil {
		t.Fatalf("Calculate() unexpected error = %v", err)
	}
	if got != a+b {
		t.Fatalf("Calculate() = %v, want %v", got, a+b)
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		a         float64
		b         float64
		want      float64
		wantErr   error
	}{
		{name: "addition", operation: "add", a: 12, b: 5, want: 17},
		{name: "addition of negatives", operation: "add", a: -1.5, b: -2.5, want: -4},
		{name: "addition of zeros", operation: "add", a: 0, b: 0, want: 0},
		{name: "subtraction", operation: "subtract", a: 10, b: 4, want: 6},
		{name: "subtraction to negative", operation: "subtract", a: 2, b: 5, want: -3},
		{name: "multiplication", operation: "multiply", a: 12, b: 5, want: 60},
		{name: "multiplication of negatives", operation: "multiply", a: -3, b: -4, want: 12},
		{name: "division", operation: "divide", a: 5, b: 2, want: 2.5},
		{name: "division of negatives", operation: "divide", a: -9, b: 3, want: -3},
		{name: "division by zero", operation: "divide", a: 1, b: 0, wantErr: ErrDivisionByZero},
		{name: "division by negative zero", operation: "divide", a: 1, b: math.Copysign(0, -1), wantErr: ErrDivisionByZero},
		{name: "unsupported operation", operation: "power", a: 2, b: 3, wantErr: ErrUnsupportedOperation},
		{name: "empty operation", operation: "", a: 1, b: 1, wantErr: ErrUnsupportedOperation},
		{name: "nan operand", operation: "add", a: math.NaN(), b: 1, wantErr: ErrNonFiniteOperand},
		{name: "infinite operand", operation: "subtract", a: 1, b: math.Inf(1), wantErr: ErrNonFiniteOperand},
		{name: "overflow to infinity", operation: "multiply", a: math.MaxFloat64, b: 2, wantErr: ErrNonFiniteResult},
		{name: "addition overflow", operation: "add", a: math.MaxFloat64, b: math.MaxFloat64, wantErr: ErrNonFiniteResult},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Calculate(tc.operation, tc.a, tc.b)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("Calculate() error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Calculate() unexpected error = %v", err)
			}
			if got != tc.want {
				t.Fatalf("Calculate() = %v, want %v", got, tc.want)
			}
			if !isFinite(got) {
				t.Fatalf("Calculate() returned non-finite result %v", got)
			}
		})
	}
}
