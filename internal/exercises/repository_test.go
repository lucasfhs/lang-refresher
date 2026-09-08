package exercises

import (
	"testing"
	"testing/fstest"
)

func TestLoadTrack(t *testing.T) {
	data := `{
        "id":"test", "language":"python", "title":"Teste", "description":"x", "estimated_minutes":2,
        "exercises":[{"id":"one", "language":"python", "category":"core", "title":"Um", "topic":"x", "objective":"x", "description":"Faça", "requirements":[], "example_output":"ok", "starter_code":"", "estimated_minutes":2, "difficulty":1}]
    }`
	repo, err := Load(fstest.MapFS{"content/python/test.json": {Data: []byte(data)}}, "content/*/*.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	track, ok := repo.Track("test")
	if !ok || len(track.Exercises) != 1 || track.Exercises[0].ID != "one" {
		t.Fatalf("trilha carregada incorretamente: %#v", track)
	}
}

func TestLoadRejectsDuplicateExercises(t *testing.T) {
	data := `{
        "id":"test", "language":"python", "title":"Teste", "description":"x", "estimated_minutes":2,
        "exercises":[
          {"id":"same", "language":"python", "category":"core", "title":"Um", "topic":"x", "objective":"x", "description":"Faça", "requirements":[], "example_output":"ok", "starter_code":"", "estimated_minutes":2, "difficulty":1},
          {"id":"same", "language":"python", "category":"core", "title":"Dois", "topic":"x", "objective":"x", "description":"Faça", "requirements":[], "example_output":"ok", "starter_code":"", "estimated_minutes":2, "difficulty":1}
        ]
    }`
	_, err := Load(fstest.MapFS{"content/python/test.json": {Data: []byte(data)}}, "content/*/*.json")
	if err == nil {
		t.Fatal("Load() deveria rejeitar ids duplicados")
	}
}
