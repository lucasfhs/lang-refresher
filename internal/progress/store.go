package progress

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ExerciseStatus string

const (
	StatusCompleted ExerciseStatus = "completed"
	StatusSkipped   ExerciseStatus = "skipped"
)

type State struct {
	TrackID      string                    `json:"track_id"`
	CurrentIndex int                       `json:"current_index"`
	Statuses     map[string]ExerciseStatus `json:"statuses"`
	Drafts       map[string]string         `json:"drafts"`
	StartedAt    time.Time                 `json:"started_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

type Store struct {
	path string
	mu   sync.Mutex
}

var unsafeTrackID = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func NewStore(path string) *Store { return &Store{path: path} }

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("localizar diretório de configuração: %w", err)
	}
	return filepath.Join(dir, "LanguageRefresher", "progress.json"), nil
}

func (s *Store) Load(trackID string) (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return s.loadTrackFile(trackID)
	}
	if err != nil {
		return State{}, fmt.Errorf("ler progresso: %w", err)
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf("decodificar progresso: %w", err)
	}
	if state.TrackID != trackID {
		return s.loadTrackFile(trackID)
	}
	ensureMaps(&state)
	return state, nil
}

func (s *Store) Save(state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ensureMaps(&state)
	state.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("codificar progresso: %w", err)
	}
	target, err := s.pathForSave(state.TrackID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return fmt.Errorf("criar diretório de progresso: %w", err)
	}
	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("gravar progresso temporário: %w", err)
	}
	if err := os.Rename(tmp, target); err != nil {
		// Windows não substitui sempre um arquivo existente; o fallback direto
		// ainda mantém o arquivo temporário como cópia recuperável em caso de erro.
		if writeErr := os.WriteFile(target, data, 0o600); writeErr != nil {
			return fmt.Errorf("salvar progresso: %w", writeErr)
		}
		_ = os.Remove(tmp)
	}
	return nil
}

func (s *Store) loadTrackFile(trackID string) (State, error) {
	data, err := os.ReadFile(s.trackPath(trackID))
	if errors.Is(err, os.ErrNotExist) {
		return fresh(trackID), nil
	}
	if err != nil {
		return State{}, fmt.Errorf("ler progresso: %w", err)
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf("decodificar progresso: %w", err)
	}
	if state.TrackID != trackID {
		return fresh(trackID), nil
	}
	ensureMaps(&state)
	return state, nil
}

func (s *Store) pathForSave(trackID string) (string, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return s.path, nil
	}
	if err != nil {
		return "", fmt.Errorf("ler progresso: %w", err)
	}
	var existing State
	if json.Unmarshal(data, &existing) == nil && existing.TrackID == trackID {
		return s.path, nil
	}
	return s.trackPath(trackID), nil
}

func (s *Store) trackPath(trackID string) string {
	ext := filepath.Ext(s.path)
	base := strings.TrimSuffix(s.path, ext)
	clean := strings.Trim(unsafeTrackID.ReplaceAllString(trackID, "-"), "-")
	if clean == "" {
		clean = "track"
	}
	return base + "." + clean + ext
}

func fresh(trackID string) State {
	now := time.Now()
	return State{TrackID: trackID, Statuses: map[string]ExerciseStatus{}, Drafts: map[string]string{}, StartedAt: now, UpdatedAt: now}
}

func ensureMaps(state *State) {
	if state.Statuses == nil {
		state.Statuses = map[string]ExerciseStatus{}
	}
	if state.Drafts == nil {
		state.Drafts = map[string]string{}
	}
	if state.StartedAt.IsZero() {
		state.StartedAt = time.Now()
	}
}
