package build

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/lab47/lablog/logger"
	"github.com/lab47/lace/pkg/pkgreflect"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/blake2b"
	"gopkg.in/yaml.v3"
)

type GoImport struct {
	Path string `yaml:"path"`
	As   string `yaml:"as"`
}

type Config struct {
	Name      string     `yaml:"name"`
	Main      string     `yaml:"main"`
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

func (b *Builder) buildId() string {
	h, _ := blake2b.New256(nil)
	err := cbor.NewEncoder(h).Encode(b.cfg)
	if err != nil {
		panic(err)
	}

	return base58.Encode(h.Sum(nil))
}

func (b *Builder) Run(ctx context.Context) (string, error) {
	id := b.buildId()

	exePath := filepath.Join(b.dir, "artifacts", b.cfg.Name)

	err := os.MkdirAll(filepath.Dir(exePath), 0755)
	if err != nil {
		return "", err
	}

	lastIdPath := exePath + ".build-id.txt"
	if data, err := os.ReadFile(lastIdPath); err == nil {
		if strings.TrimSpace(string(data)) == id {
			return exePath, nil
		}
	}

	b.log.Info("beginning build", "name", b.cfg.Name, "id", id)

	dir := filepath.Join(b.dir, "_build_"+b.cfg.Name)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(dir)

	for _, imp := range b.cfg.GoImports {
		name := imp.As
		if name == "" {
			name = imp.Path
		}

		dest := filepath.Join(dir, fileCleanup.Replace(imp.Path)) + ".go"

		b.log.Info("generating binding", "path", imp.Path, "name", name)
		err = pkgreflect.Generate(imp.Path, name, b.dir, dest, "main", &pkgreflect.Match{}, pkgreflect.GenOptions{})
		if err != nil {
			return "", err
		}
	}

	b.log.Info("writing main.go")
	err = b.writeMain(dir, id)
	if err != nil {
		return "", err
	}

	b.log.Info("compiling")
	cmd := exec.CommandContext(ctx, "go", "build", "-o", exePath, ".")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	os.WriteFile(exePath+".build-id.txt", []byte(id), 0644)

	return exePath, nil
}

func (b *Builder) writeMain(dir, id string) error {
	var data []byte

	if b.cfg.Main == "" {
		data = []byte(fmt.Sprintf(`
package main

import "github.com/lab47/lace/cli"

var Program = "%s"

func main() {
    cli.Main()
}
    `, b.cfg.Name))
	} else {
		data = []byte(fmt.Sprintf(`
package main

import "github.com/lab47/lace/cli"

var Program = "%s"
var BuildId = "%s"

func main() {
    cli.MainIn(%q)
}
    `, b.cfg.Name, id, b.cfg.Main))

	}

	data, err := format.Source(data)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "main.go")

	cur, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(data, cur) {
		return os.WriteFile(filepath.Join(dir, "main.go"), data, 0644)
	}

	return nil
}
