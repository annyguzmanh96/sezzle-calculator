package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"sezzle-calculator/backend/internal/calculator"
)

const maxRequestBytes = 1 << 20

type calculateFunc func(operation string, a, b float64) (float64, error)

type resultResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// NewMux returns the HTTP handler for the calculator API.
func NewMux() http.Handler {
	return newMux(calculator.Calculate)
}

func newMux(calculate calculateFunc) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calculate", func(w http.ResponseWriter, r *http.Request) {
		handleCalculate(w, r, calculate)
	})
	return mux
}

// ListenAddr returns the TCP address for port.
// An empty port selects :8080. port must be a number from 1 to 65535.
func ListenAddr(port string) (string, error) {
	if port == "" {
		return ":8080", nil
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return "", fmt.Errorf("invalid PORT %q", port)
	}
	return fmt.Sprintf(":%d", number), nil
}

func handleCalculate(w http.ResponseWriter, r *http.Request, calculate calculateFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "600")

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPost:
		postCalculate(w, r, calculate)
	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

func postCalculate(w http.ResponseWriter, r *http.Request, calculate calculateFunc) {
	defer r.Body.Close()

	body, err := readBody(w, r)
	if err != nil {
		writeBadRequest(w, err)
		return
	}

	operation, a, b, err := parseRequest(body)
	if err != nil {
		writeBadRequest(w, err)
		return
	}

	result, err := calculate(operation, a, b)
	if err != nil {
		writeCalculateError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: result})
}

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return nil, reject("request body is too large")
		}
		return nil, reject("could not read request body")
	}
	return body, nil
}

func writeCalculateError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero),
		errors.Is(err, calculator.ErrUnsupportedOperation),
		errors.Is(err, calculator.ErrNonFiniteOperand),
		errors.Is(err, calculator.ErrNonFiniteResult):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		log.Printf("calculate: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

func writeBadRequest(w http.ResponseWriter, err error) {
	var bad *badRequest
	if errors.As(err, &bad) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: bad.Error()})
		return
	}
	log.Printf("request: %v", err)
	writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("encode json: %v", err)
		status = http.StatusInternalServerError
		body = []byte(`{"error":"internal server error"}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}
