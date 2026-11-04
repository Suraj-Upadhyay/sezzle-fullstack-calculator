import { formatNumber } from "../calculator/format";

export interface HistoryEntry {
  id: number;
  expression: string;
  result: number;
}

interface Props {
  entries: readonly HistoryEntry[];
  onSelect: (entry: HistoryEntry) => void;
  onClear: () => void;
}

export function History({ entries, onSelect, onClear }: Props) {
  if (entries.length === 0) return null;

  return (
    <section aria-labelledby="history-heading" className="flex flex-col gap-2">
      <div className="flex items-center justify-between">
        <h2
          id="history-heading"
          className="text-sm font-medium text-slate-600 dark:text-slate-300"
        >
          History
        </h2>
        <button
          type="button"
          onClick={onClear}
          className="text-sm text-slate-500 hover:text-slate-800 dark:hover:text-slate-200"
        >
          Clear history
        </button>
      </div>
      <ul className="divide-y divide-slate-200 rounded-xl border border-slate-200 dark:divide-slate-700 dark:border-slate-700">
        {entries.map((entry) => (
          <li key={entry.id}>
            <button
              type="button"
              onClick={() => {
                onSelect(entry);
              }}
              title="Edit this expression"
              className="flex w-full justify-between gap-4 px-3 py-2 text-left font-mono text-sm hover:bg-slate-100 dark:hover:bg-slate-800"
            >
              <span className="truncate text-slate-500 dark:text-slate-400">
                {entry.expression}
              </span>
              <span className="shrink-0 font-semibold text-slate-900 dark:text-slate-100">
                = {formatNumber(entry.result)}
              </span>
            </button>
          </li>
        ))}
      </ul>
    </section>
  );
}
