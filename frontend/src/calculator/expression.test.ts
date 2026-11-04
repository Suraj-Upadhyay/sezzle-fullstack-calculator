import { describe, expect, it } from "vitest";

import {
  applyKey,
  closeBrackets,
  isComplete,
  keyFromKeyboard,
  MAX_LENGTH,
  openBrackets,
  type EditKey,
} from "./expression";

const press = (keys: string, start = ""): string =>
  Array.from(keys)
    .map((char) => (char === "<" ? "backspace" : char) as EditKey)
    .reduce(applyKey, start);

describe("applyKey", () => {
  it.each([
    ["12+3", "12+3"],
    ["2×(3+4)", "2×(3+4)"],
    ["√16", "√16"],
    ["50%", "50%"],
    ["2^−1", "2^−1"],
    ["−5", "−5"],
    ["(−2)", "(−2)"],
    ["3.25", "3.25"],
  ])("builds %s", (keys, expected) => {
    expect(press(keys)).toBe(expected);
  });

  it.each([
    ["operators replace each other", "2+×", "2×"],
    ["minus after × is a sign", "2×−", "2×−"],
    ["an operator replaces a trailing sign", "2×−+", "2+"],
    ["no double minus", "2−−", "2−"],
    ["no leading binary operator", "×", ""],
    ["only a sign after (", "(×", "("],
    ["no second decimal point", "1.2.", "1.2"],
    ["a bare decimal point gets a 0", ".5", "0.5"],
    ["a decimal point after an operator gets a 0", "1+.", "1+0."],
    ["leading zeros collapse", "007", "7"],
    ["zeros after a decimal point stay", "0.0", "0.0"],
    ["no ) without an open (", "1)", "1"],
    ["no ) straight after an operator", "(1+)", "(1+"],
    ["no ) straight after (", "()", "("],
    ["no % without a value", "%", ""],
    ["no double %", "5%%", "5%"],
    ["( after a value multiplies", "2(", "2×("],
    ["√ after a value multiplies", "2√", "2×√"],
    ["a digit after ) multiplies", "(2)3", "(2)×3"],
    ["a decimal point after % multiplies", "5%.", "5%×0."],
    ["backspace removes one character", "12+<", "12"],
    ["backspace on empty input", "<", ""],
  ])("%s", (_, keys, expected) => {
    expect(press(keys)).toBe(expected);
  });

  it("does not add a decimal point to a number in e-notation", () => {
    expect(press(".", "1e+21")).toBe("1e+21");
  });

  it("stops growing at the maximum length", () => {
    const full = "1".repeat(MAX_LENGTH);
    expect(applyKey(full, "2")).toBe(full);
    expect(applyKey(full, "backspace")).toHaveLength(MAX_LENGTH - 1);
  });
});

describe("brackets", () => {
  it("counts unclosed brackets", () => {
    expect(openBrackets("((1)+(2")).toBe(2);
  });

  it("closes what is open", () => {
    expect(closeBrackets("2×(3+(4")).toBe("2×(3+(4))");
    expect(closeBrackets("2")).toBe("2");
  });
});

describe("isComplete", () => {
  it.each([
    ["2+3", true],
    ["(2+3", true],
    ["5%", true],
    ["2.", true],
    ["", false],
    ["2+", false],
    ["(", false],
    ["√", false],
  ])("%j is %s", (expression, expected) => {
    expect(isComplete(expression)).toBe(expected);
  });
});

describe("keyFromKeyboard", () => {
  it.each([
    ["7", "7"],
    ["+", "+"],
    ["-", "−"],
    ["*", "×"],
    ["x", "×"],
    ["/", "÷"],
    [",", "."],
    ["^", "^"],
    ["%", "%"],
    ["(", "("],
    ["r", "√"],
    ["Enter", "="],
    ["=", "="],
    ["Backspace", "backspace"],
    ["Escape", "clear"],
    ["Delete", "clear"],
  ])("maps %j to %j", (key, expected) => {
    expect(keyFromKeyboard(key)).toBe(expected);
  });

  it.each(["a", "Tab", "F5", " ", "ArrowLeft"])("ignores %j", (key) => {
    expect(keyFromKeyboard(key)).toBeNull();
  });
});
