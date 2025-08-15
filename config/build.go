package config

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/rowan-gud/pakk/utils"
)

type Build struct {
	Bin map[string]Bin `toml:"bin"`
	Lib map[string]Lib `toml:"lib"`

	path *utils.PathResolver
}

func ParseBuild(rootDir string, buildFile string) (*Build, error) {
	buildDir := filepath.Dir(buildFile)

	var b Build
	if _, err := toml.DecodeFile(buildFile, &b); err != nil {
		return nil, fmt.Errorf("failed to decode build file: %w", err)
	}

	b.path = utils.NewPathResolver(rootDir, buildDir)

	return &b, nil
}

func (b *Build) Resolve(path string) ([]string, error) {
	return b.path.Resolve(path)
}
