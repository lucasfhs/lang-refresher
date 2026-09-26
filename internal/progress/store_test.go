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

func TestStoreKeepsProgressForMultipleTracks(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "progress.json"))
	python, _ := store.Load("python")
	python.CurrentIndex = 7
	python.Drafts["py-1"] = "python draft"
	if err := store.Save(python); err != nil {
		t.Fatal(err)
	}

	numpy, _ := store.Load("numpy")
	numpy.CurrentIndex = 3
	numpy.Drafts["np-1"] = "numpy draft"
	if err := store.Save(numpy); err != nil {
		t.Fatal(err)
	}

	loadedPython, err := store.Load("python")
	if err != nil {
		t.Fatal(err)
	}
	loadedNumPy, err := store.Load("numpy")
	if err != nil {
		t.Fatal(err)
	}
	if loadedPython.CurrentIndex != 7 || loadedPython.Drafts["py-1"] != "python draft" {
		t.Fatalf("progresso Python foi perdido: %#v", loadedPython)
	}
	if loadedNumPy.CurrentIndex != 3 || loadedNumPy.Drafts["np-1"] != "numpy draft" {
		t.Fatalf("progresso NumPy foi perdido: %#v", loadedNumPy)
	}
}
