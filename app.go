package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"language-refresher/internal/execution"
	"language-refresher/internal/exercises"
	"language-refresher/internal/languages"
	"language-refresher/internal/progress"
	"language-refresher/internal/sessions"
	"language-refresher/internal/settings"
)

type ExerciseView struct {
	ID               string   `json:"id"`
	Category         string   `json:"category"`
	Title            string   `json:"title"`
	Topic            string   `json:"topic"`
	Objective        string   `json:"objective"`
	Description      string   `json:"description"`
	Requirements     []string `json:"requirements"`
	Example          string   `json:"example"`
	ExampleOutput    string   `json:"exampleOutput"`
	Hint             string   `json:"hint"`
	StarterCode      string   `json:"starterCode"`
	EstimatedMinutes int      `json:"estimatedMinutes"`
	Difficulty       int      `json:"difficulty"`
	HasValidator     bool     `json:"hasValidator"`
}

type StateView struct {
	TrackID          string            `json:"trackId"`
	TrackTitle       string            `json:"trackTitle"`
	TrackDescription string            `json:"trackDescription"`
	Language         string            `json:"language"`
	CurrentIndex     int               `json:"currentIndex"`
	Total            int               `json:"total"`
	Completed        int               `json:"completed"`
	Skipped          int               `json:"skipped"`
	ElapsedSeconds   int64             `json:"elapsedSeconds"`
	Exercise         ExerciseView      `json:"exercise"`
	Draft            string            `json:"draft"`
	Status           string            `json:"status"`
	Statuses         map[string]string `json:"statuses"`
	Locale           string            `json:"locale"`
}

type App struct {
	ctx       context.Context
	sessions  *sessions.Service
	runner    *execution.Runner
	languages *languages.Registry
	exercises *exercises.Repository
	settings  *settings.Store
	localeMu  sync.RWMutex
	locale    string
}

func NewApp(session *sessions.Service, runner *execution.Runner, registry *languages.Registry, repository *exercises.Repository, settingsStore *settings.Store) *App {
	return &App{sessions: session, runner: runner, languages: registry, exercises: repository, settings: settingsStore, locale: "en"}
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx }
func (a *App) shutdown(_ context.Context)  { a.runner.Cancel() }

func (a *App) Initialize(systemLanguage string) (StateView, error) {
	stored, found, err := a.settings.Load()
	if err != nil {
		return StateView{}, err
	}
	locale := resolveLocale(systemLanguage, stored, found)
	a.setLocale(locale)
	return a.stateView(), nil
}

func (a *App) SetLanguage(language string) (StateView, error) {
	locale := normalizeLocale(language)
	if err := a.settings.Save(settings.State{Language: locale}); err != nil {
		return StateView{}, err
	}
	a.setLocale(locale)
	return a.stateView(), nil
}

func (a *App) GetState() (StateView, error) { return a.stateView(), nil }

func (a *App) SaveDraft(exerciseID, code string) error {
	return a.sessions.SaveDraft(exerciseID, code)
}

func (a *App) Navigate(delta int) (StateView, error) {
	if err := a.sessions.Navigate(delta); err != nil {
		return StateView{}, err
	}
	return a.stateView(), nil
}

func (a *App) CompleteCurrent(code string) (StateView, error) {
	current := a.sessions.CurrentExercise()
	if err := a.sessions.SaveDraft(current.ID, code); err != nil {
		return StateView{}, err
	}
	if err := a.sessions.Complete(); err != nil {
		return StateView{}, err
	}
	return a.stateView(), nil
}

func (a *App) SkipCurrent(code string) (StateView, error) {
	current := a.sessions.CurrentExercise()
	if err := a.sessions.SaveDraft(current.ID, code); err != nil {
		return StateView{}, err
	}
	if err := a.sessions.Skip(); err != nil {
		return StateView{}, err
	}
	return a.stateView(), nil
}

func (a *App) RestartSession() (StateView, error) {
	if err := a.sessions.Restart(); err != nil {
		return StateView{}, err
	}
	return a.stateView(), nil
}

func (a *App) RunCode(code string) (execution.Result, error) {
	return a.run(code, false)
}

func (a *App) ValidateCode(code string) (execution.Result, error) {
	return a.run(code, true)
}

func (a *App) CancelExecution() bool { return a.runner.Cancel() }

func (a *App) run(code string, validate bool) (execution.Result, error) {
	exercise := a.sessions.CurrentExercise()
	definition, err := a.languages.Get(exercise.Language)
	if err != nil {
		return execution.Result{}, err
	}
	if err := a.sessions.SaveDraft(exercise.ID, code); err != nil {
		return execution.Result{}, err
	}
	var validator *exercises.Validator
	timeout := 8 * time.Second
	if exercise.TimeoutSeconds > 0 {
		timeout = time.Duration(exercise.TimeoutSeconds) * time.Second
	}
	if validate {
		if exercise.Validator == nil {
			return execution.Result{}, fmt.Errorf("este exercício não possui testes automáticos")
		}
		validator = exercise.Validator
		if validator.TimeoutSeconds > 0 {
			timeout = time.Duration(validator.TimeoutSeconds) * time.Second
		}
	}
	return a.runner.Run(execution.Request{Language: definition, Code: code, SupportFiles: exercise.SupportFiles, Validator: validator, Timeout: timeout})
}

func (a *App) stateView() StateView {
	track, state := a.sessions.Snapshot()
	locale := a.getLocale()
	if localized, ok := a.exercises.LocalizedTrack(track.ID, locale); ok {
		track = localized
	}
	exercise := track.Exercises[state.CurrentIndex]
	completed, skipped := 0, 0
	statuses := make(map[string]string, len(state.Statuses))
	for id, status := range state.Statuses {
		statuses[id] = string(status)
		switch status {
		case progress.StatusCompleted:
			completed++
		case progress.StatusSkipped:
			skipped++
		}
	}
	draft, exists := state.Drafts[exercise.ID]
	if !exists {
		draft = exercise.StarterCode
	}
	return StateView{
		TrackID: track.ID, TrackTitle: track.Title, TrackDescription: track.Description,
		Language: track.Language, CurrentIndex: state.CurrentIndex, Total: len(track.Exercises),
		Completed: completed, Skipped: skipped, ElapsedSeconds: int64(time.Since(state.StartedAt).Seconds()),
		Exercise: toExerciseView(exercise), Draft: draft, Status: statuses[exercise.ID], Statuses: statuses,
		Locale: locale,
	}
}

func (a *App) setLocale(locale string) {
	a.localeMu.Lock()
	a.locale = locale
	a.localeMu.Unlock()
}

func (a *App) getLocale() string {
	a.localeMu.RLock()
	defer a.localeMu.RUnlock()
	return a.locale
}

func normalizeLocale(language string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(language)), "pt") {
		return "pt-BR"
	}
	return "en"
}

func resolveLocale(systemLanguage string, stored settings.State, found bool) string {
	if found {
		return normalizeLocale(stored.Language)
	}
	return normalizeLocale(systemLanguage)
}

func toExerciseView(exercise exercises.Exercise) ExerciseView {
	return ExerciseView{
		ID: exercise.ID, Category: exercise.Category, Title: exercise.Title, Topic: exercise.Topic,
		Objective: exercise.Objective, Description: exercise.Description, Requirements: exercise.Requirements,
		Example: exercise.Example, ExampleOutput: exercise.ExampleOutput, Hint: exercise.Hint, StarterCode: exercise.StarterCode,
		EstimatedMinutes: exercise.EstimatedMinutes, Difficulty: exercise.Difficulty, HasValidator: exercise.Validator != nil,
	}
}
