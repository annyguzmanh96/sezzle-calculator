import { beforeEach, describe, expect, it, vi } from "vitest";
import { ApiError, calculate } from "./calculate";

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

describe("calculate", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  it("returns the finite result from the API", async () => {
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockResolvedValueOnce(jsonResponse({ result: 60 }));

    await expect(calculate({ operation: "multiply", a: 12, b: 5 })).resolves.toEqual({
      result: 60,
    });
    expect(fetchMock).toHaveBeenCalledWith("/api/calculate", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ operation: "multiply", a: 12, b: 5 }),
      cache: "no-store",
    });
  });

  it("uses the API error message and status", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(
      jsonResponse({ error: "  division by zero is not allowed  " }, 400),
    );

    const error = await calculate({ operation: "divide", a: 1, b: 0 }).catch((caught: unknown) => caught);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      message: "division by zero is not allowed",
      status: 400,
    });
  });

  it("reports a network failure", async () => {
    vi.mocked(fetch).mockRejectedValueOnce(new TypeError("Failed to fetch"));

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "Unable to reach the calculator service. Start the backend and try again.",
      status: 0,
    });
  });

  it("reports an unreadable response", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(
      new Response("<html>nope</html>", { status: 502 }),
    );

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculator service returned an unreadable response.",
      status: 502,
    });
  });

  it("reports an empty response", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(new Response("   ", { status: 500 }));

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculator service returned an unreadable response.",
      status: 500,
    });
  });

  it("reports an unexpected success payload", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: "60" }));

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculator service returned an unexpected response.",
      status: 200,
    });
  });

  it("rejects a null result", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(new Response('{"result":null}', { status: 200 }));

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculator service returned an unexpected response.",
    });
  });

  it("reports a body that cannot be read", async () => {
    const response = new Response("{}", { status: 502 });
    vi.spyOn(response, "text").mockRejectedValueOnce(new Error("read failed"));
    vi.mocked(fetch).mockResolvedValueOnce(response);

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculator service returned an unreadable response.",
      status: 502,
    });
  });

  it("uses a generic message when the payload is not an error object", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ detail: "nope" }, 400));

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculation could not be completed.",
      status: 400,
    });
  });

  it("rejects a success payload without a result", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ ok: true }));

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculator service returned an unexpected response.",
      status: 200,
    });
  });

  it("uses a generic message when the error payload has no message", async () => {
    vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ error: "   " }, 400));

    await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toMatchObject({
      message: "The calculation could not be completed.",
      status: 400,
    });
  });
});
