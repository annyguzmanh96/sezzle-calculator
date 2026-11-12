# Sezzle Calculator

A small full-stack calculator. The React client collects two numbers and an operation. The Go API performs the arithmetic and returns JSON.

Optional operations (exponentiation, square root, and percentage) are not implemented. The four required operations, their tests, and this documentation were the scope of the project.

## Technology stack

| Layer | Choice |
| --- | --- |
| Frontend | React 19, TypeScript 5.9, Vite 6.4 |
| Styling | Global CSS and one CSS module |
| Frontend tests | Vitest 3.2, React Testing Library, jsdom |
| Backend | Go 1.22+ standard library (`net/http`, `encoding/json`) |
| Backend tests | `go test` |
| Data | None. There is no database, account system, or stored history. |

Vite 6 is pinned because Vite 8 requires Node.js 20.19 or newer. This workspace was verified with Node.js 20.18.0. `@testing-library/jest-dom` is pinned to 6.9.1 for the same reason: 6.10 and 7 require Node.js 22.

## Architecture

```text
sezzle-calculator/
├── backend/
│   ├── main.go                         Process entrypoint
│   ├── go.mod
│   └── internal/
│       ├── calculator/                 Arithmetic only
│       └── httpapi/                    JSON parsing, HTTP, CORS
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── src/
│       ├── App.tsx                     Page shell
│       ├── components/Calculator.tsx   Form and result display
│       ├── services/calculate.ts       API client
│       └── types.ts                    Request and response types
├── README.md
├── AI_USAGE.md
├── PROMPT.md
└── .gitignore
```

`calculator` does not import `net/http`. `httpapi` parses the request, calls `calculator.Calculate`, and maps errors to status codes. The browser calls `POST /api/calculate`. During `npm run dev` and `npm run preview`, Vite proxies `/api` to the Go server, so the page stays on one origin. The API also sends CORS headers for direct calls.

The client does not add, subtract, multiply, or divide. It displays the `result` field from the response.

## Prerequisites

- Go 1.22 or newer. Verified with Go 1.27.4.
- Node.js 20.18 or newer, and npm. Verified with Node.js 20.18.0 and npm 10.8.2.

No database, Docker, or environment file is required.

## Installation

From the repository root, in PowerShell:

```powershell
Set-Location frontend
npm install
```

The Go module has no third-party dependencies, so the backend does not need a separate install.

## Run the backend

From `backend`:

```powershell
Set-Location backend
go run .
```

The process listens on `:8080`. Set `PORT` to a number from 1 through 65535 to change it, for example `$env:PORT = "9090"`.

## Run the frontend

From `frontend`, in a second terminal, after the API is running:

```powershell
Set-Location frontend
npm run dev
```

Open http://localhost:5173/.

`npm run preview` serves the production build on http://localhost:4173/ and uses the same `/api` proxy.

## Default ports

| Process | URL |
| --- | --- |
| Go API | http://127.0.0.1:8080 |
| Vite dev server | http://localhost:5173 |
| Vite preview | http://localhost:4173 |

## REST API

`POST /api/calculate`

`Content-Type: application/json` is what clients should send. The handler reads the body and does not require the header.

Request fields:

| Field | Type | Rule |
| --- | --- | --- |
| `operation` | string | Required. Exactly `add`, `subtract`, `multiply`, or `divide`. |
| `a` | number | Required. Zero is valid. Missing, null, strings, booleans, arrays, and objects are rejected. |
| `b` | number | Same rules as `a`. |

Field names are case-sensitive. Unknown fields are rejected. The body must be one JSON object. Arithmetic uses IEEE-754 `float64`.

### Success

```powershell
Invoke-RestMethod -Method Post -Uri http://127.0.0.1:8080/api/calculate -ContentType "application/json" -Body '{"operation":"multiply","a":12,"b":5}'
```

```json
{"result":60}
```

`{"operation":"add","a":0,"b":0}` is valid and returns `{"result":0}`.

### Errors

Every error body has one string field:

```json
{"error":"division by zero is not allowed"}
```

| Status | When | Example `error` |
| --- | --- | --- |
| 400 | Missing `a` | `a is required` |
| 400 | `a` is null, a string, or another non-number | `a must be a number` |
| 400 | A numeric operand overflows `float64` | `a must be a finite number` |
| 400 | Empty `operation` | `operation is required` |
| 400 | Unknown operation, including `Add` | `unsupported operation` |
| 400 | Divisor is zero | `division by zero is not allowed` |
| 400 | The result is infinite | `calculation result is not a finite number` |
| 400 | Extra field, including `A` instead of `a` | `request contains unexpected fields` |
| 400 | Empty or malformed JSON | `request body must be valid JSON` |
| 400 | Body is not an object | `request body must be a JSON object` |
| 400 | More than one JSON value | `request body must contain a single JSON object` |
| 400 | Body larger than 1 MiB | `request body is too large` |
| 405 | Any method other than POST or OPTIONS | `method not allowed` |
| 500 | Unexpected server failure | `internal server error` |

