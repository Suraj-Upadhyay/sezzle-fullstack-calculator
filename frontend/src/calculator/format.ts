const standard = new Intl.NumberFormat("en-US", {
  maximumSignificantDigits: 15,
});

const scientific = new Intl.NumberFormat("en-US", {
  notation: "scientific",
  maximumSignificantDigits: 10,
});

// Rounds to 15 significant digits, the most a float64 holds exactly, so
// 0.1 + 0.2 displays as 0.3.
export function formatNumber(value: number): string {
  if (value === 0) return "0";
  const magnitude = Math.abs(value);
  if (magnitude >= 1e15 || magnitude < 1e-9) {
    return scientific.format(value);
  }
  return standard.format(value);
}

export function toExpression(value: number): string {
  return String(Number(value.toPrecision(15))).replace(/^-/, "−");
}

export function toSentence(message: string): string {
  const text = message.trim();
  if (text === "") return text;
  const capitalised = text.charAt(0).toUpperCase() + text.slice(1);
  return /[.!?]$/.test(capitalised) ? capitalised : `${capitalised}.`;
}

export function errorMessage(message: string, reason?: string): string {
  return toSentence(reason ? `${message}: ${reason}` : message);
}
