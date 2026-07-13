package acceptance

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestPullRequestScenariosAreValidAndCoverOptions(t *testing.T) {
	covered := map[string]bool{}
	for _, scenario := range PullRequestScenarios() {
		if err := scenario.Validate(); err != nil {
			t.Fatalf("%s: %v", scenario.Name, err)
		}
		covered[scenario.Database] = true
		for _, feature := range scenario.Features {
			covered[feature] = true
		}
	}
	for _, option := range []string{"postgres", "mysql", "auth", "rbac", "jobs", "cron", "cms", "crm"} {
		if !covered[option] {
			t.Errorf("option %q has no PR scenario", option)
		}
	}
}

func TestScenarioGeneratorArgsAreStable(t *testing.T) {
	s := Scenario{Name: "sample", Database: "postgres", Seed: "minimal", Frontend: true, Features: []string{"rbac", "auth"}}
	args := s.GeneratorArgs("out")
	if !slices.Equal(args[len(args)-4:], []string{"--feature", "auth", "--feature", "rbac"}) {
		t.Fatalf("features are not sorted: %v", args)
	}
}

type recordingRunner struct{ calls []string }

func (r *recordingRunner) Run(_ context.Context, dir, name string, args ...string) error {
	r.calls = append(r.calls, dir+"|"+name+" "+strings.Join(args, " "))
	if name == "make" {
		return errors.New("verification failed")
	}
	return nil
}

func TestHarnessAlwaysCleansUpAfterFailure(t *testing.T) {
	runner := &recordingRunner{}
	h := Harness{Runner: runner, BlueprintBinary: "blueprint", ReadinessProbe: func(context.Context, Scenario, string) error { return nil }}
	err := h.Run(context.Background(), Scenario{Name: "sample", Database: "postgres", Seed: "none"}, t.TempDir())
	if err == nil {
		t.Fatal("expected verification failure")
	}
	if got := runner.calls[len(runner.calls)-1]; !strings.Contains(got, "docker compose --project-name blueprint-sample down --volumes --remove-orphans") {
		t.Fatalf("last call was not cleanup: %s", got)
	}
}

func TestNightlyScenariosCoverEveryValidOptionPair(t *testing.T) {
	want := map[string]bool{}
	for _, scenario := range allCandidates() {
		for pair := range optionPairs(scenario) {
			want[pair] = true
		}
	}
	for _, scenario := range NightlyScenarios() {
		if err := scenario.Validate(); err != nil {
			t.Fatalf("%s: %v", scenario.Name, err)
		}
		for pair := range optionPairs(scenario) {
			delete(want, pair)
		}
	}
	if len(want) != 0 {
		t.Fatalf("nightly matrix misses %d valid pairs", len(want))
	}
}
