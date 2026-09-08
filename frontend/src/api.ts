import type { ExecutionResult, SessionState } from "./types";

function appMethod<T>(name: string, ...args: unknown[]): Promise<T> {
  const app = window.go?.main?.App;
  const method = app?.[name];
  if (!method) {
    return Promise.reject(new Error("Backend Wails indisponível. Inicie a interface com `wails dev`."));
  }
  return method(...args) as Promise<T>;
}

export const api = {
  initialize: (systemLanguage: string) => appMethod<SessionState>("Initialize", systemLanguage),
  getState: () => appMethod<SessionState>("GetState"),
  setLanguage: (language: string) => appMethod<SessionState>("SetLanguage", language),
  saveDraft: (id: string, code: string) => appMethod<void>("SaveDraft", id, code),
  navigate: (delta: number) => appMethod<SessionState>("Navigate", delta),
  complete: (code: string) => appMethod<SessionState>("CompleteCurrent", code),
  skip: (code: string) => appMethod<SessionState>("SkipCurrent", code),
  restart: () => appMethod<SessionState>("RestartSession"),
  run: (code: string) => appMethod<ExecutionResult>("RunCode", code),
  validate: (code: string) => appMethod<ExecutionResult>("ValidateCode", code),
  cancel: () => appMethod<boolean>("CancelExecution"),
};
