# AI usage

This project was implemented with an AI coding assistant. Generated code still has to be read, run, and judged by a person before it is treated as submission-ready. Passing tests and a production build are checks, not a substitute for that review.

## Tools

- Cursor, using the Grok 4.7 coding agent

No other AI tools were used.

## Purpose

The assistant was asked to implement the Sezzle Software Engineer II calculator assessment in this workspace: the Go API, the React client, tests, coverage, documentation, and the verification commands in the prompt.

## Initial prompt

The complete initial prompt is stored verbatim in [PROMPT.md](PROMPT.md). It is the only user prompt used for this work.

## Additional prompts

No additional prompts were used for debugging, testing, or refinement. Fixes made after the first test run, including stricter JSON parsing and the Node-compatible frontend toolchain, were done in the same session while carrying out the initial prompt.

## Review

The assistant ran the Go tests, the frontend tests, the TypeScript check, the production build, and a live browser session against the running API. Those results are recorded in [README.md](README.md). They do not replace a human review of the design, the edge cases, and the decision to leave the optional operations out.
