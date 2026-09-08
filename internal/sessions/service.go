package sessions

import (
	"fmt"
	"sync"
	"time"

	"language-refresher/internal/exercises"
	"language-refresher/internal/progress"
)

type Service struct {
	mu    sync.Mutex
	track exercises.Track
	store *progress.Store
	state progress.State
}

func New(track exercises.Track, store *progress.Store) (*Service, error) {
	state, err := store.Load(track.ID)
	if err != nil {
		return nil, err
	}
	if state.CurrentIndex < 0 || state.CurrentIndex >= len(track.Exercises) {
		state.CurrentIndex = 0
	}
	return &Service{track: track, store: store, state: state}, nil
}

func (s *Service) Snapshot() (exercises.Track, progress.State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.track, cloneState(s.state)
}

func (s *Service) CurrentExercise() exercises.Exercise {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.track.Exercises[s.state.CurrentIndex]
}

func (s *Service) SaveDraft(id, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.track.Exercises[s.state.CurrentIndex].ID != id {
		return fmt.Errorf("o exercício %s não está ativo", id)
	}
	s.state.Drafts[id] = code
	return s.store.Save(s.state)
}

func (s *Service) Navigate(delta int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := s.state.CurrentIndex + delta
	if next < 0 {
		next = 0
	}
	if next >= len(s.track.Exercises) {
		next = len(s.track.Exercises) - 1
	}
	s.state.CurrentIndex = next
	return s.store.Save(s.state)
}

func (s *Service) Complete() error { return s.setStatus(progress.StatusCompleted) }
func (s *Service) Skip() error     { return s.setStatus(progress.StatusSkipped) }

func (s *Service) setStatus(status progress.ExerciseStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.track.Exercises[s.state.CurrentIndex].ID
	s.state.Statuses[id] = status
	if s.state.CurrentIndex < len(s.track.Exercises)-1 {
		s.state.CurrentIndex++
	}
	return s.store.Save(s.state)
}

func (s *Service) Restart() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.state = progress.State{TrackID: s.track.ID, Statuses: map[string]progress.ExerciseStatus{}, Drafts: map[string]string{}, StartedAt: now, UpdatedAt: now}
	return s.store.Save(s.state)
}

func cloneState(state progress.State) progress.State {
	clone := state
	clone.Statuses = make(map[string]progress.ExerciseStatus, len(state.Statuses))
	for key, value := range state.Statuses {
		clone.Statuses[key] = value
	}
	clone.Drafts = make(map[string]string, len(state.Drafts))
	for key, value := range state.Drafts {
		clone.Drafts[key] = value
	}
	return clone
}
