package acceptance

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CommandRunner makes the acceptance lifecycle testable without invoking Docker.
type CommandRunner interface {
	Run(context.Context, string, string, ...string) error
}

type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, dir, name string, args ...string) error {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = dir
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	return command.Run()
}

type Harness struct {
	Runner          CommandRunner
	BlueprintBinary string
	ReadinessProbe  func(context.Context, Scenario, string) error
	Timeout         time.Duration
}

// Run generates, verifies, boots, probes, and always tears down an isolated
// Compose project. The generated project owns the concrete verification scripts.
func (h Harness) Run(ctx context.Context, scenario Scenario, workspace string) (runErr error) {
	if err := scenario.Validate(); err != nil {
		return err
	}
	if h.Runner == nil || h.BlueprintBinary == "" || h.ReadinessProbe == nil {
		return errors.New("acceptance harness is not fully configured")
	}
	projectDir := filepath.Join(workspace, scenario.Name)
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return fmt.Errorf("create acceptance workspace: %w", err)
	}
	if err := h.Runner.Run(ctx, workspace, h.BlueprintBinary, scenario.GeneratorArgs(projectDir)...); err != nil {
		return fmt.Errorf("generate project: %w", err)
	}
	projectName := "blueprint-" + scenario.Name
	compose := []string{"compose", "--project-name", projectName}
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		cleanupErr := h.Runner.Run(cleanupCtx, projectDir, "docker", append(compose, "down", "--volumes", "--remove-orphans")...)
		if runErr == nil && cleanupErr != nil {
			runErr = fmt.Errorf("cleanup acceptance environment: %w", cleanupErr)
		}
	}()

	steps := [][]string{
		{"make", "verify-generated"},
		{"docker", "compose", "--project-name", projectName, "up", "--detach", "--build"},
	}
	for _, step := range steps {
		if err := h.Runner.Run(ctx, projectDir, step[0], step[1:]...); err != nil {
			return fmt.Errorf("run %v: %w", step, err)
		}
	}
	probeCtx := ctx
	if h.Timeout > 0 {
		var cancel context.CancelFunc
		probeCtx, cancel = context.WithTimeout(ctx, h.Timeout)
		defer cancel()
	}
	if err := h.ReadinessProbe(probeCtx, scenario, projectDir); err != nil {
		return fmt.Errorf("readiness/behavior probe: %w", err)
	}
	return nil
}
