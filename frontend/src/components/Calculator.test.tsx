import { act, fireEvent, render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { jsonResponse, parseBody, requestBody, stubFetch } from "../test/fetch";
import { Calculator } from "./Calculator";

const RESULTS: Record<string, number> = {
  "1+2": 3,
  "2×(3+4)": 14,
  "(2+3)": 5,
  "√16": 4,
  "3×2": 6,
  "6−1": 5,
  "10÷4": 2.5,
};

function stubServer() {
  return stubFetch((_, init) => {
    const { expression } = parseBody(init) as { expression: string };
    const result = RESULTS[expression];
    return Promise.resolve(
      result === undefined
        ? jsonResponse(400, {
            error: "err_invalid_expression",
            message: "expression is not valid",
            reason: `no stub for ${expression}`,
          })
        : jsonResponse(200, { expression, result }),
    );
  });
}

function setup() {
  const user = userEvent.setup();
  render(<Calculator />);
  const display = () => screen.getByRole("region", { name: "Display" });
  return {
    user,
    display,
    input: () => within(display()).getByLabelText("Input"),
    result: () => within(display()).findByLabelText("Result"),
    key: (name: string) => screen.getByRole("button", { name }),
  };
}

describe("Calculator", () => {
  it("evaluates what is typed on the keyboard", async () => {
    const fetchMock = stubServer();
    const { user, display, result } = setup();

    await user.keyboard("2*(3+4)");
    expect(within(display()).getByLabelText("Input")).toHaveTextContent(
      "2×(3+4)",
    );

    await user.keyboard("{Enter}");

    expect(await result()).toHaveTextContent("14");
    expect(within(display()).getByLabelText("Expression")).toHaveTextContent(
      "2×(3+4) =",
    );
    expect(requestBody(fetchMock)).toEqual({ expression: "2×(3+4)" });
  });

  it("evaluates what is tapped on the keypad", async () => {
    stubServer();
    const { user, key, result } = setup();

    for (const name of ["Square root", "1", "6", "Equals"]) {
      await user.click(key(name));
    }

    expect(await result()).toHaveTextContent("4");
  });

  it("shows a key as pressed when its keyboard key goes down", () => {
    vi.useFakeTimers();
    try {
      const { key } = setup();

      fireEvent.keyDown(window, { key: "7" });
      expect(key("7")).toHaveAttribute("data-pressed", "true");
      fireEvent.keyDown(window, { key: "*" });
      expect(key("Multiply")).toHaveAttribute("data-pressed", "true");
      expect(key("7")).toHaveAttribute("data-pressed", "false");

      act(() => {
        vi.advanceTimersByTime(200);
      });
      expect(key("Multiply")).toHaveAttribute("data-pressed", "false");
    } finally {
      vi.useRealTimers();
    }
  });

  it("leaves browser shortcuts alone", () => {
    const { input } = setup();

    fireEvent.keyDown(window, { key: "r", ctrlKey: true });
    fireEvent.keyDown(window, { key: "1", metaKey: true });

    expect(input()).toHaveTextContent("0");
  });

  it("closes open brackets before evaluating", async () => {
    const fetchMock = stubServer();
    const { user, result } = setup();

    await user.keyboard("(2+3{Enter}");

    expect(await result()).toHaveTextContent("5");
    expect(requestBody(fetchMock)).toEqual({ expression: "(2+3)" });
  });

  it("asks for a complete expression without calling the API", async () => {
    const fetchMock = stubServer();
    const { user, display } = setup();

    await user.keyboard("{Enter}");
    expect(within(display()).getByRole("alert")).toHaveTextContent(
      "Enter an expression.",
    );

    await user.keyboard("2+{Enter}");
    expect(within(display()).getByRole("alert")).toHaveTextContent(
      "Finish the expression.",
    );
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("supports backspace and clear", async () => {
    stubServer();
    const { user, input, key } = setup();

    await user.keyboard("12+3{Backspace}");
    expect(input()).toHaveTextContent("12+");

    await user.click(key("Backspace"));
    expect(input()).toHaveTextContent("12");

    await user.keyboard("{Escape}");
    expect(input()).toHaveTextContent("0");
  });

  it("continues from a result with an operator", async () => {
    const fetchMock = stubServer();
    const { user, result } = setup();

    await user.keyboard("1+2{Enter}");
    await result();
    await user.keyboard("*2{Enter}");

    expect(await result()).toHaveTextContent("6");
    expect(requestBody(fetchMock)).toEqual({ expression: "3×2" });
  });

  it("starts over when a digit follows a result", async () => {
    stubServer();
    const { user, input, result } = setup();

    await user.keyboard("1+2{Enter}");
    await result();
    await user.keyboard("9");

    expect(input()).toHaveTextContent("9");
  });

  it("goes back to editing on backspace after a result", async () => {
    stubServer();
    const { user, input, result } = setup();

    await user.keyboard("1+2{Enter}");
    await result();
    await user.keyboard("{Backspace}");

    expect(input()).toHaveTextContent("1+");
  });

  it("shows the server's error, such as division by zero", async () => {
    stubFetch(() =>
      Promise.resolve(
        jsonResponse(422, {
          error: "err_division_by_zero",
          message: "division by zero is undefined",
        }),
      ),
    );
    const { user, display, input } = setup();

    await user.keyboard("5/0{Enter}");

    expect(await within(display()).findByRole("alert")).toHaveTextContent(
      "Division by zero is undefined.",
    );

    await user.keyboard("{Backspace}2");
    expect(input()).toHaveTextContent("5÷2");
    expect(within(display()).queryByRole("alert")).not.toBeInTheDocument();
  });

  it("explains when the service is unreachable", async () => {
    stubFetch(() => Promise.reject(new TypeError("Failed to fetch")));
    const { user, display } = setup();

    await user.keyboard("1+2{Enter}");

    expect(await within(display()).findByRole("alert")).toHaveTextContent(
      "Could not reach the calculator service.",
    );
  });

  it("accepts a pasted expression, keeping only what the keypad allows", () => {
    const { input } = setup();

    const paste = new Event("paste", { bubbles: true, cancelable: true });
    Object.assign(paste, {
      clipboardData: { getData: () => "2 * (3 + abc4)" },
    });
    act(() => {
      window.dispatchEvent(paste);
    });

    expect(input()).toHaveTextContent("2×(3+4)");
    expect(paste.defaultPrevented).toBe(true);
  });

  it("keeps a history and restores an expression from it", async () => {
    stubServer();
    const { user, input, result } = setup();

    await user.keyboard("10/4{Enter}");
    await result();
    await user.keyboard("6-1{Enter}");
    await result();

    const history = screen.getByRole("region", { name: "History" });
    const entries = within(history).getAllByRole("listitem");
    expect(entries.map((li) => li.textContent)).toEqual([
      "6−1= 5",
      "10÷4= 2.5",
    ]);

    await user.click(within(history).getByRole("button", { name: /10÷4/ }));
    expect(input()).toHaveTextContent("10÷4");

    await user.click(
      within(history).getByRole("button", { name: "Clear history" }),
    );
    expect(
      screen.queryByRole("region", { name: "History" }),
    ).not.toBeInTheDocument();
  });
});

describe("Keypad", () => {
  it("keeps focus where it was when a key is clicked", async () => {
    const { user, key } = setup();

    await user.click(key("7"));

    expect(key("7")).not.toHaveFocus();
  });
});
