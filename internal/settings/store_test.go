package settings

import (
	"path/filepath"
	"testing"
)

func TestStoreLanguageRoundTrip(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "nested", "settings.json"))
	if _, found, err := store.Load(); err != nil || found {
		t.Fatalf("configuração inicial inesperada: found=%v err=%v", found, err)
	}
	if err := store.Save(State{Language: "en"}); err != nil {
		t.Fatal(err)
	}
	state, found, err := store.Load()
	if err != nil || !found || state.Language != "en" {
		t.Fatalf("round trip inválido: state=%#v found=%v err=%v", state, found, err)
	}
}
