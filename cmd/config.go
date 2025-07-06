package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rowan-gud/pakk/config/mod"
	"github.com/rowan-gud/pakk/config/proj"
	"github.com/rowan-gud/pakk/config/renderctx"
)

var defaultExcludeDirs = []string{
	".git",
}

type parseConfigOptions struct {
	RootDir     string
	OutDir      string
	ExcludeDirs []string
}

func parseConfig(opts *parseConfigOptions) (*proj.Project, error) {
	if len(opts.ExcludeDirs) == 0 {
		opts.ExcludeDirs = defaultExcludeDirs
	}

	projectFile, err := getProjFileName(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get project file name: %w", err)
	}

	project, err := proj.Parse(projectFile, opts.OutDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project file: %w", err)
	}

	modules, err := findModules(opts.RootDir, opts, &project.Ctx)
	if err != nil {
		return nil, err
	}

	project.Modules = modules

	return project, nil
}

func findModules(rootDir string, opts *parseConfigOptions, projCtx *renderctx.ProjectContext) ([]*mod.Mod, error) {
	abs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	excludeDirs := make([]string, len(opts.ExcludeDirs))

	for idx, dir := range opts.ExcludeDirs {
		if !strings.HasSuffix(dir, "/") {
			excludeDirs[idx] = dir + "/"
		} else {
			excludeDirs[idx] = dir
		}
	}

	modFiles, err := getMatchingFiles(abs, func(path string) bool {
		for _, dir := range opts.ExcludeDirs {
			if strings.Contains(path, dir) {
				return false
			}
		}

		return isModFile(path)
	})
	if err != nil {
		return nil, err
	}

	modules := make([]*mod.Mod, len(modFiles))

	for idx, file := range modFiles {
		module, err := mod.Parse(file, projCtx)
		if err != nil {
			return nil, fmt.Errorf("module %s: %w", file, err)
		}

		modules[idx] = module
	}

	return modules, nil
}

func getMatchingFiles(rootDir string, fn func(path string) bool) ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && fn(path) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func getProjFileName(opts *parseConfigOptions) (string, error) {
	if opts.RootDir != "" {
		entries, err := os.ReadDir(opts.RootDir)
		if err != nil {
			return "", err
		}

		for _, file := range entries {
			if !file.IsDir() && isProjFile(file.Name()) {
				return filepath.Join(opts.RootDir, file.Name()), nil
			}
		}

		return "", errors.New("no project file found")
	}

	return "", errors.New("no project file provided")
}

func isProjFile(path string) bool {
	name := filepath.Base(path)

	return name == "pakk.proj.toml"
}

func isModFile(path string) bool {
	name := filepath.Base(path)

	return name == "pakk.mod.toml"
}
