package acceptance

import (
	"fmt"
	"sort"
	"strings"
)

// Scenario describes one generated-project acceptance environment.
type Scenario struct {
	Name     string
	Database string
	Frontend bool
	Seed     string
	Features []string
	Browser  bool
}

func (s Scenario) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("scenario name is required")
	}
	if s.Database != "postgres" && s.Database != "mysql" {
		return fmt.Errorf("scenario %q has unsupported database %q", s.Name, s.Database)
	}
	if s.Seed != "none" && s.Seed != "minimal" && s.Seed != "demo" {
		return fmt.Errorf("scenario %q has unsupported seed %q", s.Name, s.Seed)
	}
	if s.Browser && !s.Frontend {
		return fmt.Errorf("scenario %q enables browser tests without a frontend", s.Name)
	}
	selected := make(map[string]bool, len(s.Features))
	for _, feature := range s.Features {
		if selected[feature] {
			return fmt.Errorf("scenario %q repeats feature %q", s.Name, feature)
		}
		selected[feature] = true
	}
	dependencies := map[string][]string{
		"rbac": {"auth"}, "cron": {"jobs"},
		"content": {"auth", "rbac"}, "customers": {"auth", "rbac"}, "sales": {"customers"},
		"workspace": {"auth", "rbac"}, "audit": {"auth", "rbac"}, "files": {"auth", "rbac"}, "email": {"jobs"},
	}
	for feature, required := range dependencies {
		if !selected[feature] {
			continue
		}
		for _, dependency := range required {
			if !selected[dependency] {
				return fmt.Errorf("scenario %q: feature %q requires %q", s.Name, feature, dependency)
			}
		}
	}
	return nil
}

func (s Scenario) GeneratorArgs(output string) []string {
	args := []string{"create", "--name", s.Name, "--module", "example.com/" + s.Name,
		"--database", s.Database, "--seed", s.Seed, "--output", output}
	if s.Frontend {
		args = append(args, "--frontend")
	} else {
		args = append(args, "--no-frontend")
	}
	features := append([]string(nil), s.Features...)
	sort.Strings(features)
	for _, feature := range features {
		args = append(args, "--feature", feature)
	}
	return args
}

// PullRequestScenarios is deliberately bounded while covering both dialects
// and every feature at least once. Full pairwise expansion belongs in nightly CI.
func PullRequestScenarios() []Scenario {
	return []Scenario{
		{Name: "pg-default", Database: "postgres", Frontend: true, Seed: "minimal", Features: []string{"auth", "rbac", "workspace", "jobs", "cron"}, Browser: true},
		{Name: "mysql-full", Database: "mysql", Frontend: true, Seed: "demo", Features: []string{"auth", "rbac", "workspace", "jobs", "cron", "content", "customers", "sales", "audit", "files", "email"}, Browser: true},
		{Name: "minimal-api", Database: "postgres", Frontend: false, Seed: "none"},
		{Name: "pg-content", Database: "postgres", Frontend: true, Seed: "demo", Features: []string{"auth", "rbac", "content"}, Browser: true},
		{Name: "mysql-customers", Database: "mysql", Frontend: true, Seed: "demo", Features: []string{"auth", "rbac", "customers"}, Browser: true},
		{Name: "pg-jobs", Database: "postgres", Frontend: false, Seed: "minimal", Features: []string{"jobs", "cron"}},
		{Name: "pg-cache", Database: "postgres", Frontend: false, Seed: "minimal", Features: []string{"cache"}},
	}
}

// NightlyScenarios greedily selects a deterministic pairwise-covering matrix
// from every valid database/frontend/seed/feature combination.
func NightlyScenarios() []Scenario {
	candidates := allCandidates()
	uncovered := map[string]bool{}
	for _, candidate := range candidates {
		for pair := range optionPairs(candidate) {
			uncovered[pair] = true
		}
	}
	selected := make([]Scenario, 0)
	for len(uncovered) > 0 {
		bestIndex, bestCoverage := -1, -1
		for index, candidate := range candidates {
			coverage := 0
			for pair := range optionPairs(candidate) {
				if uncovered[pair] {
					coverage++
				}
			}
			if coverage > bestCoverage {
				bestIndex, bestCoverage = index, coverage
			}
		}
		if bestIndex < 0 || bestCoverage == 0 {
			break
		}
		chosen := candidates[bestIndex]
		chosen.Name = fmt.Sprintf("nightly-%02d-%s", len(selected)+1, chosen.Name)
		selected = append(selected, chosen)
		for pair := range optionPairs(chosen) {
			delete(uncovered, pair)
		}
		candidates = append(candidates[:bestIndex], candidates[bestIndex+1:]...)
	}
	return selected
}

func allCandidates() []Scenario {
	var candidates []Scenario
	for graphIndex, selected := range ResolvedFeatureGraphs() {
		for _, database := range []string{"postgres", "mysql"} {
			for _, frontend := range []bool{false, true} {
				for _, seed := range []string{"none", "minimal", "demo"} {
					scenario := Scenario{Name: fmt.Sprintf("%s-%t-%s-%02x", database, frontend, seed, graphIndex), Database: database, Frontend: frontend, Seed: seed, Features: selected}
					if scenario.Validate() == nil {
						candidates = append(candidates, scenario)
					}
				}
			}
		}
	}
	return candidates
}

// ResolvedFeatureGraphs is the generator's canonical valid feature graph set.
// It keeps acceptance coverage aligned with dependency resolution: the
// feature pack dependencies. It keeps acceptance coverage aligned with the
// generator without forcing the acceptance suite to enumerate every subset.
func ResolvedFeatureGraphs() [][]string {
	workspace := [][]string{
		nil,
		{"auth"},
		{"auth", "rbac"},
		{"auth", "rbac", "workspace"},
		{"auth", "rbac", "content"},
		{"auth", "rbac", "customers"},
		{"auth", "rbac", "customers", "sales"},
		{"auth", "rbac", "audit", "files"},
		{"auth", "rbac", "workspace", "content", "customers", "sales", "audit", "files"},
	}
	automation := [][]string{nil, {"jobs"}, {"jobs", "cron"}, {"jobs", "email"}, {"jobs", "cron", "email"}}
	graphs := make([][]string, 0, len(workspace)*len(automation))
	for _, left := range workspace {
		for _, right := range automation {
			features := append(append([]string(nil), left...), right...)
			graphs = append(graphs, features)
		}
	}
	return graphs
}

func optionPairs(scenario Scenario) map[string]bool {
	selected := make(map[string]bool, len(scenario.Features))
	for _, feature := range scenario.Features {
		selected[feature] = true
	}
	values := []string{
		"database=" + scenario.Database,
		fmt.Sprintf("frontend=%t", scenario.Frontend),
		"seed=" + scenario.Seed,
	}
	for _, feature := range []string{"auth", "rbac", "workspace", "jobs", "cron", "content", "customers", "sales", "audit", "files", "email"} {
		values = append(values, fmt.Sprintf("%s=%t", feature, selected[feature]))
	}
	pairs := make(map[string]bool)
	for left := 0; left < len(values); left++ {
		for right := left + 1; right < len(values); right++ {
			pairs[values[left]+"|"+values[right]] = true
		}
	}
	return pairs
}
