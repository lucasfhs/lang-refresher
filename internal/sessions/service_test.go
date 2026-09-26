package sessions

import (
	"path/filepath"
	"testing"

	"language-refresher/internal/exercises"
	"language-refresher/internal/progress"
)

func TestSessionFlow(t *testing.T) {
	track := exercises.Track{ID: "track", Exercises: []exercises.Exercise{{ID: "one", StarterCode: "one"}, {ID: "two", StarterCode: "two"}}}
	service, err := New(track, progress.NewStore(filepath.Join(t.TempDir(), "progress.json")))
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SaveDraft("one", "edited"); err != nil {
		t.Fatal(err)
	}
	if err := service.Complete(); err != nil {
		t.Fatal(err)
	}
	_, state := service.Snapshot()
	if state.CurrentIndex != 1 || state.Statuses["one"] != progress.StatusCompleted || state.Drafts["one"] != "edited" {
		t.Fatalf("estado inesperado: %#v", state)
	}
	if err := service.Skip(); err != nil {
		t.Fatal(err)
	}
	_, state = service.Snapshot()
	if state.CurrentIndex != 1 || state.Statuses["two"] != progress.StatusSkipped {
		t.Fatalf("estado final inesperado: %#v", state)
	}
}

func TestSwitchTrackRestoresIndependentState(t *testing.T) {
	store := progress.NewStore(filepath.Join(t.TempDir(), "progress.json"))
	python := exercises.Track{ID: "python", Exercises: []exercises.Exercise{{ID: "py-one"}, {ID: "py-two"}}}
	numpy := exercises.Track{ID: "numpy", Exercises: []exercises.Exercise{{ID: "np-one"}, {ID: "np-two"}}}
	service, err := New(python, store)
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Navigate(1); err != nil {
		t.Fatal(err)
	}
	if err := service.SwitchTrack(numpy); err != nil {
		t.Fatal(err)
	}
	track, state := service.Snapshot()
	if track.ID != "numpy" || state.CurrentIndex != 0 {
		t.Fatalf("troca para NumPy inválida: %s %#v", track.ID, state)
	}
	if err := service.SaveDraft("np-one", "numpy"); err != nil {
		t.Fatal(err)
	}
	if err := service.SwitchTrack(python); err != nil {
		t.Fatal(err)
	}
	_, state = service.Snapshot()
	if state.CurrentIndex != 1 {
		t.Fatalf("posição Python não foi restaurada: %#v", state)
	}
	if err := service.SwitchTrack(numpy); err != nil {
		t.Fatal(err)
	}
	_, state = service.Snapshot()
	if state.Drafts["np-one"] != "numpy" {
		t.Fatalf("rascunho NumPy não foi restaurado: %#v", state)
	}
}
