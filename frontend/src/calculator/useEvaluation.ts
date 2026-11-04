import { useCallback, useEffect, useRef, useState } from "react";

import { ApiError, evaluate, type EvaluateResponse } from "../api/calculator";

export type EvaluationState =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "success"; response: EvaluateResponse }
  | { status: "error"; error: ApiError };

export function useEvaluation() {
  const [state, setState] = useState<EvaluationState>({ status: "idle" });
  const inFlight = useRef<AbortController | null>(null);

  useEffect(() => () => inFlight.current?.abort(), []);

  const reset = useCallback(() => {
    inFlight.current?.abort();
    inFlight.current = null;
    setState({ status: "idle" });
  }, []);

  const run = useCallback(
    async (expression: string): Promise<EvaluateResponse | null> => {
      inFlight.current?.abort();
      const controller = new AbortController();
      inFlight.current = controller;
      setState({ status: "loading" });

      try {
        const response = await evaluate(expression, controller.signal);
        setState({ status: "success", response });
        return response;
      } catch (err) {
        if (controller.signal.aborted) return null;
        setState({ status: "error", error: toApiError(err) });
        return null;
      }
    },
    [],
  );

  return { state, run, reset };
}

function toApiError(err: unknown): ApiError {
  if (err instanceof ApiError) return err;
  return new ApiError("err_unknown", "Something went wrong.", 0);
}
