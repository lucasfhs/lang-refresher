package execution

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"language-refresher/internal/exercises"
	"language-refresher/internal/languages"
)

func pythonDefinition(t *testing.T) languages.Definition {
	t.Helper()
	python, err := exec.LookPath("python")
	if err != nil {
		t.Skip("Python não está instalado")
	}
	return languages.Definition{ID: "python", DisplayName: "Python", FileName: "main.py", Command: []string{python}}
}

func TestRunnerCapturesOutput(t *testing.T) {
	result, err := NewRunner().Run(Request{Language: pythonDefinition(t), Code: "import sys\nprint('out')\nprint('err', file=sys.stderr)\n", Timeout: 2 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 || !strings.Contains(result.Stdout, "out") || !strings.Contains(result.Stderr, "err") {
		t.Fatalf("resultado inesperado: %#v", result)
	}
}

func TestRunnerTimesOut(t *testing.T) {
	result, err := NewRunner().Run(Request{Language: pythonDefinition(t), Code: "while True: pass\n", Timeout: 100 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if !result.TimedOut {
		t.Fatalf("esperava timeout: %#v", result)
	}
}

func TestRunnerValidatesObjectsWithoutStdoutComparison(t *testing.T) {
	validator := &exercises.Validator{Type: "python_tests", TestCode: "assert solution['dobro'](4) == 8"}
	result, err := NewRunner().Run(Request{Language: pythonDefinition(t), Code: "def dobro(n):\n    return n * 2\n", Validator: validator, Timeout: 2 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Fatalf("validação deveria passar: %#v", result)
	}
}
