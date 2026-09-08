package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"language-refresher/internal/exercises"
	"language-refresher/internal/languages"
)

var ErrAlreadyRunning = errors.New("já existe uma execução em andamento")

type Request struct {
	Language     languages.Definition
	Code         string
	SupportFiles []exercises.SupportFile
	Validator    *exercises.Validator
	Timeout      time.Duration
}

type Result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exitCode"`
	DurationMS int64  `json:"durationMs"`
	TimedOut   bool   `json:"timedOut"`
	Cancelled  bool   `json:"cancelled"`
	Passed     bool   `json:"passed"`
}

type Runner struct {
	mu     sync.Mutex
	cancel context.CancelFunc
}

func NewRunner() *Runner { return &Runner{} }

func (r *Runner) Run(req Request) (Result, error) {
	if req.Timeout <= 0 {
		req.Timeout = 8 * time.Second
	}

	r.mu.Lock()
	if r.cancel != nil {
		r.mu.Unlock()
		return Result{}, ErrAlreadyRunning
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	r.cancel = cancel
	r.mu.Unlock()

	defer func() {
		cancel()
		r.mu.Lock()
		r.cancel = nil
		r.mu.Unlock()
	}()

	workDir, err := os.MkdirTemp("", "language-refresher-")
	if err != nil {
		return Result{}, fmt.Errorf("criar diretório temporário: %w", err)
	}
	defer os.RemoveAll(workDir)

	if err := writeWorkspace(workDir, req); err != nil {
		return Result{}, err
	}

	args := append([]string(nil), req.Language.Command[1:]...)
	entrypoint := req.Language.FileName
	if req.Validator != nil {
		entrypoint = "__validator__.py"
	}
	args = append(args, entrypoint)
	cmd := exec.CommandContext(ctx, req.Language.Command[0], args...)
	cmd.Dir = workDir
	// Força streams previsíveis no Windows (cujo console pode usar CP-1252),
	// preservando acentos e tracebacks como UTF-8 até o frontend.
	cmd.Env = append(os.Environ(), "PYTHONUTF8=1", "PYTHONIOENCODING=utf-8")
	configureProcess(cmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	started := time.Now()
	err = cmd.Run()
	duration := time.Since(started)

	result := Result{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		ExitCode:   0,
		DurationMS: duration.Milliseconds(),
	}
	if err != nil {
		result.ExitCode = -1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
		} else if ctx.Err() == nil {
			return result, fmt.Errorf("iniciar %s: %w", req.Language.DisplayName, err)
		}
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		result.TimedOut = true
	} else if errors.Is(ctx.Err(), context.Canceled) {
		result.Cancelled = true
	}
	result.Passed = req.Validator != nil && result.ExitCode == 0 && !result.TimedOut && !result.Cancelled
	return result, nil
}

func (r *Runner) Cancel() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cancel == nil {
		return false
	}
	r.cancel()
	return true
}

func writeWorkspace(workDir string, req Request) error {
	mainPath := filepath.Join(workDir, req.Language.FileName)
	if err := os.WriteFile(mainPath, []byte(req.Code), 0o600); err != nil {
		return fmt.Errorf("gravar código temporário: %w", err)
	}
	for _, support := range req.SupportFiles {
		clean := filepath.Clean(support.Path)
		if filepath.IsAbs(clean) || clean == "." || clean == ".." || filepath.Dir(clean) != "." {
			return fmt.Errorf("arquivo de apoio inválido: %s", support.Path)
		}
		if err := os.WriteFile(filepath.Join(workDir, clean), []byte(support.Content), 0o600); err != nil {
			return fmt.Errorf("gravar arquivo de apoio %s: %w", clean, err)
		}
	}
	if req.Validator != nil {
		if req.Validator.Type != "python_tests" {
			return fmt.Errorf("validador desconhecido: %s", req.Validator.Type)
		}
		validator := "import runpy\nsolution = runpy.run_path('" + req.Language.FileName + "', run_name='solution')\n" + req.Validator.TestCode + "\nprint('✓ Testes automáticos concluídos com sucesso.')\n"
		if err := os.WriteFile(filepath.Join(workDir, "__validator__.py"), []byte(validator), 0o600); err != nil {
			return fmt.Errorf("gravar validador temporário: %w", err)
		}
	}
	return nil
}
