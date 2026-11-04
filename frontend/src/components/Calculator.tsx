import { useCallback, useEffect, useRef, useState } from "react";

import {
  applyKey,
  closeBrackets,
  isComplete,
  keyFromKeyboard,
  type EditKey,
  type Key,
} from "../calculator/expression";
import { toExpression } from "../calculator/format";
import { useEvaluation } from "../calculator/useEvaluation";
import { Display } from "./Display";
import { History, type HistoryEntry } from "./History";
import { Keypad } from "./Keypad";

const MAX_HISTORY = 10;
const PRESS_FLASH_MS = 120;
const STARTS_NEW_EXPRESSION = new Set(Array.from("0123456789.(√") as Key[]);

export function Calculator() {
  const [input, setInput] = useState("");
  const [hint, setHint] = useState<string | null>(null);
  const [history, setHistory] = useState<HistoryEntry[]>([]);
  const [pressed, setPressed] = useState<Key | null>(null);
  const flashTimer = useRef<number | undefined>(undefined);
  const nextId = useRef(0);
  const { state, run, reset } = useEvaluation();

  const flash = useCallback((key: Key) => {
    setPressed(key);
    window.clearTimeout(flashTimer.current);
    flashTimer.current = window.setTimeout(() => {
      setPressed(null);
    }, PRESS_FLASH_MS);
  }, []);

  useEffect(
    () => () => {
      window.clearTimeout(flashTimer.current);
    },
    [],
  );

  const edit = useCallback(
    (key: EditKey) => {
      setHint(null);
      if (state.status === "success") {
        const start = STARTS_NEW_EXPRESSION.has(key)
          ? ""
          : key === "backspace"
            ? state.response.expression
            : toExpression(state.response.result);
        setInput(applyKey(start, key));
      } else {
        setInput((prev) => applyKey(prev, key));
      }
      if (state.status !== "idle") reset();
    },
    [state, reset],
  );

  const submit = useCallback(async () => {
    if (state.status === "success" || state.status === "loading") return;
    if (!isComplete(input)) {
      setHint(input === "" ? "Enter an expression." : "Finish the expression.");
      return;
    }
    setHint(null);

    const expression = closeBrackets(input);
    setInput(expression);
    const response = await run(expression);
    if (!response) return;

    const entry = { id: nextId.current++, ...response };
    setHistory((prev) => [entry, ...prev].slice(0, MAX_HISTORY));
  }, [input, state.status, run]);

  const press = useCallback(
    (key: Key) => {
      flash(key);
      if (key === "=") {
        void submit();
      } else if (key === "clear") {
        setInput("");
        setHint(null);
        reset();
      } else {
        edit(key);
      }
    },
    [flash, submit, reset, edit],
  );

  useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if (event.ctrlKey || event.metaKey || event.altKey) return;
      const key = keyFromKeyboard(event.key);
      if (!key) return;
      event.preventDefault();
      press(key);
    }

    function onPaste(event: ClipboardEvent) {
      const text = event.clipboardData?.getData("text") ?? "";
      const keys = Array.from(text)
        .map(keyFromKeyboard)
        .filter((key): key is EditKey => key !== null && isEditKey(key));
      if (keys.length === 0) return;
      event.preventDefault();
      reset();
      setHint(null);
      setInput((prev) => keys.reduce(applyKey, prev));
    }

    window.addEventListener("keydown", onKeyDown);
    window.addEventListener("paste", onPaste);
    return () => {
      window.removeEventListener("keydown", onKeyDown);
      window.removeEventListener("paste", onPaste);
    };
  }, [press, reset]);

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-4">
        <Display input={input} state={state} hint={hint} />
        <Keypad
          pressed={pressed}
          disabled={(key) => key === "=" && state.status === "loading"}
          onPress={press}
        />
        <p className="text-center text-xs text-slate-500 dark:text-slate-400">
          Type on your keyboard: Enter to calculate, Esc to clear, R for √.
        </p>
      </div>

      <History
        entries={history}
        onSelect={(entry) => {
          reset();
          setHint(null);
          setInput(entry.expression);
        }}
        onClear={() => {
          setHistory([]);
        }}
      />
    </div>
  );
}

function isEditKey(key: Key): key is EditKey {
  return key !== "=" && key !== "clear";
}
