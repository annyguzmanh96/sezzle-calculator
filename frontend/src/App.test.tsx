import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { App } from "./App";

describe("App", () => {
  it("renders the calculator", () => {
    render(<App />);

    expect(screen.getByRole("heading", { name: "Sezzle Calculator" })).toBeInTheDocument();
    expect(
      screen.getByText("Add, subtract, multiply, and divide. Every result is computed by the API."),
    ).toBeInTheDocument();
    expect(screen.getByLabelText("First number")).toHaveValue("");
    expect(screen.getByLabelText("Second number")).toHaveValue("");
    expect(screen.getByRole("radio", { name: "Add" })).toBeChecked();
    expect(screen.getByRole("radio", { name: "Subtract" })).not.toBeChecked();
    expect(screen.getByRole("radio", { name: "Multiply" })).not.toBeChecked();
    expect(screen.getByRole("radio", { name: "Divide" })).not.toBeChecked();
    expect(screen.getByRole("button", { name: "Calculate" })).toBeEnabled();
    expect(screen.getByRole("button", { name: "Clear" })).toBeEnabled();
    expect(screen.getByTestId("result")).toHaveTextContent("—");
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });
});
