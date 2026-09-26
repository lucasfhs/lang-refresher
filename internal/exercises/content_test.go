package exercises

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestPythonCoreCatalog(t *testing.T) {
	repo, err := Load(os.DirFS("../.."), "content/*/*-refresher.json")
	if err != nil {
		t.Fatalf("catálogo real inválido: %v", err)
	}
	if err := repo.LoadTranslations(os.DirFS("../.."), "content/*/*.i18n.json"); err != nil {
		t.Fatalf("traduções inválidas: %v", err)
	}
	track, ok := repo.Track("python-core-refresher")
	if !ok {
		t.Fatal("Python Core Refresher não encontrada")
	}
	if len(track.Exercises) < 30 || len(track.Exercises) > 50 {
		t.Fatalf("esperava 30 a 50 exercícios, recebeu %d", len(track.Exercises))
	}
	totalMinutes := 0
	validatorCount := 0
	categories := map[string]bool{}
	for _, exercise := range track.Exercises {
		totalMinutes += exercise.EstimatedMinutes
		categories[exercise.Category] = true
		if strings.TrimSpace(exercise.ExampleOutput) == "" {
			t.Fatalf("exercício %s não possui exemplo de saída", exercise.ID)
		}
		if exercise.Validator != nil {
			validatorCount++
		}
	}
	if totalMinutes < 100 || totalMinutes > 150 {
		t.Fatalf("duração estimada fora da sessão de ~2h: %d minutos", totalMinutes)
	}
	if len(categories) < 12 {
		t.Fatalf("cobertura de categorias insuficiente: %v", categories)
	}
	if validatorCount < 10 {
		t.Fatalf("poucos exercícios com validação automática: %d", validatorCount)
	}
	english, ok := repo.LocalizedTrack(track.ID, "en")
	if !ok || english.Title != "Python Core Refresher" || len(english.Exercises) != len(track.Exercises) {
		t.Fatalf("trilha em inglês inválida: %#v", english)
	}
	for i, exercise := range english.Exercises {
		if exercise.Title == track.Exercises[i].Title || exercise.Description == "" || exercise.ExampleOutput == "" {
			t.Fatalf("exercício %s não foi traduzido por completo", exercise.ID)
		}
	}
}

func TestNumPyCatalog(t *testing.T) {
	repo, err := Load(os.DirFS("../.."), "content/*/*-refresher.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.LoadTranslations(os.DirFS("../.."), "content/*/*.i18n.json"); err != nil {
		t.Fatal(err)
	}
	track, ok := repo.Track("numpy-refresher")
	if !ok {
		t.Fatal("trilha NumPy não encontrada")
	}
	if len(track.Exercises) != 45 {
		t.Fatalf("esperava 45 exercícios NumPy, recebeu %d", len(track.Exercises))
	}
	totalMinutes := 0
	categories := map[string]bool{}
	for _, exercise := range track.Exercises {
		totalMinutes += exercise.EstimatedMinutes
		categories[exercise.Category] = true
		if exercise.Validator == nil {
			t.Fatalf("exercício NumPy sem validação: %s", exercise.ID)
		}
	}
	if totalMinutes < 120 || totalMinutes > 150 {
		t.Fatalf("duração NumPy inesperada: %d", totalMinutes)
	}
	if len(categories) < 12 {
		t.Fatalf("cobertura NumPy insuficiente: %v", categories)
	}
	english, ok := repo.LocalizedTrack(track.ID, "en")
	if !ok || len(english.Exercises) != 45 || english.Title == track.Title {
		t.Fatalf("tradução NumPy inválida: %#v", english)
	}
}

func TestTrackSnippetsCompile(t *testing.T) {
	python, err := exec.LookPath("python")
	if err != nil {
		t.Skip("Python não está instalado")
	}
	repo, err := Load(os.DirFS("../.."), "content/*/*-refresher.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, track := range repo.Tracks() {
		for _, exercise := range track.Exercises {
			t.Run(exercise.ID, func(t *testing.T) {
				assertPythonCompiles(t, python, exercise.StarterCode)
				if exercise.Validator != nil {
					validator := "import runpy\nsolution = runpy.run_path('main.py', run_name='solution')\n" + exercise.Validator.TestCode
					assertPythonCompiles(t, python, validator)
				}
			})
		}
	}
}

func assertPythonCompiles(t *testing.T, python, source string) {
	t.Helper()
	cmd := exec.Command(python, "-c", "import sys; compile(sys.stdin.read(), '<content>', 'exec')")
	cmd.Stdin = strings.NewReader(source)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Python inválido: %v\n%s\n--- fonte ---\n%s", err, output, fmt.Sprintf("%s", source))
	}
}
