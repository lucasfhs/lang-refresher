import { Check, Clock3, Settings, X } from "lucide-react";
import { getCopy } from "./i18n";
import type { TrackSummary } from "./types";

interface Props {
  locale: "pt-BR" | "en";
  currentTrackId: string;
  tracks: TrackSummary[];
  onChange: (locale: "pt-BR" | "en") => void;
  onTrackChange: (trackId: string) => void;
  onClose: () => void;
}

export function SettingsModal({ locale, currentTrackId, tracks, onChange, onTrackChange, onClose }: Props) {
  const t = getCopy(locale);
  return (
    <div className="modal-backdrop" onMouseDown={onClose}>
      <section className="settings-modal" role="dialog" aria-modal="true" aria-labelledby="settings-title" onMouseDown={(event) => event.stopPropagation()}>
        <header>
          <div className="modal-title"><span className="modal-icon"><Settings size={18} /></span><div><h2 id="settings-title">{t.settings}</h2><p>{t.settingsDescription}</p></div></div>
          <button className="icon-button" onClick={onClose} title={t.close}><X size={18} /></button>
        </header>
        <div className="settings-body">
          <section className="settings-section">
            <h3>{t.practice}</h3>
            <div className="practice-options">
              {tracks.map((track) => (
                <button key={track.id} className={track.id === currentTrackId ? "practice-option selected" : "practice-option"} onClick={() => onTrackChange(track.id)}>
                  <span className="practice-copy"><strong>{track.title}</strong><small>{track.description}</small><em><Clock3 size={12} /> ~{track.estimatedMinutes} min · {track.exerciseCount} {t.exercises}</em></span>
                  {track.id === currentTrackId && <Check size={18} />}
                </button>
              ))}
            </div>
          </section>
          <section className="settings-section">
            <h3>{t.interfaceLanguage}</h3>
            <div className="language-options">
              <button className={locale === "pt-BR" ? "language-option selected" : "language-option"} onClick={() => onChange("pt-BR")}>
                <span><strong>{t.portuguese}</strong><small>Português (Brasil)</small></span>{locale === "pt-BR" && <Check size={18} />}
              </button>
              <button className={locale === "en" ? "language-option selected" : "language-option"} onClick={() => onChange("en")}>
                <span><strong>{t.english}</strong><small>English</small></span>{locale === "en" && <Check size={18} />}
              </button>
            </div>
          </section>
        </div>
        <p className="detection-note">{t.detected}</p>
      </section>
    </div>
  );
}
