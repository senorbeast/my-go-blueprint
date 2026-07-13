package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/melkeydev/go-blueprint/internal/generator"
	"github.com/melkeydev/go-blueprint/internal/spec"
	"github.com/melkeydev/go-blueprint/internal/templates"
)

type CLI struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
	Gen *generator.Generator
}

func New() *CLI {
	return &CLI{In: os.Stdin, Out: os.Stdout, Err: os.Stderr, Gen: generator.New(templates.Renderer{})}
}

func (cli *CLI) Run(args []string) error {
	if len(args) == 0 {
		return cli.usageError("command is required")
	}
	switch args[0] {
	case "create":
		return cli.create(args[1:])
	case "add":
		return cli.add(args[1:])
	case "doctor":
		return cli.doctor(args[1:])
	case "verify":
		return cli.verify(args[1:])
	case "help", "-h", "--help":
		cli.printUsage()
		return nil
	default:
		return cli.usageError(fmt.Sprintf("unknown command %q", args[0]))
	}
}

type featureFlags []spec.Feature

func (f *featureFlags) String() string {
	values := make([]string, len(*f))
	for i, v := range *f {
		values[i] = string(v)
	}
	return strings.Join(values, ",")
}
func (f *featureFlags) Set(value string) error { *f = append(*f, spec.Feature(value)); return nil }

func (cli *CLI) create(args []string) error {
	if len(args) == 0 {
		return cli.createInteractive()
	}
	defaults := spec.DefaultConfig()
	flags := flag.NewFlagSet("create", flag.ContinueOnError)
	flags.SetOutput(cli.Err)
	name := flags.String("name", "", "project name")
	module := flags.String("module", "", "Go module path")
	database := flags.String("database", string(defaults.Database), "postgres or mysql")
	seed := flags.String("seed", string(defaults.Seed), "none, minimal, or demo")
	frontend := flags.Bool("frontend", defaults.Frontend, "include the React frontend")
	noFrontend := flags.Bool("no-frontend", false, "omit the React frontend")
	docker := flags.Bool("docker", defaults.Docker, "include Docker configuration")
	noDocker := flags.Bool("no-docker", false, "omit Docker configuration")
	output := flags.String("output", "", "output directory")
	dryRun := flags.Bool("dry-run", false, "print planned files without writing")
	features := featureFlags{}
	flags.Var(&features, "feature", "feature to include (repeatable)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return cli.usageError("create does not accept positional arguments")
	}
	if *noFrontend {
		*frontend = false
	}
	if *noDocker {
		*docker = false
	}
	if len(features) == 0 {
		features = append(features, defaults.Features...)
	}
	config := spec.Config{Name: *name, Module: *module, Database: spec.Database(*database), Seed: spec.SeedProfile(*seed), Frontend: *frontend, Docker: *docker, Features: features, OutputDir: *output, DryRun: *dryRun}
	result, err := cli.Gen.Create(config)
	if err != nil {
		return err
	}
	return cli.printResult("create", result)
}

func (cli *CLI) createInteractive() error {
	defaults := spec.DefaultConfig()
	input := cli.In
	if input == nil {
		input = os.Stdin
	}
	reader := bufio.NewReader(input)
	read := func(label, fallback string) (string, error) {
		if fallback == "" {
			fmt.Fprintf(cli.Out, "%s: ", label)
		} else {
			fmt.Fprintf(cli.Out, "%s [%s]: ", label, fallback)
		}
		value, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}
		value = strings.TrimSpace(value)
		if value == "" {
			return fallback, nil
		}
		return value, nil
	}
	name, err := read("Project name", "")
	if err != nil {
		return err
	}
	module, err := read("Go module", "example.com/"+name)
	if err != nil {
		return err
	}
	database, err := read("Database (postgres/mysql)", string(defaults.Database))
	if err != nil {
		return err
	}
	seed, err := read("Seed profile (none/minimal/demo)", string(defaults.Seed))
	if err != nil {
		return err
	}
	featureDefault := make([]string, len(defaults.Features))
	for index, feature := range defaults.Features {
		featureDefault[index] = string(feature)
	}
	featureText, err := read("Features (comma-separated)", strings.Join(featureDefault, ","))
	if err != nil {
		return err
	}
	frontendText, err := read("Include React frontend (yes/no)", "yes")
	if err != nil {
		return err
	}
	output, err := read("Output directory", name)
	if err != nil {
		return err
	}
	features := make([]spec.Feature, 0)
	for _, feature := range strings.Split(featureText, ",") {
		if feature = strings.TrimSpace(feature); feature != "" {
			features = append(features, spec.Feature(feature))
		}
	}
	frontend := strings.EqualFold(frontendText, "yes") || strings.EqualFold(frontendText, "y") || strings.EqualFold(frontendText, "true")
	result, err := cli.Gen.Create(spec.Config{
		Name: name, Module: module, Database: spec.Database(database), Seed: spec.SeedProfile(seed),
		Frontend: frontend, Docker: true, Features: features, OutputDir: output,
	})
	if err != nil {
		return err
	}
	return cli.printResult("create", result)
}

func (cli *CLI) add(args []string) error {
	if len(args) < 2 || args[0] != "feature" {
		return cli.usageError("usage: blueprint add feature <name> [--dir path] [--dry-run]")
	}
	flags := flag.NewFlagSet("add feature", flag.ContinueOnError)
	flags.SetOutput(cli.Err)
	root := flags.String("dir", ".", "generated project directory")
	dryRun := flags.Bool("dry-run", false, "print planned files without writing")
	if err := flags.Parse(args[2:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return cli.usageError("add feature accepts one feature name before its options")
	}
	result, err := cli.Gen.AddFeature(*root, spec.Feature(args[1]), *dryRun)
	if err != nil {
		return err
	}
	return cli.printResult("add feature", result)
}

func (cli *CLI) verify(args []string) error {
	flags := flag.NewFlagSet("verify", flag.ContinueOnError)
	flags.SetOutput(cli.Err)
	root := flags.String("dir", ".", "generated project directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return cli.usageError("verify does not accept positional arguments")
	}
	if err := generator.Verify(*root); err != nil {
		return err
	}
	_, err := fmt.Fprintln(cli.Out, "verification passed")
	return err
}

func (cli *CLI) doctor(args []string) error {
	flags := flag.NewFlagSet("doctor", flag.ContinueOnError)
	flags.SetOutput(cli.Err)
	root := flags.String("dir", ".", "generated project directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return cli.usageError("doctor does not accept positional arguments")
	}
	findings := generator.Doctor(*root)
	if len(findings) == 0 {
		_, err := fmt.Fprintln(cli.Out, "doctor found no problems")
		return err
	}
	for _, finding := range findings {
		fmt.Fprintf(cli.Err, "- %v\n", finding)
	}
	return fmt.Errorf("doctor found %d problem(s)", len(findings))
}

func (cli *CLI) printResult(action string, result generator.Result) error {
	mode := "wrote"
	if result.DryRun {
		mode = "would write"
	}
	if _, err := fmt.Fprintf(cli.Out, "%s: %s %d file(s)\n", action, mode, len(result.Files)); err != nil {
		return err
	}
	for _, path := range result.Files {
		if _, err := fmt.Fprintf(cli.Out, "  %s\n", path); err != nil {
			return err
		}
	}
	return nil
}

func (cli *CLI) usageError(message string) error { cli.printUsage(); return errors.New(message) }
func (cli *CLI) printUsage() {
	fmt.Fprintln(cli.Err, "usage: blueprint <create|add feature|doctor|verify> [options]")
}
