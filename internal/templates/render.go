package templates

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/melkeydev/go-blueprint/internal/spec"
)

//go:embed backend frontend all:project backend/.air.toml.tmpl backend/.env.example.tmpl frontend/.env.example.tmpl
var assets embed.FS

type Renderer struct{}

type templateData struct {
	ProjectName string
	Name        string
	Module      string
	Database    string
	Seed        string
	Auth        bool
	RBAC        bool
	Jobs        bool
	Cron        bool
	Content     bool
	Customers   bool
	Sales       bool
	Workspace   bool
	Audit       bool
	Files       bool
	Email       bool
	Cache       bool
	Frontend    bool
}

func (Renderer) Render(config spec.Config) (map[string][]byte, error) {
	data := templateData{
		ProjectName: config.Name,
		Name:        config.Name,
		Module:      config.Module,
		Database:    string(config.Database),
		Seed:        string(config.Seed),
		Auth:        config.Has(spec.FeatureAuth),
		RBAC:        config.Has(spec.FeatureRBAC),
		Jobs:        config.Has(spec.FeatureJobs),
		Cron:        config.Has(spec.FeatureCron),
		Content:     config.Has(spec.FeatureContent),
		Customers:   config.Has(spec.FeatureCustomers),
		Sales:       config.Has(spec.FeatureSales),
		Workspace:   config.Has(spec.FeatureWorkspace),
		Audit:       config.Has(spec.FeatureAudit),
		Files:       config.Has(spec.FeatureFiles),
		Email:       config.Has(spec.FeatureEmail),
		Cache:       config.Has(spec.FeatureCache),
		Frontend:    config.Frontend,
	}
	files := make(map[string][]byte)
	if err := renderTree(files, "project", data, config); err != nil {
		return nil, err
	}
	if err := renderTree(files, "backend", data, config); err != nil {
		return nil, err
	}
	if config.Frontend {
		if err := renderTree(files, "frontend", data, config); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func renderTree(output map[string][]byte, root string, data templateData, config spec.Config) error {
	return fs.WalkDir(assets, root, func(assetPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(assetPath, ".tmpl") {
			return nil
		}
		relative := strings.TrimPrefix(assetPath, root+"/")
		if !includeFeatureAsset(root, relative, config) {
			return nil
		}
		if root == "backend" && !includeBackend(relative, config.Database) {
			return nil
		}
		target := path.Join(root, strings.TrimSuffix(relative, ".tmpl"))
		if root == "project" {
			target = strings.TrimPrefix(target, "project/")
		}
		if root == "backend" {
			target = backendTarget(target, config.Database)
		}
		contents, err := assets.ReadFile(assetPath)
		if err != nil {
			return err
		}
		parsed, err := template.New(assetPath).Option("missingkey=error").Parse(string(contents))
		if err != nil {
			return fmt.Errorf("parse %s: %w", assetPath, err)
		}
		var rendered bytes.Buffer
		if err := parsed.Execute(&rendered, data); err != nil {
			return fmt.Errorf("render %s: %w", assetPath, err)
		}
		output[target] = rendered.Bytes()
		return nil
	})
}

func includeFeatureAsset(root, relative string, config spec.Config) bool {
	if root == "frontend" {
		if strings.HasPrefix(relative, "src/features/auth/") {
			return config.Has(spec.FeatureAuth)
		}
		if strings.HasPrefix(relative, "src/features/cms/") {
			return config.Has(spec.FeatureContent)
		}
		if strings.HasPrefix(relative, "src/features/crm/") {
			return config.Has(spec.FeatureCustomers) || config.Has(spec.FeatureSales)
		}
		if strings.HasPrefix(relative, "src/features/organizations/") {
			return config.Has(spec.FeatureWorkspace)
		}
		if strings.HasPrefix(relative, "src/features/audit/") {
			return config.Has(spec.FeatureAudit)
		}
		if strings.HasPrefix(relative, "src/features/files/") {
			return config.Has(spec.FeatureFiles)
		}
		if strings.HasPrefix(relative, "src/features/operations/") {
			return config.Has(spec.FeatureJobs)
		}
		return true
	}
	if root != "backend" {
		return true
	}
	if strings.HasPrefix(relative, "internal/features/users/") || strings.HasPrefix(relative, "internal/features/auth/") || strings.HasPrefix(relative, "internal/features/organizations/") {
		return config.Has(spec.FeatureAuth)
	}
	if strings.HasPrefix(relative, "internal/features/cms/") {
		return config.Has(spec.FeatureContent)
	}
	if strings.HasPrefix(relative, "internal/features/audit/") {
		return config.Has(spec.FeatureAudit)
	}
	if strings.HasPrefix(relative, "internal/features/files/") {
		return config.Has(spec.FeatureFiles)
	}
	if strings.HasPrefix(relative, "internal/features/email/") {
		return config.Has(spec.FeatureEmail)
	}
	if strings.HasPrefix(relative, "internal/features/jobs/") {
		return config.Has(spec.FeatureJobs)
	}
	if strings.HasPrefix(relative, "internal/features/jobs/") || strings.HasPrefix(relative, "internal/platform/jobs/") {
		return config.Has(spec.FeatureJobs)
	}
	if strings.HasPrefix(relative, "internal/platform/cache/") {
		return config.Has(spec.FeatureCache)
	}
	if strings.HasPrefix(relative, "internal/features/crm/") {
		return config.Has(spec.FeatureCustomers) || config.Has(spec.FeatureSales)
	}
	if strings.HasPrefix(relative, "cmd/worker/") {
		return config.Has(spec.FeatureJobs)
	}
	if strings.HasPrefix(relative, "cmd/scheduler/") {
		return config.Has(spec.FeatureCron)
	}
	base := path.Base(relative)
	switch {
	case strings.HasPrefix(base, "00001_"), strings.HasPrefix(base, "00002_"), strings.HasPrefix(base, "00003_"):
		return config.Has(spec.FeatureAuth)
	case strings.HasPrefix(base, "00004_"), strings.HasPrefix(base, "00014_"):
		return config.Has(spec.FeatureJobs)
	case strings.HasPrefix(base, "00006_"), strings.HasPrefix(base, "00007_"), strings.HasPrefix(base, "00008_"), strings.HasPrefix(base, "00009_"):
		return config.Has(spec.FeatureRBAC)
	case strings.HasPrefix(base, "00010_"), strings.HasPrefix(base, "00011_"), strings.HasPrefix(base, "00012_"), strings.HasPrefix(base, "00013_"):
		return config.Has(spec.FeatureAuth)
	case strings.HasPrefix(base, "0002"):
		return config.Has(spec.FeatureContent)
	case strings.HasPrefix(base, "0003"):
		return config.Has(spec.FeatureCustomers)
	default:
		return true
	}
}

func includeBackend(relative string, database spec.Database) bool {
	if strings.HasPrefix(relative, "migrations/postgres/") {
		return database == spec.DatabasePostgres
	}
	if strings.HasPrefix(relative, "migrations/mysql/") {
		return database == spec.DatabaseMySQL
	}
	if strings.HasPrefix(relative, "docker-compose.postgres.") {
		return database == spec.DatabasePostgres
	}
	if strings.HasPrefix(relative, "docker-compose.mysql.") {
		return database == spec.DatabaseMySQL
	}
	return true
}

func backendTarget(target string, database spec.Database) string {
	dialectMigrations := "backend/migrations/" + string(database) + "/"
	if strings.HasPrefix(target, dialectMigrations) {
		target = "backend/migrations/" + strings.TrimPrefix(target, dialectMigrations)
	}
	compose := "backend/docker-compose." + string(database) + ".yml"
	if target == compose {
		return "backend/docker-compose.yml"
	}
	return target
}
