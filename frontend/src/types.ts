export type Operation = "add" | "subtract" | "multiply" | "divide";

export interface CalculateRequest {
  operation: Operation;
  a: number;
  b: number;
}

export interface CalculateResponse {
  result: number;
}

export interface ErrorResponse {
  error: string;
}
