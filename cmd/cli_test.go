package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/melkeydev/go-blueprint/internal/generator"
)

func testCLI() (*CLI, *bytes.Buffer, *bytes.Buffer) {
	out, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	return &CLI{Out: out, Err: stderr, Gen: generator.New(nil)}, out, stderr
}

func TestCreateFlagsAndVerify(t *testing.T) {
	cli, out, _ := testCLI()
	root := filepath.Join(t.TempDir(), "project")
	err := cli.Run([]string{"create", "--name", "acme", "--module", "example.com/acme", "--database", "mysql", "--seed", "demo", "--no-frontend", "--feature", "cms", "--output", root})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "create: wrote 1 file(s)") {
		t.Fatalf("output = %q", out.String())
	}
	out.Reset()
	if err := cli.Run([]string{"verify", "--dir", root}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "verification passed") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestAddFeatureDryRun(t *testing.T) {
	cli, out, _ := testCLI()
	root := filepath.Join(t.TempDir(), "project")
	if err := cli.Run([]string{"create", "--name", "acme", "--module", "example.com/acme", "--feature", "jobs", "--output", root}); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := cli.Run([]string{"add", "feature", "cms", "--dir", root, "--dry-run"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "would write") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestCreateRejectsUnknownFeature(t *testing.T) {
	cli, _, _ := testCLI()
	err := cli.Run([]string{"create", "--name", "acme", "--module", "example.com/acme", "--feature", "wat", "--output", t.TempDir()})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInteractiveCreateUsesDefaults(t *testing.T) {
	cli, out, _ := testCLI()
	root := filepath.Join(t.TempDir(), "interactive")
	cli.In = strings.NewReader("acme\nexample.com/acme\n\n\n\n\n" + root + "\n")
	if err := cli.Run([]string{"create"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "create: wrote") {
		t.Fatalf("output = %q", out.String())
	}
}
