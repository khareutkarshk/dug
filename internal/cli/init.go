package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/khareutkarshk/dug/internal/cli/templates"
)

const defaultGHCRImage = "ghcr.io/khareutkarshk/dug:latest"

// scaffoldFile describes one file written by `dug init`.
type scaffoldFile struct {
	// Source is the path inside templates.FS.
	Source string
	// Dest is the path relative to the project root.
	Dest string
	// Templated selects whether the destination is rendered as a text/template.
	Templated bool
}

var initScaffold = []scaffoldFile{
	{Source: "edge.yaml", Dest: filepath.Join("configs", "edge.yaml")},
	{Source: "docker-compose.yml", Dest: "docker-compose.yml", Templated: true},
	{Source: "gitignore", Dest: ".gitignore"},
	{Source: "gitkeep", Dest: filepath.Join("certs", ".gitkeep")},
	{Source: "README.md", Dest: "README.md", Templated: true},
}

type initTemplateData struct {
	Name  string
	Image string
}

// Init bootstraps a new DUG project directory.
//
// Usage:
//
//	dug init <directory> [--force]
func Init(args []string) error {
	return initProject(os.Stdout, os.Stderr, args)
}

func initProject(stdout, stderr io.Writer, args []string) error {
	out := newPrinter(stdout)
	errOut := newPrinter(stderr)

	target, force, err := parseInitArgs(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			errOut.println("Usage: dug init <directory> [--force]")
			errOut.println()
			errOut.println("Options:")
			errOut.println("  --force   Overwrite existing files")
			if errOut.err != nil {
				return reported(errOut.err)
			}
			return reported(err)
		}
		errOut.println("Usage: dug init <directory> [--force]")
		if errOut.err != nil {
			return reported(errOut.err)
		}
		return reported(err)
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		return reported(fmt.Errorf("resolve path: %w", err))
	}

	data := initTemplateData{
		Name:  filepath.Base(absTarget),
		Image: defaultGHCRImage,
	}

	if err := os.MkdirAll(absTarget, 0o755); err != nil {
		return reported(fmt.Errorf("create project directory: %w", err))
	}

	created := make([]string, 0, len(initScaffold))
	skipped := make([]string, 0)

	for _, file := range initScaffold {
		destPath := filepath.Join(absTarget, file.Dest)

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return reported(fmt.Errorf("create directory for %s: %w", file.Dest, err))
		}

		_, statErr := os.Stat(destPath)
		exists := statErr == nil
		if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
			return reported(fmt.Errorf("stat %s: %w", file.Dest, statErr))
		}

		if exists && !force {
			skipped = append(skipped, file.Dest)
			continue
		}

		content, err := renderScaffoldFile(file, data)
		if err != nil {
			return reported(fmt.Errorf("render %s: %w", file.Source, err))
		}

		if err := os.WriteFile(destPath, content, 0o644); err != nil {
			return reported(fmt.Errorf("write %s: %w", file.Dest, err))
		}

		created = append(created, file.Dest)
	}

	if len(created) == 0 && len(skipped) > 0 {
		errOut.printf("✗ nothing written: all files already exist in %s\n", target)
		errOut.println("  re-run with --force to overwrite")
		if errOut.err != nil {
			return reported(errOut.err)
		}
		return reported(errors.New("project already initialized"))
	}

	out.printf("✓ Created DUG project %q\n", data.Name)
	out.println()
	if len(created) > 0 {
		out.println("Generated:")
		for _, path := range created {
			out.printf("  + %s\n", path)
		}
	}
	if len(skipped) > 0 {
		out.println()
		out.println("Skipped (already exists):")
		for _, path := range skipped {
			out.printf("  · %s\n", path)
		}
	}
	out.println()
	out.println("Next steps:")
	out.printf("  cd %s\n", target)
	out.println("  dug run -config configs/edge.yaml")
	out.println()
	out.println("Or with Docker:")
	out.printf("  cd %s\n", target)
	out.println("  docker compose up")

	if out.err != nil {
		return reported(out.err)
	}
	return nil
}

func parseInitArgs(args []string) (dir string, force bool, err error) {
	var positional []string

	for _, arg := range args {
		switch arg {
		case "-force", "--force":
			force = true
		case "-h", "-help", "--help":
			return "", false, flag.ErrHelp
		default:
			if strings.HasPrefix(arg, "-") {
				return "", false, fmt.Errorf("unknown flag: %s", arg)
			}
			positional = append(positional, arg)
		}
	}

	if len(positional) != 1 {
		return "", false, errors.New("directory argument is required")
	}

	dir = strings.TrimSpace(positional[0])
	if dir == "" {
		return "", false, errors.New("directory must be a non-empty project path")
	}

	return dir, force, nil
}

func renderScaffoldFile(file scaffoldFile, data initTemplateData) ([]byte, error) {
	raw, err := templates.FS.ReadFile(file.Source)
	if err != nil {
		return nil, err
	}

	if !file.Templated {
		return raw, nil
	}

	tmpl, err := template.New(file.Source).Parse(string(raw))
	if err != nil {
		return nil, err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}
