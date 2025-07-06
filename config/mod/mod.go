package mod

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/config/parser"
	"github.com/rowan-gud/pakk/config/renderctx"
)

type ModKind string

type Mod struct {
	Name string

	Bin *Bin
	Pkg *Pkg
	raw map[string]any
}

func Parse(filePath string, projCtx *renderctx.ProjectContext) (*Mod, error) {
	var mod Mod

	if _, err := toml.DecodeFile(filePath, &mod); err != nil {
		return nil, err
	}

	dir := filepath.Dir(filePath)

	modCtx := &renderctx.ModContext{
		Name: mod.Name,
		Path: filePath,
		Dir:  dir,
	}

	ctx := renderctx.New(projCtx, modCtx)

	if mod.raw["bin"] != nil {
		bin, err := parseBin(mod.raw["bin"], ctx)
		if err != nil {
			return nil, err
		}

		mod.Bin = bin
	}

	if mod.raw["pkg"] != nil {
		pkg, err := parsePkg(mod.raw["pkg"], ctx)
		if err != nil {
			return nil, err
		}

		mod.Pkg = pkg
	}

	return &mod, nil
}

func (m *Mod) UnmarshalTOML(data any) error {
	root, err := parser.ParseMap(data)
	if err != nil {
		return fmt.Errorf("failed to parse `Mod`: %w", err)
	}

	name, err := parser.ParseString(root["name"])
	if err != nil {
		return fmt.Errorf("failed to parse `Mod.name`: %w", err)
	}

	m.Name = name
	m.raw = root

	return nil
}
