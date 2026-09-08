import { Check, Settings, X } from "lucide-react";
import { getCopy } from "./i18n";

interface Props {
  locale: "pt-BR" | "en";
  onChange: (locale: "pt-BR" | "en") => void;
  onClose: () => void;
}

export function SettingsModal({ locale, onChange, onClose }: Props) {
  const t = getCopy(locale);
  return (
    <div className="modal-backdrop" onMouseDown={onClose}>
      <section className="settings-modal" role="dialog" aria-modal="true" aria-labelledby="settings-title" onMouseDown={(event) => event.stopPropagation()}>
        <header>
          <div className="modal-title"><span className="modal-icon"><Settings size={18} /></span><div><h2 id="settings-title">{t.settings}</h2><p>{t.settingsDescription}</p></div></div>
          <button className="icon-button" onClick={onClose} title={t.close}><X size={18} /></button>
        </header>
        <div className="language-options">
          <button className={locale === "pt-BR" ? "language-option selected" : "language-option"} onClick={() => onChange("pt-BR")}>
            <span><strong>{t.portuguese}</strong><small>Português (Brasil)</small></span>{locale === "pt-BR" && <Check size={18} />}
          </button>
          <button className={locale === "en" ? "language-option selected" : "language-option"} onClick={() => onChange("en")}>
            <span><strong>{t.english}</strong><small>English</small></span>{locale === "en" && <Check size={18} />}
          </button>
        </div>
        <p className="detection-note">{t.detected}</p>
      </section>
    </div>
  );
}
