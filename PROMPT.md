# SEZZLE — SOFTWARE ENGINEER II TECHNICAL ASSESSMENT

You are acting as a senior full-stack software engineer with strong expertise in React, TypeScript, Go, REST APIs, software architecture, and automated testing.

Your task is to build a complete, production-quality but appropriately scoped full-stack calculator application for Sezzle's Software Engineer II technical assessment.

The assessment recommends approximately 2–4 hours of development. Prioritize correctness, clean architecture, maintainability, testability, and a polished user experience over unnecessary complexity.

Implement the solution directly in the currently opened `sezzle-calculator` workspace.

## 1. Technology Stack

**Frontend**
- React with TypeScript
- Vite
- Plain CSS or CSS Modules
- Vitest
- React Testing Library

**Backend**
- Go (installed locally)
- Go standard library (`net/http`)
- Go built-in testing package
- REST API returning JSON

**General**
- Git-compatible monorepo structure
- Minimal dependencies
- No database
- No authentication
- No Docker unless all mandatory requirements are completed
- No unnecessary frameworks or architectural abstractions

## 2. Project Structure

Use a clear, maintainable structure similar to:

```text
sezzle-calculator/
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── services/
│   │   ├── App.tsx
│   │   ├── App.css
│   │   └── main.tsx
│   ├── package.json
│   └── vite.config.ts
├── backend/
│   ├── internal/
│   │   ├── calculator/
│   │   └── httpapi/
│   ├── main.go
│   └── go.mod
├── README.md
├── AI_USAGE.md
└── .gitignore
```

You may adjust the structure if a simpler, more idiomatic solution is appropriate. Avoid unnecessary layers.

## 3. Functional Requirements

Implement the following mandatory arithmetic operations:

- Addition
- Subtraction
- Multiplication
- Division

Optional operations:
- Exponentiation
- Square root
- Percentage

Only implement optional operations after the mandatory functionality, testing, and documentation are complete.

All arithmetic calculations must be performed by the Go backend, not locally in the React frontend.

## 4. Backend — Go REST API

Implement a REST API using Go's standard library.

Primary endpoint:

`POST /api/calculate`

Example request:

```json
{
  "operation": "multiply",
  "a": 12,
  "b": 5
}
```

Successful response:

```json
{
  "result": 60
}
```

Support the operation identifiers:

- `add`
- `subtract`
- `multiply`
- `divide`

Implement:

1. JSON request parsing.
2. Strict input validation.
3. Required field validation, including distinguishing missing numeric fields from zero.
4. Unsupported operation handling.
5. Division-by-zero protection.
6. Rejection of malformed JSON and unexpected fields.
7. Rejection of multiple JSON objects in one request.
8. Appropriate HTTP status codes.
9. Consistent JSON error responses.
10. Clear separation between arithmetic logic and HTTP transport.
11. Proper CORS handling or Vite development proxy configuration.
12. Safe handling of non-finite arithmetic results so JSON encoding never fails unexpectedly.

Example error response:

```json
{
  "error": "division by zero is not allowed"
}
```

Use appropriate HTTP status codes such as 200, 400, 405, and 500 when applicable.

Avoid overengineering. Keep the backend understandable and idiomatic.

## 5. Frontend — React + TypeScript

Build a clean, modern, responsive calculator interface.

Design requirements:

- Professional, minimal visual design.
- Clear typography and spacing.
- Responsive layout for desktop and mobile.
- Accessible labels and keyboard navigation.
- Clearly visible result area.
- Consistent button styles and interaction states.
- Loading indicator while awaiting the API.
- User-friendly error messages.

Functional requirements:

1. Allow users to enter two numeric operands.
2. Allow users to select an arithmetic operation.
3. Submit the calculation to the Go API.
4. Display the result returned by the backend.
5. Validate empty and invalid inputs.
6. Prevent accidental duplicate submissions while loading.
7. Display errors returned by the API.
8. Handle backend connection failures gracefully.
9. Support negative numbers and decimal values.
10. Provide a clear/reset action.

Use React functional components and hooks.

Use TypeScript types for request and response structures.

Do not calculate arithmetic results directly in the frontend.

Do not add unnecessary state management libraries.

## 6. Automated Testing

Testing is mandatory for both layers.

### Backend Tests

Use Go's built-in testing package.

Cover:

- Addition
- Subtraction
- Multiplication
- Division
- Division by zero
- Unsupported operations
- Missing fields
- Malformed JSON
- Invalid JSON types
- Unexpected JSON fields
- Multiple JSON objects
- HTTP status codes
- JSON response structures
- Non-finite arithmetic results

