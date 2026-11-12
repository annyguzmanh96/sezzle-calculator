import { useRef, useState, type FormEvent } from "react";
import { ApiError, calculate } from "../services/calculate";
import type { Operation } from "../types";
import styles from "./Calculator.module.css";

const SYMBOLS: Record<Operation, string> = {
  add: "+",
  subtract: "−",
  multiply: "×",
  divide: "÷",
};

const OPERATIONS: { id: Operation; label: string; symbol: string }[] = (
  Object.keys(SYMBOLS) as Operation[]
).map((id) => ({
  id,
  label: id.charAt(0).toUpperCase() + id.slice(1),
  symbol: SYMBOLS[id],
}));

const OPERAND_PATTERN = /^[+-]?(?:\d+\.?\d*|\.\d+)$/;

function parseOperand(input: string): number | null {
  const value = input.trim();
  if (value === "") {
    return null;
  }
  if (!OPERAND_PATTERN.test(value)) {
    return Number.NaN;
  }
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : Number.NaN;
}

function isCompleteOperand(input: string): boolean {
  const parsed = parseOperand(input);
  return parsed !== null && !Number.isNaN(parsed);
}

function operationSymbol(operation: Operation): string {
  return SYMBOLS[operation];
}

export function Calculator() {
  const [a, setA] = useState("");
  const [b, setB] = useState("");
  const [operation, setOperation] = useState<Operation>("add");
  const [result, setResult] = useState<number | null>(null);
  const [expression, setExpression] = useState<string | null>(null);
  const [fieldError, setFieldError] = useState<string | null>(null);
  const [requestError, setRequestError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const requestId = useRef(0);
  const busy = useRef(false);

  const message = fieldError ?? requestError;
  const aInvalid = fieldError !== null && !isCompleteOperand(a);
  const bInvalid = fieldError !== null && !isCompleteOperand(b);

  function updateA(value: string) {
    setA(value);
    setFieldError(null);
    setRequestError(null);
  }

  function updateB(value: string) {
    setB(value);
    setFieldError(null);
    setRequestError(null);
  }

  function updateOperation(value: Operation) {
    setOperation(value);
    setFieldError(null);
    setRequestError(null);
  }

  function reset() {
    requestId.current += 1;
    busy.current = false;
    setA("");
    setB("");
    setOperation("add");
    setResult(null);
    setExpression(null);
    setFieldError(null);
    setRequestError(null);
    setLoading(false);
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busy.current) {
      return;
    }

    const left = parseOperand(a);
    const right = parseOperand(b);
    if (left === null || right === null) {
      setFieldError("Enter both numbers.");
      setRequestError(null);
      setResult(null);
      setExpression(null);
      return;
    }
    if (Number.isNaN(left) || Number.isNaN(right)) {
      setFieldError("Enter valid numbers. Decimals and negative values are allowed.");
      setRequestError(null);
      setResult(null);
      setExpression(null);
      return;
    }

    const id = ++requestId.current;
    busy.current = true;
    setLoading(true);
    setFieldError(null);
    setRequestError(null);

    try {
      const response = await calculate({ operation, a: left, b: right });
      if (id !== requestId.current) {
        return;
      }
      setResult(response.result);
      setExpression(`${a.trim()} ${operationSymbol(operation)} ${b.trim()}`);
    } catch (caught) {
      if (id !== requestId.current) {
        return;
      }
      setResult(null);
      setExpression(null);
      setRequestError(
        caught instanceof ApiError ? caught.message : "Something went wrong. Try again.",
      );
    } finally {
      if (id === requestId.current) {
        busy.current = false;
        setLoading(false);
      }
    }
  }

  return (
    <form
      className={styles.card}
      aria-label="Calculator"
      aria-busy={loading}
      noValidate
      onSubmit={onSubmit}
    >
      <div className={styles.display} role="status" aria-live="polite" aria-atomic="true">
        <p className={styles.displayLabel}>Result</p>
        <p className={styles.expression}>{!loading && expression ? expression : "\u00a0"}</p>
        <p className={styles.displayValue} data-testid="result">
          {loading ? "Calculating..." : result === null ? "—" : String(result)}
        </p>
      </div>

      {message ? (
        <p id="calculator-error" className={styles.error} role="alert">
          {message}
        </p>
      ) : null}

      <div className={styles.operands}>
        <div className={styles.field}>
          <label htmlFor="operand-a">First number</label>
          <input
            id="operand-a"
            name="a"
            type="text"
            inputMode="decimal"
            autoComplete="off"
            spellCheck={false}
            maxLength={100}
            value={a}
            disabled={loading}
            aria-invalid={aInvalid}
            aria-describedby={fieldError ? "calculator-error" : undefined}
            onChange={(event) => updateA(event.target.value)}
          />
        </div>
        <div className={styles.field}>
          <label htmlFor="operand-b">Second number</label>
          <input
            id="operand-b"
            name="b"
            type="text"
            inputMode="decimal"
            autoComplete="off"
            spellCheck={false}
            maxLength={100}
            value={b}
            disabled={loading}
            aria-invalid={bInvalid}
            aria-describedby={fieldError ? "calculator-error" : undefined}
            onChange={(event) => updateB(event.target.value)}
          />
        </div>
      </div>

      <fieldset className={styles.operations} disabled={loading}>
        <legend>Operation</legend>
        <div className={styles.operationGrid}>
          {OPERATIONS.map((item) => (
            <label key={item.id} className={styles.operation}>
              <input
                type="radio"
                name="operation"
                value={item.id}
                aria-label={item.label}
                checked={operation === item.id}
                onChange={() => updateOperation(item.id)}
              />
              <span className={styles.symbol} aria-hidden="true">
                {item.symbol}
              </span>
              <span className={styles.name}>{item.label}</span>
            </label>
          ))}
        </div>
      </fieldset>

      <div className={styles.actions}>
        <button type="submit" className={styles.primary} disabled={loading}>
          {loading ? <span className={styles.spinner} aria-hidden="true" /> : null}
          {loading ? "Calculating..." : "Calculate"}
        </button>
        <button type="button" className={styles.secondary} onClick={reset}>
          Clear
        </button>
      </div>
      <p className={styles.note}>Computed by the Go API.</p>
    </form>
  );
}
