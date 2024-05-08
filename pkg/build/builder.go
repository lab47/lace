package build

import (
	"context"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lab47/lace/pkg/pkgreflect"
	"github.com/lab47/lsvd/logger"
	"gopkg.in/yaml.v3"
)

type GoImport struct {
	Path string `yaml:"path"`
	As   string `yaml:"as"`
}

type Config struct {
	Name      string     `yaml:"name"`
	GoImports []GoImport `yaml:"go-imports"`
}

type Builder struct {
	log logger.Logger
	dir string
	cfg *Config
}

func LoadBuilder(log logger.Logger, dir string) (*Builder, error) {
	f, err := os.Open(filepath.Join(dir, "lace.yml"))
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &Builder{log: log, cfg: &cfg, dir: dir}, nil
}

var fileCleanup = strings.NewReplacer(
	"/", "_",
	".", "__",
	"~", "___",
	"\\", "____",
)

func (b *Builder) Run(ctx context.Context) error {
	b.log.Info("beginning build", "name", b.cfg.Name)

	dir := filepath.Join(b.dir, "_build_"+b.cfg.Name)
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return err
	}

	// defer os.RemoveAll(dir)

	for _, imp := range b.cfg.GoImports {
		name := imp.As
		if name == "" {
			name = imp.Path
		}

		dest := filepath.Join(dir, fileCleanup.Replace(imp.Path)) + ".go"

		b.log.Info("generating binding", "path", imp.Path, "name", name)
		err = pkgreflect.Generate(imp.Path, name, b.dir, dest, "main", &pkgreflect.Match{})
		if err != nil {
			return err
		}
	}

	b.log.Info("writing main.go")
	err = b.writeMain(dir)
	if err != nil {
		return err
	}

	b.log.Info("compiling")
	cmd := exec.CommandContext(ctx, "go", "build", "-o", filepath.Join(b.dir, b.cfg.Name), ".")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()

}

func (b *Builder) writeMain(dir string) error {
	data := []byte(fmt.Sprintf(`
package main

import "github.com/lab47/lace/cli"

var Program = "%s"

func main() {
    cli.Main()
}
    `, b.cfg.Name))

	data, err := format.Source(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "main.go"), data, 0644)
}
