package generator

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/melkeydev/go-blueprint/internal/spec"
)

func config(root string) spec.Config {
	value := spec.DefaultConfig()
	value.Name = "acme"
	value.Module = "example.com/acme"
	value.OutputDir = root
	return value
}

func TestCreateWritesDeterministicManifest(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	generator := New(RendererFunc(func(spec.Config) (map[string][]byte, error) {
		return map[string][]byte{"z.txt": []byte("z"), "a.txt": []byte("a")}, nil
	}))
	generator.now = func() time.Time { return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC) }
	result, err := generator.Create(config(root))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a.txt", "z.txt"}; !reflect.DeepEqual(result.Files, want) {
		t.Fatalf("files = %v, want %v", result.Files, want)
	}
	manifest, err := LoadManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.GeneratedAt != generator.now() {
		t.Fatalf("generatedAt = %v", manifest.GeneratedAt)
	}
	if len(manifest.Files) != 2 || manifest.Files[0].Path != "a.txt" {
		t.Fatalf("managed files = %#v", manifest.Files)
	}
	if err := Verify(root); err != nil {
		t.Fatalf("Verify() = %v", err)
	}
}

func TestCreateDryRunDoesNotWrite(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	value := config(root)
	value.DryRun = true
	result, err := New(nil).Create(value)
	if err != nil {
		t.Fatal(err)
	}
	if !result.DryRun {
		t.Fatal("expected dry run")
	}
	if _, err := os.Stat(root); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("root should not exist, err = %v", err)
	}
}

func TestVerifyDetectsManagedFileModification(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	if _, err := New(nil).Create(config(root)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Verify(root); !errors.Is(err, ErrConflict) {
		t.Fatalf("Verify() = %v, want ErrConflict", err)
	}
}

func TestAddFeatureResolvesDependenciesAndProtectsChanges(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	value := config(root)
	value.Features = []spec.Feature{spec.FeatureJobs}
	generator := New(RendererFunc(func(config spec.Config) (map[string][]byte, error) {
		files := map[string][]byte{"base.txt": []byte("base")}
		if config.Has(spec.FeatureContent) {
			files["content.txt"] = []byte("content")
		}
		return files, nil
	}))
	if _, err := generator.Create(value); err != nil {
		t.Fatal(err)
	}
	result, err := generator.AddFeature(root, spec.FeatureContent, false)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Config.Has(spec.FeatureAuth) || !result.Config.Has(spec.FeatureRBAC) || !result.Config.Has(spec.FeatureContent) {
		t.Fatalf("dependencies not resolved: %v", result.Config.Features)
	}
	if _, err := os.Stat(filepath.Join(root, "content.txt")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "base.txt"), []byte("user change"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := generator.AddFeature(root, spec.FeatureCustomers, false); err != nil {
		t.Fatalf("AddFeature() = %v", err)
	}
	contents, err := os.ReadFile(filepath.Join(root, "base.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "user change" {
		t.Fatalf("customized file = %q", contents)
	}
	manifest, err := LoadManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if !manifest.Files[0].Customized {
		t.Fatalf("customized file was not recorded: %#v", manifest.Files)
	}
	if err := Verify(root); err != nil {
		t.Fatalf("Verify() after customization = %v", err)
	}
}

func TestRejectsUnsafeRendererPath(t *testing.T) {
	generator := New(RendererFunc(func(spec.Config) (map[string][]byte, error) {
		return map[string][]byte{"../escape": []byte("bad")}, nil
	}))
	if _, err := generator.Create(config(filepath.Join(t.TempDir(), "project"))); err == nil {
		t.Fatal("expected unsafe path error")
	}
}