`OPTIONS /api/calculate` returns 204 and the CORS headers used for a cross-origin call. Responses set `Access-Control-Allow-Origin: *`.

## Tests

### Backend

From `backend`:

```powershell
go test ./...
```

Coverage:

```powershell
go test ./... -coverprofile coverage.out
go tool cover -func coverage.out
```

On PowerShell, keep a space between `-coverprofile` and the file name, and between `-func` and the file name. `coverage.out` is gitignored.

### Frontend

From `frontend`:

```powershell
npm test
npm run test:coverage
npm run typecheck
```

Vitest coverage uses the v8 provider. The report excludes `src/main.tsx` (the bootstrap), `src/types.ts` (type declarations only), `src/vite-env.d.ts`, and test files. `coverage/` is gitignored.

## Verified results

These commands were run in this workspace with Go 1.27.4, Node.js 20.18.0, and npm 10.8.2.

Backend:

```text
ok   sezzle-calculator/backend                       coverage: 0.0% of statements
ok   sezzle-calculator/backend/internal/calculator   coverage: 100.0% of statements
ok   sezzle-calculator/backend/internal/httpapi      coverage: 100.0% of statements
total: (statements) 94.9%
```

`main` is the uncovered package. It only resolves `PORT` and starts `http.Server`. Starting `go run .` logged `calculator API listening on :8080`.

Frontend:

```text
Test Files  3 passed (3)
     Tests  31 passed (31)

Statements   100%   (318/318)
Branches     97.77% (88/90)
Functions    100%   (20/20)
Lines        100%   (318/318)
```

The two unexercised branches are defensive: an operand that matches the decimal pattern but is not finite, and a thrown value that is not an `ApiError`. `npx tsc -b` completed through `npm run typecheck` and `npm run build` with exit code 0. The production build emitted `dist/index.html`, one CSS file, and one JS file.

`npm audit --omit=dev` reported 0 vulnerabilities. A full `npm audit` reported 7 vulnerabilities, all in the Vitest toolchain (`vitest`, `@vitest/mocker`, `tinypool`, `glob`, and `minimatch`). `npm audit fix --force` would upgrade Vitest to a major version this Node and Vite pin were not verified with, so it was not applied.

Live check: with both servers running, the page at http://localhost:5173/ displayed `60` for 12 × 5, `division by zero is not allowed` for 12 ÷ 0, `Enter both numbers.` for an empty submit, and `-6.25` for -12.5 × .5. A direct `POST` to http://127.0.0.1:8080/api/calculate returned `{"result":60}`. The layout was checked at desktop width and at a 390×844 viewport.

## Design decisions

- Arithmetic lives in `internal/calculator`. HTTP status codes and JSON live in `internal/httpapi`.
- The handler decodes the object into `map[string]json.RawMessage`. That preserves a missing number as an absent key, so `0` is not confused with a field that was left out, and `A` is not accepted as `a`. `encoding/json` otherwise matches struct fields without regard to case.
- Numeric text is parsed with `strconv.ParseFloat`. Overflow and a non-finite result become a 400 response, so `json.Marshal` is never asked to encode `NaN` or `Inf`.
- The form disables Calculate while a request is in flight and ignores a response that arrives after Clear.
- The UI accepts an optional sign and a decimal point. It does not treat `1e2` as input. The API still accepts exponent notation because that is a normal JSON number.

## Assumptions and limitations

- Results are `float64`. `0.1 + 0.2` is the IEEE-754 sum, not decimal `0.3`.
- Duplicate keys in one object follow Go's JSON decoder: the last value wins.
- There is no history, rounding mode, or thousands-separator input.
- A static host that does not proxy `/api` will not reach the API. Use `npm run dev`, `npm run preview`, or serve the build behind a proxy to port 8080.
- Exponentiation, square root, and percentage are intentionally absent.
- Git was not initialized. Coverage profiles and `frontend/dist` are generated locally and ignored.

## AI assistance

Cursor's Grok 4.7 agent wrote this implementation from the assessment prompt. The prompt is in [PROMPT.md](PROMPT.md), and how it was used is described in [AI_USAGE.md](AI_USAGE.md). Read the code and run the tests before submitting it.