Use table-driven tests where appropriate.

Generate coverage using:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Frontend Tests

Use Vitest and React Testing Library.

Cover:

- Initial rendering
- Operand input
- Operation selection
- Successful API calculation
- Input validation
- API error responses
- Network failures
- Loading behavior
- Clear/reset functionality

Mock API responses appropriately.

Generate frontend coverage using the appropriate Vitest coverage provider.

Do not invent test results or coverage percentages.

If tests fail, investigate and fix the actual cause.

## 7. Documentation — README.md

Create a professional, concise README in English.

Include:

1. Project overview.
2. Technology stack.
3. Architecture and folder structure.
4. Prerequisites.
5. Installation instructions.
6. How to start the Go backend.
7. How to start the React frontend.
8. Default local ports.
9. REST API documentation.
10. Request and response JSON examples.
11. Error handling examples.
12. How to run backend tests.
13. How to run frontend tests.
14. How to generate coverage reports.
15. Actual verified test results, if available.
16. Key design decisions.
17. Assumptions and limitations.
18. AI assistance disclosure.

All commands must be accurate and reproducible on a clean development environment.

Document any limitations honestly.

## 8. AI Usage Documentation

Create `AI_USAGE.md`.

The assessment explicitly allows AI tools and requires sharing the prompts used.

Document:

- AI tools used.
- Purpose of AI assistance.
- The actual initial prompt used to generate this project.
- A section for additional prompts used during debugging, testing, and refinement.
- A brief explanation that generated code must be reviewed and validated.

Do not fabricate prompts, development history, manual reviews, or actions that have not occurred.

Include the complete text of this prompt in the AI usage documentation, or reference a separate file containing it verbatim.

## 9. Code Quality

Follow these principles:

- Single responsibility.
- Clear naming conventions.
- Small, readable functions.
- Explicit error handling.
- Minimal duplication.
- Appropriate TypeScript types.
- Idiomatic Go.
- No unused dependencies.
- No unnecessary comments.
- No premature optimization.
- No excessive abstractions.

The result should demonstrate mid-level software engineering competence.

## 10. Security and Reliability

Ensure:

- No secrets or credentials are committed.
- Invalid inputs do not crash the backend.
- Division by zero is handled safely.
- Unsupported HTTP methods are rejected.
- Backend errors are communicated clearly.
- CORS or development proxy configuration works.
- The application handles backend unavailability.
- Generated coverage files and build artifacts are excluded from Git.

Do not introduce unnecessary security infrastructure for this small assessment.

## 11. Execution and Verification

After implementing the application:

1. Install frontend dependencies.
2. Format Go code.
3. Run Go tests.
4. Generate Go coverage.
5. Run frontend tests.
6. Generate frontend coverage.
7. Run TypeScript checks.
8. Run the frontend production build.
9. Verify the backend starts successfully.
10. Verify that the frontend can communicate with the backend.
11. Fix any errors found.

Use commands compatible with Windows PowerShell where appropriate.

Do not claim successful execution unless commands were actually run.

If any verification cannot be completed, clearly explain why.

## 12. Git and Submission Preparation

Prepare the repository for GitHub.

Create an appropriate `.gitignore`.

Exclude:

- `node_modules`
- Build output
- Coverage artifacts
- Environment secrets
- Temporary files
- Editor-specific unnecessary files

Do not initialize Git, create commits, create remote repositories, or push code automatically.

We will review the implementation and create meaningful commits manually afterward.

Do not invent development timestamps or claim the project was completed in a specific amount of time.

## 13. Final Report

When implementation and verification are complete, provide:

1. Summary of implemented features.
2. Final project structure.
3. Exact commands to run backend and frontend.
4. Actual backend test results.
5. Actual backend coverage percentage.
6. Actual frontend test results.
7. Actual frontend coverage percentage.
8. TypeScript/build verification results.
9. Remaining errors or limitations.
10. List of files created or modified.
11. Suggested logical Git commits.
12. Confirmation of whether all mandatory Sezzle requirements are satisfied.

Do not stop after generating files. Execute the verification steps and resolve issues wherever possible.

## FINAL PRIORITY

The objective is not to build the most complex calculator.

The objective is to deliver a clean, reliable, maintainable, well-tested full-stack application that demonstrates strong engineering fundamentals and can be confidently explained during a technical interview.

Begin implementation now.
