package config

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/utils"
)

type Package struct {
	Plugins []Plugin `toml:"plugins"`

	path *utils.PathResolver
}

type Plugin struct {
	Name    string         `toml:"name"`
	Version string         `toml:"version"`
	Options map[string]any `toml:"options"`
}

func ParsePackage(rootDir string, packageFile string) (*Package, error) {
	packageDir := filepath.Dir(packageFile)

	var j map[string]any
	if _, err := toml.DecodeFile(packageFile, &j); err != nil {
		return nil, fmt.Errorf("failed to decode package file: %w", err)
	}

	fmt.Printf("Package: %+v\n", j)

	var p Package
	if _, err := toml.DecodeFile(packageFile, &p); err != nil {
		return nil, fmt.Errorf("failed to decode package file: %w", err)
	}

	p.path = utils.NewPathResolver(rootDir, packageDir)

	return &p, nil
}
