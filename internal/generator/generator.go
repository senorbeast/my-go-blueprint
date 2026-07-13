package generator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/melkeydev/go-blueprint/internal/spec"
)

const (
	ManifestPath  = ".blueprint/manifest.json"
	GeneratorName = "opinionated-go-blueprint"
)

var ErrConflict = errors.New("generated file has been modified")

// Renderer isolates feature composition from filesystem mutation. Template packs can
// replace it without changing the CLI or manifest safety rules.
type Renderer interface {
	Render(spec.Config) (map[string][]byte, error)
}

type RendererFunc func(spec.Config) (map[string][]byte, error)

func (f RendererFunc) Render(config spec.Config) (map[string][]byte, error) { return f(config) }

type Generator struct {
	renderer Renderer
	now      func() time.Time
}

func New(renderer Renderer) *Generator {
	if renderer == nil {
		renderer = RendererFunc(renderFoundation)
	}
	return &Generator{renderer: renderer, now: time.Now}
}

type Result struct {
	Config spec.Config
	Files  []string
	DryRun bool
}

func (g *Generator) Create(config spec.Config) (Result, error) {
	resolved, err := spec.Resolve(config)
	if err != nil {
		return Result{}, err
	}
	root, err := projectRoot(resolved)
	if err != nil {
		return Result{}, err
	}
	files, err := g.renderer.Render(resolved)
	if err != nil {
		return Result{}, fmt.Errorf("render project: %w", err)
	}
	paths, err := validateFiles(files)
	if err != nil {
		return Result{}, err
	}
	if _, err := os.Stat(filepath.Join(root, ManifestPath)); err == nil {
		return Result{}, fmt.Errorf("%s already contains a generated project", root)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Result{}, err
	}
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(path))); err == nil {
			return Result{}, fmt.Errorf("refusing to overwrite existing file %q", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return Result{}, err
		}
	}
	if resolved.DryRun {
		return Result{Config: resolved, Files: paths, DryRun: true}, nil
	}
	if err := writeRendered(root, files, paths); err != nil {
		return Result{}, err
	}
	if err := g.writeManifest(root, resolved, files, paths, nil); err != nil {
		return Result{}, err
	}
	return Result{Config: resolved, Files: paths}, nil
}

func (g *Generator) AddFeature(root string, feature spec.Feature, dryRun bool) (Result, error) {
	manifest, err := LoadManifest(root)
	if err != nil {
		return Result{}, err
	}
	if manifest.Config.Has(feature) {
		return Result{Config: manifest.Config, DryRun: dryRun}, nil
	}
	config := manifest.Config
	config.Features = append(config.Features, feature)
	config.DryRun = dryRun
	config.OutputDir = root
	config, err = spec.Resolve(config)
	if err != nil {
		return Result{}, err
	}
	files, err := g.renderer.Render(config)
	if err != nil {
		return Result{}, fmt.Errorf("render project: %w", err)
	}
	paths, err := validateFiles(files)
	if err != nil {
		return Result{}, err
	}
	owned := make(map[string]spec.ManagedFile, len(manifest.Files))
	for _, file := range manifest.Files {
		owned[file.Path] = file
	}
	customized := make(map[string]bool)
	for path, managed := range owned {
		current, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if readErr != nil {
			return Result{}, fmt.Errorf("read managed file %s: %w", path, readErr)
		}
		if managed.Customized || digest(current) != managed.SHA256 {
			customized[path] = true
			files[path] = current
		}
	}
	for _, path := range paths {
		if _, isOwned := owned[path]; isOwned {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(path))); err == nil {
			return Result{}, fmt.Errorf("refusing to overwrite user-owned file %q", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return Result{}, err
		}
	}
	if dryRun {
		return Result{Config: config, Files: paths, DryRun: true}, nil
	}
	if err := writeRendered(root, files, paths); err != nil {
		return Result{}, err
	}
	if err := g.writeManifest(root, config, files, paths, customized); err != nil {
		return Result{}, err
	}
	return Result{Config: config, Files: paths}, nil
}

func Verify(root string) error {
	manifest, err := LoadManifest(root)
	if err != nil {
		return err
	}
	if manifest.Version != spec.ManifestVersion {
		return fmt.Errorf("unsupported manifest version %d", manifest.Version)
	}
	if _, err := spec.Resolve(manifest.Config); err != nil {
		return fmt.Errorf("invalid manifest config: %w", err)
	}
	return verifyManaged(root, manifest)
}

func Doctor(root string) []error {
	var findings []error
	if _, err := os.Stat(root); err != nil {
		findings = append(findings, fmt.Errorf("project directory: %w", err))
		return findings
	}
	if err := Verify(root); err != nil {
		findings = append(findings, err)
	}
	return findings
}

func LoadManifest(root string) (spec.Manifest, error) {
	data, err := os.ReadFile(filepath.Join(root, ManifestPath))
	if err != nil {
		return spec.Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	var manifest spec.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return spec.Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}
	return manifest, nil
}

func (g *Generator) writeManifest(root string, config spec.Config, files map[string][]byte, paths []string, customized map[string]bool) error {
	managed := make([]spec.ManagedFile, 0, len(paths))
	for _, path := range paths {
		managed = append(managed, spec.ManagedFile{Path: path, SHA256: digest(files[path]), Customized: customized[path]})
	}
	config.OutputDir = ""
	config.DryRun = false
	manifest := spec.Manifest{Version: spec.ManifestVersion, Generator: GeneratorName, GeneratedAt: g.now().UTC(), Config: config, Files: managed}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeFile(root, ManifestPath, data)
}

func verifyManaged(root string, manifest spec.Manifest) error {
	for _, managed := range manifest.Files {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(managed.Path)))
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrConflict, managed.Path, err)
		}
		if digest(data) != managed.SHA256 {
			return fmt.Errorf("%w: %s", ErrConflict, managed.Path)
		}
	}
	return nil
}

func validateFiles(files map[string][]byte) ([]string, error) {
	paths := make([]string, 0, len(files))
	for path := range files {
		clean := filepath.ToSlash(filepath.Clean(path))
		if path == "" || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || filepath.IsAbs(path) || clean == ManifestPath {
			return nil, fmt.Errorf("unsafe generated path %q", path)
		}
		if clean != path {
			return nil, fmt.Errorf("generated path %q is not canonical", path)
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func writeRendered(root string, files map[string][]byte, paths []string) error {
	for _, path := range paths {
		if err := writeFile(root, path, files[path]); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(root, path string, data []byte) error {
	target := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func projectRoot(config spec.Config) (string, error) {
	root := config.OutputDir
	if root == "" {
		root = config.Name
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func renderFoundation(config spec.Config) (map[string][]byte, error) {
	features := make([]string, len(config.Features))
	for i, feature := range config.Features {
		features[i] = string(feature)
	}
	slices.Sort(features)
	readme := fmt.Sprintf("# %s\n\nGenerated by %s.\n\n- Module: `%s`\n- Database: `%s`\n- Features: `%s`\n", config.Name, GeneratorName, config.Module, config.Database, strings.Join(features, ", "))
	return map[string][]byte{"README.md": []byte(readme)}, nil
}
