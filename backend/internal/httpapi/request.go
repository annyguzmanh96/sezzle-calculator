package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"strconv"
)

type badRequest struct {
	message string
}

func (e *badRequest) Error() string {
	return e.message
}

func reject(message string) error {
	return &badRequest{message: message}
}

func parseRequest(body []byte) (string, float64, float64, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return "", 0, 0, reject("request body must be valid JSON")
	}
	if trimmed[0] != '{' {
		return "", 0, 0, reject("request body must be a JSON object")
	}

	dec := json.NewDecoder(bytes.NewReader(trimmed))
	var raw map[string]json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return "", 0, 0, reject("request body must be valid JSON")
	}

	var extra json.RawMessage
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return "", 0, 0, reject("request body must contain a single JSON object")
	}

	for key := range raw {
		switch key {
		case "operation", "a", "b":
		default:
			return "", 0, 0, reject("request contains unexpected fields")
		}
	}

	operation, err := parseOperation(raw["operation"])
	if err != nil {
		return "", 0, 0, err
	}
	a, err := parseFiniteNumber("a", raw["a"])
	if err != nil {
		return "", 0, 0, err
	}
	b, err := parseFiniteNumber("b", raw["b"])
	if err != nil {
		return "", 0, 0, err
	}
	return operation, a, b, nil
}

func parseOperation(raw json.RawMessage) (string, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return "", reject("operation is required")
	}
	if trimmed[0] != '"' {
		return "", reject("operation must be a string")
	}

	var operation string
	if err := json.Unmarshal(trimmed, &operation); err != nil {
		return "", reject("operation must be a string")
	}
	if operation == "" {
		return "", reject("operation is required")
	}
	return operation, nil
}

func parseFiniteNumber(field string, raw json.RawMessage) (float64, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return 0, reject(field + " is required")
	}
	if !isJSONNumber(trimmed) {
		return 0, reject(field + " must be a number")
	}

	value, err := strconv.ParseFloat(string(trimmed), 64)
	if err != nil {
		var numErr *strconv.NumError
		if errors.As(err, &numErr) && errors.Is(numErr.Err, strconv.ErrRange) {
			return 0, reject(field + " must be a finite number")
		}
		return 0, reject(field + " must be a number")
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, reject(field + " must be a finite number")
	}
	return value, nil
}

func isJSONNumber(raw []byte) bool {
	return raw[0] == '-' || (raw[0] >= '0' && raw[0] <= '9')
}
