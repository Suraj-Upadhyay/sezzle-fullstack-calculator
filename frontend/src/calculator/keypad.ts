import type { Key } from "./expression";

export type Variant = "digit" | "operator" | "function" | "equals";

export interface KeySpec {
  key: Key;
  label: string;
  name: string;
  variant: Variant;
  wide?: boolean;
}

const k = (
  key: Key,
  label: string,
  name: string,
  variant: Variant,
  wide = false,
): KeySpec => ({ key, label, name, variant, wide });

const digit = (d: Key) => k(d, d, d, "digit");

export const LAYOUT: readonly KeySpec[] = [
  k("clear", "AC", "Clear", "function"),
  k("backspace", "⌫", "Backspace", "function"),
  k("(", "(", "Open bracket", "function"),
  k(")", ")", "Close bracket", "function"),
  k("√", "√", "Square root", "function"),
  k("^", "xʸ", "Power", "function"),
  k("%", "%", "Percent", "function"),
  k("÷", "÷", "Divide", "operator"),
  digit("7"),
  digit("8"),
  digit("9"),
  k("×", "×", "Multiply", "operator"),
  digit("4"),
  digit("5"),
  digit("6"),
  k("−", "−", "Subtract", "operator"),
  digit("1"),
  digit("2"),
  digit("3"),
  k("+", "+", "Add", "operator"),
  k("0", "0", "0", "digit", true),
  k(".", ".", "Decimal point", "digit"),
  k("=", "=", "Equals", "equals"),
];
