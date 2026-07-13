package spec

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Database string

const (
	DatabasePostgres Database = "postgres"
	DatabaseMySQL    Database = "mysql"
)

type SeedProfile string

const (
	SeedNone    SeedProfile = "none"
	SeedMinimal SeedProfile = "minimal"
	SeedDemo    SeedProfile = "demo"
)

type Feature string

const (
	FeatureAuth Feature = "auth"
	FeatureRBAC Feature = "rbac"
	FeatureJobs Feature = "jobs"
	FeatureCron Feature = "cron"
	FeatureCMS  Feature = "cms"
	FeatureCRM  Feature = "crm"
)

var featureDependencies = map[Feature][]Feature{
	FeatureRBAC: {FeatureAuth},
	FeatureCron: {FeatureJobs},
	FeatureCMS:  {FeatureAuth, FeatureRBAC},
	FeatureCRM:  {FeatureAuth, FeatureRBAC},
}

var featureOrder = []Feature{
	FeatureAuth,
	FeatureRBAC,
	FeatureJobs,
	FeatureCron,
	FeatureCMS,
	FeatureCRM,
}

type Config struct {
	Name      string
	Module    string
	Database  Database
	Seed      SeedProfile
	Frontend  bool
	Docker    bool
	Features  []Feature
	OutputDir string
	DryRun    bool
}

func DefaultConfig() Config {
	return Config{
		Database: DatabasePostgres,
		Seed:     SeedMinimal,
		Frontend: true,
		Docker:   true,
		Features: []Feature{FeatureAuth, FeatureRBAC, FeatureJobs, FeatureCron},
	}
}

func Resolve(input Config) (Config, error) {
	if strings.TrimSpace(input.Name) == "" {
		return Config{}, errors.New("project name is required")
	}
	if strings.TrimSpace(input.Module) == "" {
		return Config{}, errors.New("Go module path is required")
	}
	if input.Database != DatabasePostgres && input.Database != DatabaseMySQL {
		return Config{}, fmt.Errorf("unsupported database %q", input.Database)
	}
	if input.Seed != SeedNone && input.Seed != SeedMinimal && input.Seed != SeedDemo {
		return Config{}, fmt.Errorf("unsupported seed profile %q", input.Seed)
	}

	selected := make(map[Feature]bool, len(input.Features))
	var include func(Feature) error
	include = func(feature Feature) error {
		if !slices.Contains(featureOrder, feature) {
			return fmt.Errorf("unsupported feature %q", feature)
		}
		if selected[feature] {
			return nil
		}
		for _, dependency := range featureDependencies[feature] {
			if err := include(dependency); err != nil {
				return err
			}
		}
		selected[feature] = true
		return nil
	}
	for _, feature := range input.Features {
		if err := include(feature); err != nil {
			return Config{}, err
		}
	}

	input.Features = input.Features[:0]
	for _, feature := range featureOrder {
		if selected[feature] {
			input.Features = append(input.Features, feature)
		}
	}
	return input, nil
}

func (config Config) Has(feature Feature) bool {
	return slices.Contains(config.Features, feature)
}

func (config Config) Auth() bool { return config.Has(FeatureAuth) }
func (config Config) RBAC() bool { return config.Has(FeatureRBAC) }
func (config Config) Jobs() bool { return config.Has(FeatureJobs) }
func (config Config) Cron() bool { return config.Has(FeatureCron) }
func (config Config) CMS() bool  { return config.Has(FeatureCMS) }
func (config Config) CRM() bool  { return config.Has(FeatureCRM) }
