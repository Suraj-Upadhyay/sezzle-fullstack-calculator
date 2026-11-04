import { errorMessage, formatNumber } from "../calculator/format";
import type { EvaluationState } from "../calculator/useEvaluation";

interface Props {
  input: string;
  state: EvaluationState;
  hint: string | null;
}

export function Display({ input, state, hint }: Props) {
  const success = state.status === "success" ? state.response : null;
  const main = success ? formatNumber(success.result) : input || "0";
  const caption = success ? `${success.expression} =` : "";

  return (
    <section
      aria-label="Display"
      className="flex min-h-36 flex-col justify-end gap-1 rounded-3xl bg-slate-900 px-5 py-4 text-right text-white dark:bg-black"
    >
      <p
        aria-label="Expression"
        className="min-h-5 font-mono text-sm break-all text-slate-400"
      >
        {caption}
      </p>
      <output
        aria-label={success ? "Result" : "Input"}
        aria-live={success ? "polite" : "off"}
        className={`font-mono font-semibold break-all ${main.length > 14 ? "text-2xl" : "text-4xl"}`}
      >
        {main}
      </output>
      <Status state={state} hint={hint} />
    </section>
  );
}

function Status({ state, hint }: Pick<Props, "state" | "hint">) {
  if (hint) {
    return (
      <p role="alert" className="text-sm text-amber-300">
        {hint}
      </p>
    );
  }
  if (state.status === "error") {
    return (
      <p role="alert" className="text-sm text-red-300">
        {errorMessage(state.error.message, state.error.detail.reason)}
      </p>
    );
  }
  return (
    <p className="min-h-5 text-sm text-slate-500">
      {state.status === "loading" ? "Calculating…" : ""}
    </p>
  );
}
