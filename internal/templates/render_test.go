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
