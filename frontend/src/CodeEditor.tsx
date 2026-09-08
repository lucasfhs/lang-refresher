import { useEffect, useRef } from "react";
import { basicSetup } from "codemirror";
import { indentWithTab } from "@codemirror/commands";
import { python } from "@codemirror/lang-python";
import { EditorState, Prec } from "@codemirror/state";
import { EditorView, keymap } from "@codemirror/view";
import { oneDark } from "@codemirror/theme-one-dark";

interface Props {
  value: string;
  onChange: (value: string) => void;
  onRun: () => void;
  onSave: () => void;
}

export function CodeEditor({ value, onChange, onRun, onSave }: Props) {
  const hostRef = useRef<HTMLDivElement>(null);
  const viewRef = useRef<EditorView | null>(null);
  const callbacks = useRef({ onChange, onRun, onSave });
  callbacks.current = { onChange, onRun, onSave };

  useEffect(() => {
    if (!hostRef.current) return;
    const state = EditorState.create({
      doc: value,
      extensions: [
        basicSetup,
        python(),
        oneDark,
        EditorState.tabSize.of(4),
        EditorView.lineWrapping,
        EditorView.theme({
          "&": { height: "100%", fontSize: "14px", backgroundColor: "#0c111a" },
          ".cm-scroller": { overflow: "auto", fontFamily: "'JetBrains Mono', 'Cascadia Code', Consolas, monospace" },
          ".cm-content": { padding: "14px 0" },
          ".cm-gutters": { backgroundColor: "#0c111a", borderRight: "1px solid #202a3b" },
          ".cm-activeLine": { backgroundColor: "#131c2b" },
          ".cm-activeLineGutter": { backgroundColor: "#172133" },
        }),
        Prec.high(keymap.of([
          indentWithTab,
          { key: "Ctrl-Enter", mac: "Cmd-Enter", preventDefault: true, run: () => { callbacks.current.onRun(); return true; } },
          { key: "Ctrl-s", mac: "Cmd-s", preventDefault: true, run: () => { callbacks.current.onSave(); return true; } },
        ])),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) callbacks.current.onChange(update.state.doc.toString());
        }),
      ],
    });
    const view = new EditorView({ state, parent: hostRef.current });
    viewRef.current = view;
    return () => view.destroy();
  }, []);

  useEffect(() => {
    const view = viewRef.current;
    if (!view) return;
    const current = view.state.doc.toString();
    if (current !== value) {
      view.dispatch({ changes: { from: 0, to: current.length, insert: value } });
    }
  }, [value]);

  return <div className="editor-host" ref={hostRef} aria-label="Editor Python" />;
}
