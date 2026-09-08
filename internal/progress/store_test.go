package progress

import (
	"path/filepath"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "nested", "progress.json"))
	state, err := store.Load("python-core")
	if err != nil {
		t.Fatal(err)
	}
	state.CurrentIndex = 3
	state.Statuses["exercise"] = StatusCompleted
	state.Drafts["exercise"] = "print('ok')"
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.Load("python-core")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.CurrentIndex != 3 || loaded.Statuses["exercise"] != StatusCompleted || loaded.Drafts["exercise"] != "print('ok')" {
		t.Fatalf("round trip inválido: %#v", loaded)
	}
}

func TestStoreStartsFreshForAnotherTrack(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "progress.json"))
	state, _ := store.Load("one")
	state.CurrentIndex = 2
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	other, err := store.Load("two")
	if err != nil {
		t.Fatal(err)
	}
	if other.TrackID != "two" || other.CurrentIndex != 0 {
		t.Fatalf("esperava novo estado, recebeu %#v", other)
	}
}
