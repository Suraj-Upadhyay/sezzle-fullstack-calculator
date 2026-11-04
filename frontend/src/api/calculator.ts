export interface EvaluateResponse {
  expression: string;
  result: number;
}

export interface ErrorDetail {
  fields?: Record<string, string>;
  position?: number;
  reason?: string;
}

interface ErrorBody extends ErrorDetail {
  error: string;
  message: string;
}

export const NETWORK_ERROR = "err_network";
export const UNEXPECTED_RESPONSE = "err_unexpected_response";

export class ApiError extends Error {
  readonly code: string;
  readonly status: number;
  readonly detail: ErrorDetail;

  constructor(
    code: string,
    message: string,
    status: number,
    detail: ErrorDetail = {},
  ) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.status = status;
    this.detail = detail;
  }
}

export async function evaluate(
  expression: string,
  signal?: AbortSignal,
): Promise<EvaluateResponse> {
  let response: Response;
  try {
    response = await fetch("/api/v1/evaluate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ expression }),
      signal,
    });
  } catch (err) {
    if (signal?.aborted) throw err;
    throw new ApiError(
      NETWORK_ERROR,
      "Could not reach the calculator service.",
      0,
    );
  }

  const body: unknown = await response.json().catch(() => null);

  if (!response.ok) {
    if (isErrorBody(body)) {
      const { error, message, ...detail } = body;
      throw new ApiError(error, message, response.status, detail);
    }
    throw unexpectedResponse(response.status);
  }
  if (!isEvaluateResponse(body)) {
    throw unexpectedResponse(response.status);
  }
  return body;
}

function unexpectedResponse(status: number): ApiError {
  return new ApiError(
    UNEXPECTED_RESPONSE,
    `Unexpected response from the server (HTTP ${String(status)}).`,
    status,
  );
}

function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function isErrorBody(value: unknown): value is ErrorBody {
  return (
    isObject(value) &&
    typeof value.error === "string" &&
    typeof value.message === "string"
  );
}

function isEvaluateResponse(value: unknown): value is EvaluateResponse {
  return (
    isObject(value) &&
    typeof value.expression === "string" &&
    typeof value.result === "number"
  );
}
