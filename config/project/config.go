package project

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/collections"
	"github.com/rowan-gud/pakk/config/mod"
	"github.com/rowan-gud/pakk/config/render"
)

type Project struct {
	Name       string   `toml:"name"`
	RootImport []string `toml:"root_import,omitempty"`

	ctx      *render.ProjectContext
	lockFile *ProjectLock
	logger   *slog.Logger
	modules  *collections.DirectedGraph[string, *mod.Mod]
}

type ProjectLock struct {
	Modules map[string]string `toml:"modules"`
}

func Parse(bytes []byte, ctx *render.ProjectContext, logFile *os.File) (*Project, error) {
	logger := slog.New(slog.NewTextHandler(logFile, nil))
	var project Project

	if err := toml.Unmarshal(bytes, &project); err != nil {
		logger.Error("failed to parse project",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	lockFilePath := filepath.Join(ctx.Out, "project.lock.toml")

	var lockFile ProjectLock
	_, err := toml.DecodeFile(lockFilePath, &lockFile)
	if os.IsNotExist(err) {
		lockFile = ProjectLock{
			Modules: make(map[string]string),
		}
	} else if err != nil {
		logger.Error("failed to decode lock file",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to decode lock file: %w", err)
	}

	project.ctx = ctx
	project.lockFile = &lockFile
	project.logger = logger
	project.modules = collections.NewDirectedGraph[string, *mod.Mod]()

	return &project, nil
}

func (p *Project) AddMod(path string, mod *mod.Mod) {
	p.modules.AddNode(path, mod)
}

func (p *Project) Cleanup() {
	marshalled, err := toml.Marshal(p.lockFile)
	if err != nil {
		p.logger.Warn("failed to marshal lock file",
			slog.Any("error", err),
		)
	}

	lockFilePath := filepath.Join(p.ctx.Out, "project.lock.toml")
	if err := os.WriteFile(lockFilePath, marshalled, os.ModePerm); err != nil {
		p.logger.Warn("failed to write lock file",
			slog.Any("error", err),
		)
	}
}

func (p *Project) Ctx() *render.ProjectContext {
	return p.ctx
}

func (p *Project) Modules() map[string]*mod.Mod {
	nodes := p.modules.Nodes()
	res := make(map[string]*mod.Mod, len(nodes))

	for path, n := range nodes {
		res[path] = n.Data()
	}

	return res
}

func (p *Project) ModulesGraph() *collections.DirectedGraph[string, *mod.Mod] {
	return p.modules
}

func (p *Project) InitializeDependencies() {
	for key, n := range p.modules.Nodes() {
		for _, dep := range n.Data().Deps {
			p.modules.AddEdge(key, dep)
		}
	}
}
