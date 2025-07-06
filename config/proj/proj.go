package proj

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/config/mod"
	"github.com/rowan-gud/pakk/config/renderctx"
)

type Project struct {
	Name       string
	RootImport string

	Modules []*mod.Mod

	Ctx renderctx.ProjectContext
}

func Parse(filePath string, outDir string) (*Project, error) {
	if !filepath.IsAbs(filePath) {
		fp, err := filepath.Abs(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create absolute path: %w", err)
		}

		filePath = fp
	}

	var proj Project

	if _, err := toml.DecodeFile(filePath, &proj); err != nil {
		return nil, err
	}

	rootDir := filepath.Dir(filePath)

	proj.Ctx = renderctx.ProjectContext{
		Name:       proj.Name,
		Dir:        rootDir,
		OutDir:     outDir,
		RootImport: proj.RootImport,
	}

	return &proj, nil
}
