import { describe, expect, it } from "vitest";

import { errorMessage, formatNumber, toExpression, toSentence } from "./format";

describe("formatNumber", () => {
  it.each([
    [0, "0"],
    [-0, "0"],
    [0.1 + 0.2, "0.3"],
    [1 / 3, "0.333333333333333"],
    [1234567.5, "1,234,567.5"],
    [-42, "-42"],
    [1e21, "1E21"],
    [1.5e-12, "1.5E-12"],
  ])("formats %d as %s", (value, expected) => {
    expect(formatNumber(value)).toBe(expected);
  });
});

describe("toExpression", () => {
  it("drops floating-point noise and uses the keypad's minus sign", () => {
    expect(toExpression(0.1 + 0.2)).toBe("0.3");
    expect(toExpression(-2.5)).toBe("−2.5");
    expect(toExpression(1e21)).toBe("1e+21");
    expect(toExpression(-1e-7)).toBe("−1e-7");
  });
});

describe("toSentence", () => {
  it("capitalises and adds a full stop", () => {
    expect(toSentence("division by zero is undefined")).toBe(
      "Division by zero is undefined.",
    );
  });

  it("keeps existing punctuation", () => {
    expect(toSentence("Could not connect.")).toBe("Could not connect.");
  });

  it("leaves empty text alone", () => {
    expect(toSentence("  ")).toBe("");
  });
});

describe("errorMessage", () => {
  it("appends the reason when there is one", () => {
    expect(errorMessage("expression is not valid", 'missing ")"')).toBe(
      'Expression is not valid: missing ")".',
    );
    expect(errorMessage("division by zero is undefined")).toBe(
      "Division by zero is undefined.",
    );
  });
});
