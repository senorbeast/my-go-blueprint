package backend_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/melkeydev/go-blueprint/internal/spec"
)

func TestTemplatesParse(t *testing.T) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") { return nil }
		contents, err := os.ReadFile(path)
		if err != nil { return err }
		if _, err := template.New(path).Parse(string(contents)); err != nil { t.Errorf("parse %s: %v", path, err) }
		return nil
	})
	if err != nil { t.Fatal(err) }
}

func TestTemplatesExecuteWithSpecConfig(t *testing.T) {
	config := spec.DefaultConfig()
	config.Name = "example"
	config.Module = "example.com/acme/example"
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") { return nil }
		contents, err := os.ReadFile(path)
		if err != nil { return err }
		parsed, err := template.New(path).Option("missingkey=error").Parse(string(contents))
		if err != nil { return err }
		if err := parsed.Execute(&bytes.Buffer{}, config); err != nil { t.Errorf("execute %s: %v", path, err) }
		return nil
	})
	if err != nil { t.Fatal(err) }
}

func TestMigrationsAreOneTableAndReversible(t *testing.T) {
	for _, dialect := range []string{"postgres", "mysql"} {
		matches, err := filepath.Glob(filepath.Join("migrations", dialect, "*.sql.tmpl"))
		if err != nil { t.Fatal(err) }
		if len(matches) == 0 { t.Fatalf("no %s migrations", dialect) }
		for _, path := range matches {
			contents, err := os.ReadFile(path)
			if err != nil { t.Fatal(err) }
			text := strings.ToUpper(string(contents))
			if strings.Count(text, "CREATE TABLE") != 1 { t.Errorf("%s must create exactly one table", path) }
			if !strings.Contains(text, "-- +GOOSE UP") || !strings.Contains(text, "-- +GOOSE DOWN") { t.Errorf("%s must have Up and Down", path) }
		}
	}
}
