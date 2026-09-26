package exercises

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

type Repository struct {
	tracks       map[string]Track
	translations map[string]map[string]TranslationCatalog
}

func Load(fsys fs.FS, pattern string) (*Repository, error) {
	paths, err := fs.Glob(fsys, pattern)
	if err != nil {
		return nil, fmt.Errorf("localizar trilhas: %w", err)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("nenhuma trilha encontrada em %q", pattern)
	}

	repo := &Repository{tracks: make(map[string]Track, len(paths)), translations: make(map[string]map[string]TranslationCatalog)}
	for _, path := range paths {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, fmt.Errorf("ler %s: %w", path, err)
		}
		var track Track
		if err := json.Unmarshal(data, &track); err != nil {
			return nil, fmt.Errorf("decodificar %s: %w", path, err)
		}
		if err := validateTrack(track); err != nil {
			return nil, fmt.Errorf("trilha %s: %w", filepath.Base(path), err)
		}
		if _, exists := repo.tracks[track.ID]; exists {
			return nil, fmt.Errorf("id de trilha duplicado: %s", track.ID)
		}
		repo.tracks[track.ID] = track
	}
	return repo, nil
}

func (r *Repository) LoadTranslations(fsys fs.FS, pattern string) error {
	paths, err := fs.Glob(fsys, pattern)
	if err != nil {
		return fmt.Errorf("localizar traduções: %w", err)
	}
	for _, path := range paths {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("ler %s: %w", path, err)
		}
		var catalog TranslationCatalog
		if err := json.Unmarshal(data, &catalog); err != nil {
			return fmt.Errorf("decodificar %s: %w", path, err)
		}
		track, ok := r.tracks[catalog.TrackID]
		if !ok {
			return fmt.Errorf("tradução %s referencia trilha desconhecida %q", path, catalog.TrackID)
		}
		if catalog.Locale == "" || catalog.Title == "" || len(catalog.Exercises) != len(track.Exercises) {
			return fmt.Errorf("tradução %s está incompleta", path)
		}
		for _, exercise := range track.Exercises {
			translated, ok := catalog.Exercises[exercise.ID]
			if !ok || translated.Title == "" || translated.Description == "" || translated.ExampleOutput == "" {
				return fmt.Errorf("tradução ausente para %s", exercise.ID)
			}
		}
		if r.translations[catalog.TrackID] == nil {
			r.translations[catalog.TrackID] = make(map[string]TranslationCatalog)
		}
		r.translations[catalog.TrackID][catalog.Locale] = catalog
	}
	return nil
}

func (r *Repository) LocalizedTrack(id, locale string) (Track, bool) {
	track, ok := r.tracks[id]
	if !ok || locale == "pt-BR" {
		return track, ok
	}
	catalog, ok := r.translations[id][locale]
	if !ok {
		return track, true
	}
	track.Title = catalog.Title
	track.Description = catalog.Description
	track.Exercises = append([]Exercise(nil), track.Exercises...)
	for i := range track.Exercises {
		translated := catalog.Exercises[track.Exercises[i].ID]
		exercise := &track.Exercises[i]
		exercise.Category = translated.Category
		exercise.Title = translated.Title
		exercise.Topic = translated.Topic
		exercise.Objective = translated.Objective
		exercise.Description = translated.Description
		exercise.Requirements = translated.Requirements
		exercise.Example = translated.Example
		exercise.ExampleOutput = translated.ExampleOutput
		exercise.Hint = translated.Hint
	}
	return track, true
}

func (r *Repository) Track(id string) (Track, bool) {
	track, ok := r.tracks[id]
	return track, ok
}

// Tracks returns every available practice in a deterministic order.
func (r *Repository) Tracks() []Track {
	tracks := make([]Track, 0, len(r.tracks))
	for _, track := range r.tracks {
		tracks = append(tracks, track)
	}
	sort.Slice(tracks, func(i, j int) bool { return tracks[i].ID < tracks[j].ID })
	return tracks
}

func validateTrack(track Track) error {
	if strings.TrimSpace(track.ID) == "" || strings.TrimSpace(track.Language) == "" || strings.TrimSpace(track.Title) == "" {
		return errors.New("id, language e title são obrigatórios")
	}
	if len(track.Exercises) == 0 {
		return errors.New("a trilha não possui exercícios")
	}
	seen := make(map[string]struct{}, len(track.Exercises))
	for i, exercise := range track.Exercises {
		if exercise.ID == "" || exercise.Title == "" || exercise.Description == "" || exercise.ExampleOutput == "" {
			return fmt.Errorf("exercício %d não possui campos obrigatórios", i+1)
		}
		if exercise.Language != track.Language {
			return fmt.Errorf("exercício %s usa linguagem %q, esperado %q", exercise.ID, exercise.Language, track.Language)
		}
		if _, ok := seen[exercise.ID]; ok {
			return fmt.Errorf("id de exercício duplicado: %s", exercise.ID)
		}
		seen[exercise.ID] = struct{}{}
		if exercise.EstimatedMinutes < 1 || exercise.EstimatedMinutes > 5 {
			return fmt.Errorf("exercício %s deve durar entre 1 e 5 minutos", exercise.ID)
		}
	}
	return nil
}
