import type { CalculateRequest, CalculateResponse, ErrorResponse } from "../types";

export class ApiError extends Error {
  readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export async function calculate(request: CalculateRequest): Promise<CalculateResponse> {
  let response: Response;
  try {
    response = await fetch("/api/calculate", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        operation: request.operation,
        a: request.a,
        b: request.b,
      }),
      cache: "no-store",
    });
  } catch {
    throw new ApiError(
      "Unable to reach the calculator service. Start the backend and try again.",
      0,
    );
  }

  const payload = await readJSON(response);
  if (!response.ok) {
    throw new ApiError(
      messageFrom(payload) ?? "The calculation could not be completed.",
      response.status,
    );
  }
  if (!isCalculateResponse(payload)) {
    throw new ApiError(
      "The calculator service returned an unexpected response.",
      response.status,
    );
  }
  return payload;
}

async function readJSON(response: Response): Promise<unknown> {
  let text: string;
  try {
    text = await response.text();
  } catch {
    throw new ApiError(
      "The calculator service returned an unreadable response.",
      response.status,
    );
  }
  if (text.trim() === "") {
    throw new ApiError(
      "The calculator service returned an unreadable response.",
      response.status,
    );
  }
  try {
    return JSON.parse(text) as unknown;
  } catch {
    throw new ApiError(
      "The calculator service returned an unreadable response.",
      response.status,
    );
  }
}

function messageFrom(payload: unknown): string | null {
  if (!isErrorResponse(payload)) {
    return null;
  }
  const message = payload.error.trim();
  return message === "" ? null : message;
}

function isErrorResponse(payload: unknown): payload is ErrorResponse {
  if (typeof payload !== "object" || payload === null || !("error" in payload)) {
    return false;
  }
  return typeof (payload as { error: unknown }).error === "string";
}

function isCalculateResponse(payload: unknown): payload is CalculateResponse {
  if (typeof payload !== "object" || payload === null || !("result" in payload)) {
    return false;
  }
  const result = (payload as { result: unknown }).result;
  return typeof result === "number" && Number.isFinite(result);
}
