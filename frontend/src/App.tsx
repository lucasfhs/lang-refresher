import { useCallback, useEffect, useRef, useState } from "react";
import { api } from "./api";
import { ExercisePanel } from "./ExercisePanel";
import { Workspace } from "./Workspace";
import { SettingsModal } from "./SettingsModal";
import { getCopy } from "./i18n";
import { Settings } from "lucide-react";
import { ResizeHandle } from "./ResizeHandle";
import type { ExecutionResult, SessionState } from "./types";
import "./styles.css";

export default function App() {
  const [state, setState] = useState<SessionState | null>(null);
  const [code, setCode] = useState("");
  const codeRef = useRef("");
  const [result, setResult] = useState<ExecutionResult | null>(null);
  const [error, setError] = useState("");
  const [running, setRunning] = useState(false);
  const [validating, setValidating] = useState(false);
  const [saved, setSaved] = useState(true);
  const [elapsedSeconds, setElapsedSeconds] = useState(0);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const runningRef = useRef(false);
  const [exerciseWidth, setExerciseWidth] = useState(() => clampExerciseWidth(readStoredNumber("layout.exerciseWidth", 430)));

  const loadState = useCallback((next: SessionState) => {
    setState(next);
    setCode(next.draft);
    codeRef.current = next.draft;
    setResult(null);
    setError("");
    setSaved(true);
    setElapsedSeconds(next.elapsedSeconds);
  }, []);

  useEffect(() => {
    api.initialize(navigator.language || "en").then(loadState).catch((err) => setError(readError(err)));
  }, [loadState]);

  useEffect(() => {
    const timer = window.setInterval(() => setElapsedSeconds((value) => value + 1), 1000);
    return () => window.clearInterval(timer);
  }, []);

  useEffect(() => {
    const handleResize = () => setExerciseWidth((width) => clampExerciseWidth(width));
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    storeNumber("layout.exerciseWidth", exerciseWidth);
  }, [exerciseWidth]);

  useEffect(() => {
    if (!state || saved) return;
    const timer = window.setTimeout(() => {
      api.saveDraft(state.exercise.id, codeRef.current).then(() => setSaved(true)).catch(() => undefined);
    }, 900);
    return () => window.clearTimeout(timer);
  }, [code, saved, state]);

  const onCodeChange = (value: string) => {
    codeRef.current = value;
    setCode(value);
    setSaved(false);
  };

  const save = useCallback(async () => {
    if (!state) return;
    try {
      await api.saveDraft(state.exercise.id, codeRef.current);
      setSaved(true);
    } catch (err) {
      setError(readError(err));
    }
  }, [state]);

  const execute = useCallback(async (validate: boolean) => {
    if (runningRef.current) return;
    runningRef.current = true;
    setRunning(true);
    setValidating(validate);
    setResult(null);
    setError("");
    try {
      const nextResult = validate ? await api.validate(codeRef.current) : await api.run(codeRef.current);
      setResult(nextResult);
      setSaved(true);
    } catch (err) {
      setError(readError(err));
    } finally {
      runningRef.current = false;
      setRunning(false);
      setValidating(false);
    }
  }, []);

  useEffect(() => {
    const handleRunShortcut = (event: KeyboardEvent) => {
      if (settingsOpen || event.repeat) return;
      if ((event.ctrlKey || event.metaKey) && event.key === "Enter") {
        event.preventDefault();
        event.stopPropagation();
        void execute(false);
      }
    };
    window.addEventListener("keydown", handleRunShortcut, { capture: true });
    return () => window.removeEventListener("keydown", handleRunShortcut, { capture: true });
  }, [execute, settingsOpen]);

  const navigate = async (delta: number) => {
    if (!state) return;
    try {
      await api.saveDraft(state.exercise.id, codeRef.current);
      loadState(await api.navigate(delta));
    } catch (err) { setError(readError(err)); }
  };

  const advance = async (kind: "complete" | "skip") => {
    try {
      loadState(kind === "complete" ? await api.complete(codeRef.current) : await api.skip(codeRef.current));
    } catch (err) { setError(readError(err)); }
  };

  const restart = async () => {
    if (!state || !window.confirm(getCopy(state.locale).restartConfirm)) return;
    try { loadState(await api.restart()); } catch (err) { setError(readError(err)); }
  };

  const changeLanguage = async (locale: "pt-BR" | "en") => {
    try {
      loadState(await api.setLanguage(locale));
      setSettingsOpen(false);
    } catch (err) { setError(readError(err)); }
  };

  const changeTrack = async (trackId: string) => {
    if (!state || trackId === state.trackId) return;
    try {
      await api.saveDraft(state.exercise.id, codeRef.current);
      loadState(await api.setTrack(trackId));
      setSettingsOpen(false);
    } catch (err) { setError(readError(err)); }
  };

  if (!state) {
    const initialLocale = (navigator.language || "en").toLowerCase().startsWith("pt") ? "pt-BR" : "en";
    return <div className="loading-screen"><div className="brand-mark large"><Settings size={22} /></div><p>{error || getCopy(initialLocale).loading}</p></div>;
  }

  return (
    <div className="app-shell" style={{ gridTemplateColumns: `${exerciseWidth}px 7px minmax(0, 1fr)` }}>
      <ExercisePanel state={{ ...state, elapsedSeconds }} busy={running} onPrevious={() => navigate(-1)} onNext={() => navigate(1)} onComplete={() => advance("complete")} onSkip={() => advance("skip")} onRestart={restart} onSettings={() => setSettingsOpen(true)} />
      <ResizeHandle orientation="vertical" ariaLabel={getCopy(state.locale).resizeExercise} onDrag={(delta) => setExerciseWidth((width) => clampExerciseWidth(width + delta))} onReset={() => setExerciseWidth(clampExerciseWidth(430))} />
      <Workspace state={state} code={code} result={result} error={error} running={running} validating={validating} saved={saved} onCodeChange={onCodeChange} onRun={() => execute(false)} onValidate={() => execute(true)} onCancel={() => api.cancel()} onSave={save} />
      {settingsOpen && <SettingsModal locale={state.locale} currentTrackId={state.trackId} tracks={state.availableTracks} onChange={changeLanguage} onTrackChange={changeTrack} onClose={() => setSettingsOpen(false)} />}
    </div>
  );
}

function clampExerciseWidth(width: number) {
  const viewport = window.innerWidth;
  const minimum = viewport < 800 ? 260 : 300;
  const workspaceMinimum = viewport < 800 ? 330 : 420;
  const maximum = Math.max(minimum, viewport - workspaceMinimum - 7);
  return Math.round(Math.min(maximum, Math.max(minimum, width)));
}

function readStoredNumber(key: string, fallback: number) {
  try {
    const value = Number(window.localStorage.getItem(key));
    return Number.isFinite(value) && value > 0 ? value : fallback;
  } catch {
    return fallback;
  }
}

function storeNumber(key: string, value: number) {
  try {
    window.localStorage.setItem(key, String(value));
  } catch {
    // O layout continua funcional quando o armazenamento do WebView está indisponível.
  }
}

function readError(error: unknown) {
  return error instanceof Error ? error.message : String(error);
}
