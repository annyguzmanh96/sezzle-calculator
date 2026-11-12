package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateSuccess(t *testing.T) {
	tests := []struct {
		name string
		body string
		want float64
		raw  string
	}{
		{name: "addition", body: `{"operation":"add","a":12,"b":5}`, want: 17, raw: `{"result":17}`},
		{name: "subtraction", body: `{"operation":"subtract","a":10,"b":4}`, want: 6, raw: `{"result":6}`},
		{name: "multiplication", body: `{"operation":"multiply","a":12,"b":5}`, want: 60, raw: `{"result":60}`},
		{name: "division", body: `{"operation":"divide","a":5,"b":2}`, want: 2.5, raw: `{"result":2.5}`},
		{name: "negative decimals", body: `{"operation":"multiply","a":-1.5,"b":0.5}`, want: -0.75, raw: `{"result":-0.75}`},
		{name: "zeros are present", body: `{"operation":"add","a":0,"b":0}`, want: 0, raw: `{"result":0}`},
		{name: "surrounding whitespace", body: " \n{\"operation\":\"add\",\"a\":1,\"b\":2}\n ", want: 3, raw: `{"result":3}`},
		{name: "exact decimals", body: `{"operation":"add","a":0.5,"b":0.25}`, want: 0.75, raw: `{"result":0.75}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := post(t, tc.body)
			payload := assertJSONResponse(t, rec, http.StatusOK)
			if _, ok := payload["error"]; ok {
				t.Fatalf("unexpected error field: %v", payload)
			}
			if got := resultOf(t, payload); got != tc.want {
				t.Fatalf("result = %v, want %v", got, tc.want)
			}
			if tc.raw != "" && strings.TrimSpace(rec.Body.String()) != tc.raw {
				t.Fatalf("body = %s, want %s", rec.Body.String(), tc.raw)
			}
			if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Fatal("missing CORS allow origin")
			}
		})
	}
}

func TestCalculateInexactDecimal(t *testing.T) {
	rec := post(t, `{"operation":"add","a":0.1,"b":0.2}`)
	payload := assertJSONResponse(t, rec, http.StatusOK)
	a := 0.1
	b := 0.2
	if got := resultOf(t, payload); got != a+b {
		t.Fatalf("result = %v, want runtime sum %v", got, a+b)
	}
}

func TestCalculateClientErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "division by zero", body: `{"operation":"divide","a":1,"b":0}`, want: "division by zero is not allowed"},
		{name: "division by negative zero", body: `{"operation":"divide","a":1,"b":-0}`, want: "division by zero is not allowed"},
		{name: "unsupported operation", body: `{"operation":"power","a":2,"b":3}`, want: "unsupported operation"},
		{name: "operation case sensitive", body: `{"operation":"Add","a":2,"b":3}`, want: "unsupported operation"},
		{name: "empty operation", body: `{"operation":"","a":2,"b":3}`, want: "operation is required"},
		{name: "missing operation", body: `{"a":1,"b":2}`, want: "operation is required"},
		{name: "null operation", body: `{"operation":null,"a":1,"b":2}`, want: "operation must be a string"},
		{name: "numeric operation", body: `{"operation":1,"a":1,"b":2}`, want: "operation must be a string"},
		{name: "missing a", body: `{"operation":"add","b":0}`, want: "a is required"},
		{name: "missing b", body: `{"operation":"subtract","a":0}`, want: "b is required"},
		{name: "null a", body: `{"operation":"add","a":null,"b":1}`, want: "a must be a number"},
		{name: "string a", body: `{"operation":"add","a":"12","b":1}`, want: "a must be a number"},
		{name: "boolean b", body: `{"operation":"add","a":1,"b":true}`, want: "b must be a number"},
		{name: "array a", body: `{"operation":"add","a":[1],"b":2}`, want: "a must be a number"},
		{name: "object b", body: `{"operation":"add","a":1,"b":{}}`, want: "b must be a number"},
		{name: "unexpected field", body: `{"operation":"add","a":1,"b":2,"c":3}`, want: "request contains unexpected fields"},
		{name: "wrong field case", body: `{"operation":"add","A":1,"b":2}`, want: "request contains unexpected fields"},
		{name: "malformed json", body: `{`, want: "request body must be valid JSON"},
		{name: "not json", body: `not-json`, want: "request body must be a JSON object"},
		{name: "trailing comma", body: `{"operation":"add","a":1,"b":2,}`, want: "request body must be valid JSON"},
		{name: "empty body", body: ``, want: "request body must be valid JSON"},
		{name: "whitespace body", body: " \n\t ", want: "request body must be valid JSON"},
		{name: "null body", body: `null`, want: "request body must be a JSON object"},
		{name: "array body", body: `[1,2]`, want: "request body must be a JSON object"},
		{name: "multiple objects", body: `{"operation":"add","a":1,"b":2}{"operation":"add","a":3,"b":4}`, want: "request body must contain a single JSON object"},
		{name: "multiple objects with newline", body: "{\"operation\":\"add\",\"a\":1,\"b\":2}\n{\"operation\":\"subtract\",\"a\":3,\"b\":4}", want: "request body must contain a single JSON object"},
		{name: "trailing junk", body: `{"operation":"add","a":1,"b":2} trailing`, want: "request body must contain a single JSON object"},
		{name: "braces inside string are one object", body: `{"operation":"}{","a":1,"b":2}`, want: "unsupported operation"},
		{name: "non finite operand", body: `{"operation":"add","a":1e309,"b":1}`, want: "a must be a finite number"},
		{name: "negative non finite operand", body: `{"operation":"add","a":1,"b":-1e309}`, want: "b must be a finite number"},
		{name: "non finite result", body: `{"operation":"multiply","a":1e308,"b":1e308}`, want: "calculation result is not a finite number"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := post(t, tc.body)
			payload := assertJSONResponse(t, rec, http.StatusBadRequest)
			if _, ok := payload["result"]; ok {
				t.Fatalf("error response includes result: %v", payload)
			}
			if got := errorOf(t, payload); got != tc.want {
				t.Fatalf("error = %q, want %q", got, tc.want)
			}
			if strings.Contains(rec.Body.String(), "Inf") || strings.Contains(rec.Body.String(), "NaN") {
				t.Fatalf("response leaked a non-finite token: %s", rec.Body.String())
			}
		})
	}
}

func TestCalculateMethods(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/calculate", strings.NewReader(`{"operation":"add","a":1,"b":2}`))
			rec := httptest.NewRecorder()
			NewMux().ServeHTTP(rec, req)

			payload := assertJSONResponse(t, rec, http.StatusMethodNotAllowed)
			if errorOf(t, payload) != "method not allowed" {
				t.Fatalf("error = %v", payload)
			}
			if rec.Header().Get("Allow") != "POST, OPTIONS" {
				t.Fatalf("Allow = %q", rec.Header().Get("Allow"))
			}
		})
	}
}

func TestCalculateOptions(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/calculate", nil)
	rec := httptest.NewRecorder()
	NewMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("body = %q", rec.Body.String())
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS origin")
	}
	if rec.Header().Get("Access-Control-Allow-Methods") != "POST, OPTIONS" {
		t.Fatalf("methods = %q", rec.Header().Get("Access-Control-Allow-Methods"))
	}
	if rec.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Fatalf("headers = %q", rec.Header().Get("Access-Control-Allow-Headers"))
	}
}

func TestCalculateRequestLimits(t *testing.T) {
	t.Run("too large", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(bytes.Repeat([]byte("a"), maxRequestBytes+1)))
		rec := httptest.NewRecorder()
		NewMux().ServeHTTP(rec, req)
		payload := assertJSONResponse(t, rec, http.StatusBadRequest)
		if errorOf(t, payload) != "request body is too large" {
			t.Fatalf("error = %v", payload)
		}
	})

	t.Run("unreadable body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/calculate", errReader{})
		rec := httptest.NewRecorder()
		NewMux().ServeHTTP(rec, req)
		payload := assertJSONResponse(t, rec, http.StatusBadRequest)
		if errorOf(t, payload) != "could not read request body" {
			t.Fatalf("error = %v", payload)
		}
	})
}

func TestCalculateInternalError(t *testing.T) {
	handler := newMux(func(string, float64, float64) (float64, error) {
		return 0, errors.New("boom")
	})
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(`{"operation":"add","a":1,"b":2}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	payload := assertJSONResponse(t, rec, http.StatusInternalServerError)
	if errorOf(t, payload) != "internal server error" {
		t.Fatalf("error = %v", payload)
	}
	if strings.Contains(rec.Body.String(), "boom") {
		t.Fatalf("response leaked internal error: %s", rec.Body.String())
	}
}

func TestWriteJSONRejectsNonFiniteNumbers(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		rec := httptest.NewRecorder()
		writeJSON(rec, http.StatusOK, resultResponse{Result: value})
		payload := assertJSONResponse(t, rec, http.StatusInternalServerError)
		if errorOf(t, payload) != "internal server error" {
			t.Fatalf("error = %v", payload)
		}
		if strings.Contains(rec.Body.String(), "NaN") || strings.Contains(rec.Body.String(), "Inf") {
			t.Fatalf("response contained a non-finite token: %s", rec.Body.String())
		}
	}
}

func TestWriteBadRequestHidesUnexpectedErrors(t *testing.T) {
	rec := httptest.NewRecorder()
	writeBadRequest(rec, errors.New("disk failed"))
	payload := assertJSONResponse(t, rec, http.StatusInternalServerError)
	if errorOf(t, payload) != "internal server error" {
		t.Fatalf("error = %v", payload)
	}
	if strings.Contains(rec.Body.String(), "disk failed") {
		t.Fatalf("response leaked internal error: %s", rec.Body.String())
	}
}

func TestParserGuards(t *testing.T) {
	_, err := parseOperation([]byte(`"\q"`))
	if err == nil || err.Error() != "operation must be a string" {
		t.Fatalf("parseOperation() error = %v", err)
	}

	_, err = parseFiniteNumber("a", []byte("1.2.3"))
	if err == nil || err.Error() != "a must be a number" {
		t.Fatalf("parseFiniteNumber() error = %v", err)
	}

	_, err = parseFiniteNumber("a", []byte("-Inf"))
	if err == nil || err.Error() != "a must be a finite number" {
		t.Fatalf("parseFiniteNumber() error = %v", err)
	}
}

func TestListenAddr(t *testing.T) {
	tests := []struct {
		port    string
		want    string
		wantErr bool
	}{
		{port: "", want: ":8080"},
		{port: "8080", want: ":8080"},
		{port: "80", want: ":80"},
		{port: "65535", want: ":65535"},
		{port: "08080", want: ":8080"},
		{port: "0", wantErr: true},
		{port: "-1", wantErr: true},
		{port: "65536", wantErr: true},
		{port: "http", wantErr: true},
		{port: ":8080", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.port, func(t *testing.T) {
			got, err := ListenAddr(tc.port)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ListenAddr(%q) = %q, want error", tc.port, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ListenAddr(%q) error = %v", tc.port, err)
			}
			if got != tc.want {
				t.Fatalf("ListenAddr(%q) = %q, want %q", tc.port, got, tc.want)
			}
		})
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (errReader) Close() error {
	return nil
}

func post(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	NewMux().ServeHTTP(rec, req)
	return rec
}

func assertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, status int) map[string]any {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, status, rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("content-type = %q", rec.Header().Get("Content-Type"))
	}

	dec := json.NewDecoder(bytes.NewReader(rec.Body.Bytes()))
	var payload map[string]any
	if err := dec.Decode(&payload); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rec.Body.String())
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		t.Fatalf("response contains extra JSON: %v", err)
	}
	if len(payload) != 1 {
		t.Fatalf("response keys = %#v", payload)
	}
	return payload
}

func resultOf(t *testing.T, payload map[string]any) float64 {
	t.Helper()
	result, ok := payload["result"].(float64)
	if !ok {
		t.Fatalf("result type = %T", payload["result"])
	}
	return result
}

func errorOf(t *testing.T, payload map[string]any) string {
	t.Helper()
	message, ok := payload["error"].(string)
	if !ok {
		t.Fatalf("error type = %T", payload["error"])
	}
	return message
}
