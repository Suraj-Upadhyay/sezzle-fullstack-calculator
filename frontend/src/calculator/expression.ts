export type Digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";
export type Operator = "+" | "−" | "×" | "÷" | "^";
export type EditKey =
  Digit | Operator | "." | "%" | "√" | "(" | ")" | "backspace";
export type Key = EditKey | "clear" | "=";

export const MAX_LENGTH = 200;

const TRAILING_OPERATORS = /[+−×÷^]+$/;
const TRAILING_NUMBER = /[\d.]+(e[+-]?\d+)?$/;

export function endsWithOperand(expression: string): boolean {
  return /[\d.)%]$/.test(expression);
}

export function openBrackets(expression: string): number {
  let depth = 0;
  for (const char of expression) {
    if (char === "(") depth++;
    if (char === ")") depth--;
  }
  return depth;
}

export function isComplete(expression: string): boolean {
  return endsWithOperand(expression);
}

export function closeBrackets(expression: string): string {
  return expression + ")".repeat(Math.max(0, openBrackets(expression)));
}

// applyKey only ever produces input the API can parse once its brackets
// are closed: operators replace each other, decimal points are not
// repeated and a value next to a bracket gets an explicit ×.
export function applyKey(expression: string, key: EditKey): string {
  if (key === "backspace") return expression.slice(0, -1);
  if (expression.length >= MAX_LENGTH) return expression;

  const afterValue = /[)%]$/.test(expression);
  const number = TRAILING_NUMBER.exec(expression)?.[0] ?? "";

  switch (key) {
    case ".":
      if (number.includes(".") || number.includes("e")) return expression;
      if (afterValue) return `${expression}×0.`;
      return number === "" ? `${expression}0.` : `${expression}.`;
    case "(":
    case "√":
      return endsWithOperand(expression)
        ? `${expression}×${key}`
        : expression + key;
    case ")":
      return endsWithOperand(expression) && openBrackets(expression) > 0
        ? expression + key
        : expression;
    case "%":
      return endsWithOperand(expression) && !expression.endsWith("%")
        ? expression + key
        : expression;
    case "+":
    case "−":
    case "×":
    case "÷":
    case "^":
      return applyOperator(expression, key);
    default:
      if (afterValue) return `${expression}×${key}`;
      if (number === "0") return expression.slice(0, -1) + key;
      return expression + key;
  }
}

function applyOperator(expression: string, operator: Operator): string {
  if (endsWithOperand(expression)) return expression + operator;
  if (operator === "−" && !expression.endsWith("−")) {
    return expression + operator;
  }
  const base = expression.replace(TRAILING_OPERATORS, "");
  return endsWithOperand(base) ? base + operator : expression;
}

const KEY_ALIASES: Record<string, Key> = {
  ",": ".",
  "-": "−",
  "*": "×",
  x: "×",
  X: "×",
  "/": "÷",
  r: "√",
  R: "√",
  Enter: "=",
  Backspace: "backspace",
  Escape: "clear",
  Delete: "clear",
};

const KEYS = new Set<string>(Array.from("0123456789.+−×÷^%√()="));

export function keyFromKeyboard(key: string): Key | null {
  if (key in KEY_ALIASES) return KEY_ALIASES[key] ?? null;
  return KEYS.has(key) ? (key as Key) : null;
}
