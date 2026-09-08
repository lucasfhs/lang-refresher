import { Check, ChevronLeft, ChevronRight, Clock3, RotateCcw, Settings, SkipForward } from "lucide-react";
import { getCopy } from "./i18n";
import type { SessionState } from "./types";

interface Props {
  state: SessionState;
  busy: boolean;
  onPrevious: () => void;
  onNext: () => void;
  onComplete: () => void;
  onSkip: () => void;
  onRestart: () => void;
  onSettings: () => void;
}

export function ExercisePanel({ state, busy, onPrevious, onNext, onComplete, onSkip, onRestart, onSettings }: Props) {
  const exercise = state.exercise;
  const t = getCopy(state.locale);
  return (
    <aside className="exercise-panel">
      <header className="track-header">
        <div className="brand-row">
          <button className="settings-button" onClick={onSettings} title={t.settings}><Settings size={18} /></button>
          <div>
            <div className="eyebrow">LANGUAGE REFRESHER</div>
            <h1>{state.trackTitle}</h1>
          </div>
        </div>
        <button className="icon-button" onClick={onRestart} title={t.restart} disabled={busy}>
          <RotateCcw size={17} />
        </button>
      </header>

      <section className="progress-card">
        <div className="progress-copy">
          <strong>{state.completed} / {state.total} {t.exercises}</strong>
          <span>{Math.round((state.completed / state.total) * 100)}%</span>
        </div>
        <div className="progress-track"><div style={{ width: `${(state.completed / state.total) * 100}%` }} /></div>
        <div className="progress-meta">
          <span><Clock3 size={14} /> {formatDuration(state.elapsedSeconds)}</span>
          {state.skipped > 0 && <span>{state.skipped} {state.skipped === 1 ? t.skipped : t.skippedPlural}</span>}
        </div>
      </section>

      <div className="exercise-scroll">
        <div className="task-kicker">
          <span>{t.task} {state.currentIndex + 1} {t.of} {state.total}</span>
          <span className="category-pill">{exercise.category}</span>
        </div>
        <h2>{exercise.title}</h2>
        <div className="task-meta">
          <span>{exercise.topic}</span><span>•</span><span>~{exercise.estimatedMinutes} min</span>
        </div>

        <ContentSection title={t.objective}><p>{exercise.objective}</p></ContentSection>
        <ContentSection title={t.assignment}><p className="statement">{exercise.description}</p></ContentSection>
        <ContentSection title={t.requirements}>
          <ul>{exercise.requirements.map((item) => <li key={item}>{item}</li>)}</ul>
        </ContentSection>
        {exercise.example && <ContentSection title={t.example}><pre className="example">{exercise.example}</pre></ContentSection>}
        <ContentSection title={t.exampleOutput}><pre className="example output-example">{exercise.exampleOutput}</pre></ContentSection>
        {exercise.hint && <details className="hint"><summary>{t.showHint}</summary><p>{exercise.hint}</p></details>}
      </div>

      <footer className="exercise-actions">
        <div className="navigation-row">
          <button className="secondary" onClick={onPrevious} disabled={busy || state.currentIndex === 0}><ChevronLeft size={17} /> {t.previous}</button>
          <button className="secondary" onClick={onNext} disabled={busy || state.currentIndex === state.total - 1}>{t.next} <ChevronRight size={17} /></button>
        </div>
        <div className="completion-row">
          <button className="ghost" onClick={onSkip} disabled={busy}><SkipForward size={16} /> {t.skip}</button>
          <button className="complete" onClick={onComplete} disabled={busy}><Check size={18} /> {t.complete}</button>
        </div>
      </footer>
    </aside>
  );
}

function ContentSection({ title, children }: { title: string; children: React.ReactNode }) {
  return <section className="content-section"><h3>{title}</h3>{children}</section>;
}

function formatDuration(totalSeconds: number) {
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  if (hours > 0) return `${hours}h ${minutes}min`;
  return `${minutes} min`;
}
