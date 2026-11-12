import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { Calculator } from "./Calculator";

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

describe("Calculator", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  it("accepts operand input", async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "-12.5");
    await user.type(screen.getByLabelText("Second number"), ".5");

    expect(screen.getByLabelText("First number")).toHaveValue("-12.5");
    expect(screen.getByLabelText("Second number")).toHaveValue(".5");
  });

  it("selects an operation", async () => {
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("radio", { name: "Multiply" }));

    expect(screen.getByRole("radio", { name: "Multiply" })).toBeChecked();
    expect(screen.getByRole("radio", { name: "Add" })).not.toBeChecked();
  });

  it("displays the result returned by the API", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockResolvedValueOnce(jsonResponse({ result: 60 }));
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "12");
    await user.type(screen.getByLabelText("Second number"), "5");
    await user.click(screen.getByRole("radio", { name: "Multiply" }));
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByText("12 × 5")).toBeInTheDocument();
    expect(screen.getByTestId("result")).toHaveTextContent("60");
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

  it("shows the server result instead of calculating locally", async () => {
    const user = userEvent.setup();
    vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: 5 }));
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "2");
    await user.type(screen.getByLabelText("Second number"), "2");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByTestId("result")).toHaveTextContent("5");
    expect(screen.getByText("2 + 2")).toBeInTheDocument();
  });

  it("sends negative and decimal operands", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockResolvedValueOnce(jsonResponse({ result: -6.25 }));
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "-12.5");
    await user.type(screen.getByLabelText("Second number"), ".5");
    await user.click(screen.getByRole("radio", { name: "Multiply" }));
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByTestId("result")).toHaveTextContent("-6.25");
    expect(fetchMock).toHaveBeenCalledWith(
      "/api/calculate",
      expect.objectContaining({
        body: JSON.stringify({ operation: "multiply", a: -12.5, b: 0.5 }),
      }),
    );
  });

  it.each([
    ["", "", "Enter both numbers."],
    ["4", "", "Enter both numbers."],
    ["", "4", "Enter both numbers."],
    ["   ", "1", "Enter both numbers."],
    ["abc", "1", "Enter valid numbers. Decimals and negative values are allowed."],
    ["1.2.3", "1", "Enter valid numbers. Decimals and negative values are allowed."],
    ["-", "1", "Enter valid numbers. Decimals and negative values are allowed."],
    ["1e2", "1", "Enter valid numbers. Decimals and negative values are allowed."],
  ])("rejects operands %j and %j", async (left, right, message) => {
    const user = userEvent.setup();
    const fetchMock = vi.mocked(fetch);
    render(<Calculator />);

    if (left !== "") {
      await user.type(screen.getByLabelText("First number"), left);
    }
    if (right !== "") {
      await user.type(screen.getByLabelText("Second number"), right);
    }
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(message);
    expect(fetchMock).not.toHaveBeenCalled();
    expect(screen.getByTestId("result")).toHaveTextContent("—");
  });

  it("displays an API error", async () => {
    const user = userEvent.setup();
    vi.mocked(fetch).mockResolvedValueOnce(
      jsonResponse({ error: "division by zero is not allowed" }, 400),
    );
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "8");
    await user.type(screen.getByLabelText("Second number"), "0");
    await user.click(screen.getByRole("radio", { name: "Divide" }));
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "division by zero is not allowed",
    );
    expect(screen.getByTestId("result")).toHaveTextContent("—");
  });

  it("displays a network failure", async () => {
    const user = userEvent.setup();
    vi.mocked(fetch).mockRejectedValueOnce(new TypeError("Failed to fetch"));
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "1");
    await user.type(screen.getByLabelText("Second number"), "2");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Unable to reach the calculator service. Start the backend and try again.",
    );
  });

  it("shows a loading state and ignores a duplicate submit", async () => {
    const user = userEvent.setup();
    let resolveFetch: (response: Response) => void = () => {};
    const pending = new Promise<Response>((resolve) => {
      resolveFetch = resolve;
    });
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockReturnValueOnce(pending);
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "10");
    await user.type(screen.getByLabelText("Second number"), "4");
    await user.click(screen.getByRole("button", { name: "Calculate" }));

    expect(screen.getByRole("button", { name: "Calculating..." })).toBeDisabled();
    expect(screen.getByLabelText("First number")).toBeDisabled();
    expect(screen.getByTestId("result")).toHaveTextContent("Calculating...");
    expect(fetchMock).toHaveBeenCalledTimes(1);

    fireEvent.submit(screen.getByRole("form", { name: "Calculator" }));
    expect(fetchMock).toHaveBeenCalledTimes(1);

    resolveFetch(jsonResponse({ result: 14 }));
    expect(await screen.findByText("14")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Calculate" })).toBeEnabled();
  });

  it("ignores a response that arrives after clear", async () => {
    const user = userEvent.setup();
    let resolveFetch: (response: Response) => void = () => {};
    vi.mocked(fetch).mockReturnValueOnce(
      new Promise<Response>((resolve) => {
        resolveFetch = resolve;
      }),
    );
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "9");
    await user.type(screen.getByLabelText("Second number"), "3");
    await user.click(screen.getByRole("radio", { name: "Divide" }));
    await user.click(screen.getByRole("button", { name: "Calculate" }));
    await user.click(screen.getByRole("button", { name: "Clear" }));

    expect(screen.getByLabelText("First number")).toHaveValue("");
    expect(screen.getByRole("radio", { name: "Add" })).toBeChecked();
    expect(screen.getByTestId("result")).toHaveTextContent("—");

    await act(async () => {
      resolveFetch(jsonResponse({ result: 99 }));
    });

    await waitFor(() => {
      expect(screen.getByTestId("result")).toHaveTextContent("—");
    });
    expect(screen.queryByText("99")).not.toBeInTheDocument();
  });

  it("ignores an error that arrives after clear", async () => {
    const user = userEvent.setup();
    let rejectFetch: (reason: unknown) => void = () => {};
    vi.mocked(fetch).mockReturnValueOnce(
      new Promise<Response>((_, reject) => {
        rejectFetch = reject;
      }),
    );
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "4");
    await user.type(screen.getByLabelText("Second number"), "2");
    await user.click(screen.getByRole("button", { name: "Calculate" }));
    await user.click(screen.getByRole("button", { name: "Clear" }));

    await act(async () => {
      rejectFetch(new TypeError("Failed to fetch"));
    });

    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
    expect(screen.getByTestId("result")).toHaveTextContent("—");
    expect(screen.getByLabelText("First number")).toHaveValue("");
  });

  it("clears inputs, the result, the operation, and errors", async () => {
    const user = userEvent.setup();
    vi.mocked(fetch).mockResolvedValueOnce(
      jsonResponse({ error: "division by zero is not allowed" }, 400),
    );
    render(<Calculator />);

    await user.type(screen.getByLabelText("First number"), "8");
    await user.type(screen.getByLabelText("Second number"), "0");
    await user.click(screen.getByRole("radio", { name: "Divide" }));
    await user.click(screen.getByRole("button", { name: "Calculate" }));
    expect(await screen.findByRole("alert")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Clear" }));

    expect(screen.getByLabelText("First number")).toHaveValue("");
    expect(screen.getByLabelText("Second number")).toHaveValue("");
    expect(screen.getByRole("radio", { name: "Add" })).toBeChecked();
    expect(screen.getByTestId("result")).toHaveTextContent("—");
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });
});
