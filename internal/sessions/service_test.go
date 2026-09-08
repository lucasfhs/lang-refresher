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
