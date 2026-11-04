import { vi } from "vitest";

export function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

export function stubFetch(
  impl: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>,
) {
  const fetchMock = vi.fn(impl);
  vi.stubGlobal("fetch", fetchMock);
  return fetchMock;
}

export function parseBody(init?: RequestInit): unknown {
  if (typeof init?.body !== "string") throw new Error("body is not JSON text");
  return JSON.parse(init.body);
}

export function requestBody(fetchMock: ReturnType<typeof stubFetch>): unknown {
  return parseBody(fetchMock.mock.calls.at(-1)?.[1]);
}
