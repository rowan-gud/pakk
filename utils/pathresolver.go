package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

type PathResolver struct {
	rootDir  string
	buildDir string
}

func NewPathResolver(rootDir string, buildDir string) *PathResolver {
	return &PathResolver{
		rootDir:  rootDir,
		buildDir: buildDir,
	}
}

func (r *PathResolver) Resolve(path string) ([]string, error) {
	var absPath string

	if trimmed, hadPrefix := strings.CutPrefix(path, "//"); hadPrefix {
		absPath = filepath.Join(r.rootDir, trimmed)
	} else {
		absPath = filepath.Join(r.buildDir, path)
	}

	rel, err := filepath.Rel(r.rootDir, absPath)
	if err != nil {
		return nil, fmt.Errorf("could not get relative path: %w", err)
	}

	parts := filepath.SplitList(rel)

	if parts[0] == "." {
		parts = parts[1:]
	}

	parts = append([]string{"/"}, parts...)

	return parts, nil
}
