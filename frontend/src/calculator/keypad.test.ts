import { describe, expect, it } from "vitest";

import { LAYOUT } from "./keypad";

describe("keypad layout", () => {
  it("fills a 4-column grid exactly", () => {
    const cells = LAYOUT.reduce((n, spec) => n + (spec.wide ? 2 : 1), 0);
    expect(cells % 4).toBe(0);
  });

  it("has every digit once", () => {
    const digits = LAYOUT.filter((spec) => spec.variant === "digit")
      .map((spec) => spec.key)
      .filter((key) => /\d/.test(key))
      .sort();
    expect(digits).toEqual(Array.from("0123456789"));
  });
});
