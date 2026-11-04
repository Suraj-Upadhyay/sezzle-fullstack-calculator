import type { Key } from "../calculator/expression";
import { LAYOUT, type Variant } from "../calculator/keypad";

const VARIANTS: Record<Variant, string> = {
  digit:
    "bg-white text-slate-900 hover:bg-slate-100 data-[pressed=true]:bg-slate-200 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700 dark:data-[pressed=true]:bg-slate-600",
  function:
    "bg-slate-200 text-slate-700 hover:bg-slate-300 data-[pressed=true]:bg-slate-400 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600 dark:data-[pressed=true]:bg-slate-500",
  operator:
    "bg-indigo-100 text-indigo-700 hover:bg-indigo-200 data-[pressed=true]:bg-indigo-300 dark:bg-indigo-950 dark:text-indigo-200 dark:hover:bg-indigo-900 dark:data-[pressed=true]:bg-indigo-800",
  equals:
    "bg-indigo-600 text-white hover:bg-indigo-700 data-[pressed=true]:bg-indigo-800",
};

interface Props {
  pressed: Key | null;
  disabled: (key: Key) => boolean;
  onPress: (key: Key) => void;
}

export function Keypad({ pressed, disabled, onPress }: Props) {
  return (
    <div role="group" aria-label="Keypad" className="grid grid-cols-4 gap-2">
      {LAYOUT.map((spec) => (
        <button
          key={spec.key}
          type="button"
          aria-label={spec.name}
          data-pressed={pressed === spec.key}
          disabled={disabled(spec.key)}
          onPointerDown={(e) => {
            e.preventDefault();
          }}
          onClick={() => {
            onPress(spec.key);
          }}
          className={`h-14 rounded-2xl text-xl font-semibold shadow-sm transition duration-75 select-none focus-visible:ring-2 focus-visible:ring-indigo-400 focus-visible:outline-none active:scale-95 disabled:opacity-50 data-[pressed=true]:scale-95 sm:h-16 ${spec.wide ? "col-span-2" : ""} ${VARIANTS[spec.variant]}`}
        >
          {spec.label}
        </button>
      ))}
    </div>
  );
}
