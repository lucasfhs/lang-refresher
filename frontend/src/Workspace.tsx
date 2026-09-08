import { CheckCircle2, CircleStop, FlaskConical, Play, Save, TerminalSquare } from "lucide-react";
import { CodeEditor } from "./CodeEditor";
import { getCopy } from "./i18n";
import type { ExecutionResult, SessionState } from "./types";

interface Props {
  state: SessionState;
  code: string;
  result: ExecutionResult | null;
  error: string;
  running: boolean;
  validating: boolean;
  saved: boolean;
  onCodeChange: (code: string) => void;
  onRun: () => void;
  onValidate: () => void;
  onCancel: () => void;
  onSave: () => void;
}

export function Workspace(props: Props) {
  const { state, code, result, error, running, validating, saved } = props;
  const t = getCopy(state.locale);
  return (
    <main className="workspace">
      <section className="editor-section">
        <header className="pane-header">
          <div className="file-tab"><span className="python-icon">Py</span> main.py <span className="dirty-indicator">{saved ? "" : "●"}</span></div>
          <div className="editor-actions">
            <button className="toolbar-button" onClick={props.onSave} title={`${t.save} (Ctrl+S)`}><Save size={16} /> {saved ? t.saved : t.save}</button>
            {state.exercise.hasValidator && (
              <button className="validate-button" onClick={props.onValidate} disabled={running}>
                <FlaskConical size={17} /> {validating ? t.validating : t.validate}
              </button>
            )}
            {running ? (
              <button className="stop-button" onClick={props.onCancel}><CircleStop size={18} /> {t.cancel}</button>
            ) : (
              <button className="run-button" onClick={props.onRun}><Play size={17} fill="currentColor" /> {t.run} <kbd>Ctrl ↵</kbd></button>
            )}
          </div>
        </header>
        <CodeEditor value={code} onChange={props.onCodeChange} onRun={props.onRun} onSave={props.onSave} />
      </section>

      <section className="terminal-section">
        <header className="pane-header terminal-header">
          <div><TerminalSquare size={17} /> {t.output}</div>
          {result && <span className={result.exitCode === 0 ? "exit-success" : "exit-error"}>exit {result.exitCode} · {result.durationMs} ms</span>}
        </header>
        <div className="terminal-output" role="log" aria-live="polite">
          {!result && !error && !running && <span className="terminal-empty">{t.outputEmpty}</span>}
          {running && <span className="terminal-running"><span className="spinner" /> {t.running}</span>}
          {error && <pre className="stderr">{error}</pre>}
          {result?.passed && <div className="passed"><CheckCircle2 size={17} /> {t.passed}</div>}
          {result?.timedOut && <pre className="stderr">{t.timedOut}</pre>}
          {result?.cancelled && <pre className="stderr">{t.cancelled}</pre>}
          {result?.stdout && <pre className="stdout">{result.stdout}</pre>}
          {result?.stderr && <pre className="stderr">{result.stderr}</pre>}
        </div>
      </section>
    </main>
  );
}
