export interface Exercise {
  id: string;
  category: string;
  title: string;
  topic: string;
  objective: string;
  description: string;
  requirements: string[];
  example: string;
  exampleOutput: string;
  hint: string;
  starterCode: string;
  estimatedMinutes: number;
  difficulty: number;
  hasValidator: boolean;
}

export interface SessionState {
  trackId: string;
  trackTitle: string;
  trackDescription: string;
  language: string;
  currentIndex: number;
  total: number;
  completed: number;
  skipped: number;
  elapsedSeconds: number;
  exercise: Exercise;
  draft: string;
  status: "" | "completed" | "skipped";
  statuses: Record<string, string>;
  locale: "pt-BR" | "en";
}

export interface ExecutionResult {
  stdout: string;
  stderr: string;
  exitCode: number;
  durationMs: number;
  timedOut: boolean;
  cancelled: boolean;
  passed: boolean;
}
