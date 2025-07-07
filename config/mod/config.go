package mod

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/config/parse"
	"github.com/rowan-gud/pakk/config/render"
)

type Mod struct {
	Name string   `toml:"name"`
	Deps []string `toml:"deps,omitempty"`

	*Bin `toml:"bin,omitempty"`
	*Pkg `toml:"pkg,omitempty"`

	ctx *render.RenderContext
}

type Bin struct {
	Artifacts []string      `toml:"artifacts"`
	Cmd       parse.Command `toml:"cmd"`
	Sources   parse.Sources `toml:"sources"`
}

type Pkg struct {
	Pre      []PkgPre     `toml:"pre,omitempty"`
	Provides *PkgProvides `toml:"provides,omitempty"`
}

type PkgPre struct {
	Each      *parse.Sources `toml:"each,omitempty"`
	Cmd       parse.Command  `toml:"cmd"`
	Generates *parse.Sources `toml:"generates,omitempty"`
}

type PkgProvides struct {
	Import *PkgProvidesImport `toml:"import,omitempty"`
}

type PkgProvidesImport struct {
	Path []string `toml:"path"`
}

func Parse(bytes []byte, ctx *render.RenderContext) (*Mod, error) {
	var mod Mod

	if err := toml.Unmarshal(bytes, &mod); err != nil {
		return nil, fmt.Errorf("failed to parse module: %w", err)
	}

	mod.ctx = ctx

	return &mod, nil
}
