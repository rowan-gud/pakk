package project

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/config/mod"
	"github.com/rowan-gud/pakk/config/render"
)

type Project struct {
	Name       string   `toml:"name"`
	RootImport []string `toml:"root_import,omitempty"`

	ctx     *render.ProjectContext
	modules []*mod.Mod
}

func Parse(bytes []byte, ctx *render.ProjectContext) (*Project, error) {
	var project Project

	if err := toml.Unmarshal(bytes, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	project.ctx = ctx

	return &project, nil
}

func (p *Project) AddMod(mod *mod.Mod) {
	p.modules = append(p.modules, mod)
}

func (p *Project) Ctx() *render.ProjectContext {
	return p.ctx
}

func (p *Project) Modules() []*mod.Mod {
	return p.modules
}
