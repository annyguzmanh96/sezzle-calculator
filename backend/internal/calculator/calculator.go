// Package calculator performs the arithmetic for the calculator API.
package calculator

import (
	"errors"
	"math"
)

var (
	// ErrDivisionByZero is returned when the divisor is zero.
	ErrDivisionByZero = errors.New("division by zero is not allowed")

	// ErrUnsupportedOperation is returned when operation is not recognized.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrNonFiniteOperand is returned when an operand is NaN or infinite.
	ErrNonFiniteOperand = errors.New("operands must be finite numbers")

	// ErrNonFiniteResult is returned when the arithmetic result is not finite.
	ErrNonFiniteResult = errors.New("calculation result is not a finite number")
)

// Calculate applies operation to a and b.
// Supported operations are add, subtract, multiply, and divide.
func Calculate(operation string, a, b float64) (float64, error) {
	if !isFinite(a) || !isFinite(b) {
		return 0, ErrNonFiniteOperand
	}

	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		result = a / b
	default:
		return 0, ErrUnsupportedOperation
	}

	if !isFinite(result) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
