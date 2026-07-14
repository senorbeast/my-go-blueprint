package templates

import (
	"strings"
	"testing"

	"github.com/melkeydev/go-blueprint/internal/spec"
)

func TestRenderSelectsDialectAndFrontend(t *testing.T) {
	config := spec.DefaultConfig()
	config.Name = "acme"
	config.Module = "example.com/acme"
	config.Database = spec.DatabaseMySQL

	files, err := (Renderer{}).Render(config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if _, ok := files["backend/docker-compose.yml"]; !ok {
		t.Fatal("selected Compose file was not rendered")
	}
	if _, ok := files["README.md"]; !ok {
		t.Fatal("root README was not rendered")
	}
	if _, ok := files["Makefile"]; !ok {
		t.Fatal("root verification commands were not rendered")
	}
	if _, ok := files[".github/workflows/ci.yml"]; !ok {
		t.Fatal("generated CI workflow was not rendered")
	}
	if _, ok := files["frontend/package.json"]; !ok {
		t.Fatal("frontend package was not rendered")
	}
	for name, contents := range files {
		if strings.Contains(name, "postgres") || strings.Contains(string(contents), "{{") {
			t.Fatalf("unexpected unresolved/dialect output in %s", name)
		}
	}
}

func TestRenderCanOmitFrontend(t *testing.T) {
	config := spec.DefaultConfig()
	config.Name = "acme"
	config.Module = "example.com/acme"
	config.Frontend = false
	config.Features = nil
	files, err := (Renderer{}).Render(config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for name := range files {
		if strings.HasPrefix(name, "frontend/") {
			t.Fatalf("frontend file rendered: %s", name)
		}
		if strings.Contains(name, "/features/auth/") || strings.Contains(name, "00001_users") {
			t.Fatalf("auth file rendered without auth: %s", name)
		}
	}
}

func TestRenderBackendVerificationRegeneratesSQLCBeforeCompilation(t *testing.T) {
	config := spec.DefaultConfig()
	config.Name = "acme"
	config.Module = "example.com/acme"
	config.Database = spec.DatabasePostgres
	config.Frontend = false

	files, err := (Renderer{}).Render(config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	makefile, ok := files["backend/Makefile"]
	if !ok {
		t.Fatal("generated backend is missing its verification command")
	}
	contents := string(makefile)
	generate := strings.Index(contents, "go tool sqlc generate")
	compile := strings.Index(contents, "go test ./...")
	if generate < 0 || compile < 0 || generate > compile {
		t.Fatalf("backend verification must regenerate sqlc output before compilation:\\n%s", contents)
	}
	goModule, ok := files["backend/go.mod"]
	if !ok || !strings.Contains(string(goModule), "tool github.com/sqlc-dev/sqlc/cmd/sqlc") {
		t.Fatalf("generated backend must pin the sqlc Go tool:\\n%s", goModule)
	}
	configFile, ok := files["backend/sqlc.yaml"]
	if !ok || strings.Contains(string(configFile), "migrations/postgres") || !strings.Contains(string(configFile), "internal/platform/database/queries") {
		t.Fatalf("sqlc must use the rendered migration directory:\\n%s", configFile)
	}
}

func TestRenderIncludesOnlySelectedBusinessPack(t *testing.T) {
	config := spec.DefaultConfig()
	config.Name = "acme"
	config.Module = "example.com/acme"
	config.Features = []spec.Feature{spec.FeatureCMS}
	config, err := spec.Resolve(config)
	if err != nil {
		t.Fatal(err)
	}
	files, err := (Renderer{}).Render(config)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := files["backend/internal/features/cms/feature.go"]; !ok {
		t.Fatal("CMS feature missing")
	}
	for name := range files {
		if strings.Contains(name, "/features/crm/") || strings.Contains(name, "00030_crm") {
			t.Fatalf("CRM file rendered without CRM: %s", name)
		}
	}
}
