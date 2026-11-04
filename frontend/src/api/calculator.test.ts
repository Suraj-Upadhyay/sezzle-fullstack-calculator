import { describe, expect, it } from "vitest";

import { jsonResponse, requestBody, stubFetch } from "../test/fetch";
import {
  ApiError,
  evaluate,
  NETWORK_ERROR,
  UNEXPECTED_RESPONSE,
} from "./calculator";

describe("evaluate", () => {
  it("posts the expression and returns the result", async () => {
    const response = { expression: "(2+3)×4", result: 20 };
    const fetchMock = stubFetch(() =>
      Promise.resolve(jsonResponse(200, response)),
    );

    await expect(evaluate("(2+3)×4")).resolves.toEqual(response);

    const [url, init] = fetchMock.mock.calls[0] ?? [];
    expect(url).toBe("/api/v1/evaluate");
    expect(init?.method).toBe("POST");
    expect(requestBody(fetchMock)).toEqual({ expression: "(2+3)×4" });
  });

  it("raises the server's error code, message and detail", async () => {
    stubFetch(() =>
      Promise.resolve(
        jsonResponse(400, {
          error: "err_invalid_expression",
          message: "expression is not valid",
          position: 5,
          reason: 'missing ")"',
        }),
      ),
    );

    const error = await evaluate("(1+2").catch((e: unknown) => e);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      code: "err_invalid_expression",
      status: 400,
      detail: { position: 5, reason: 'missing ")"' },
    });
  });

  it("reports a network failure", async () => {
    stubFetch(() => Promise.reject(new TypeError("Failed to fetch")));

    await expect(evaluate("1+2")).rejects.toMatchObject({
      code: NETWORK_ERROR,
      status: 0,
    });
  });

  it.each([
    [
      "a non-JSON error",
      new Response("<html>bad gateway</html>", { status: 502 }),
    ],
    ["a success without a result", jsonResponse(200, { ok: true })],
  ])("rejects %s", async (_, response) => {
    stubFetch(() => Promise.resolve(response));

    await expect(evaluate("1+2")).rejects.toMatchObject({
      code: UNEXPECTED_RESPONSE,
    });
  });

  it("rethrows an abort instead of reporting a network error", async () => {
    const controller = new AbortController();
    controller.abort();
    const abort = new DOMException("aborted", "AbortError");
    stubFetch(() => Promise.reject(abort));

    await expect(evaluate("1+2", controller.signal)).rejects.toBe(abort);
  });
});
