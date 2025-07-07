package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/collections"
	"github.com/rowan-gud/pakk/config/mod"
	"github.com/rowan-gud/pakk/config/project"
	"github.com/rowan-gud/pakk/config/render"
)

var defaultExcludeDirs = []string{
	".git",
}

type parseConfigOptions struct {
	RootDir     string
	ExcludeDirs []string
}

func parseValues(pakkDir string) (map[string]any, error) {
	valuesFile := filepath.Join(pakkDir, "values.toml")

	values := make(map[string]any)
	if _, err := toml.DecodeFile(valuesFile, &values); err != nil {
		return nil, err
	}

	return values, nil
}

func renderProjectFile(pakkDir string, ctx *render.ProjectContext) ([]byte, error) {
	helpersFile := filepath.Join(pakkDir, "helpers.tpl")
	projectFile := filepath.Join(pakkDir, "project.toml")

	tpl, err := template.ParseFiles(helpersFile, projectFile)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "project.toml", ctx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func renderModFile(pakkDir string, ctx *render.RenderContext) ([]byte, error) {
	helpersFile := filepath.Join(pakkDir, "helpers.tpl")
	projectFile := filepath.Join(pakkDir, "mod.toml")

	tpl, err := template.ParseFiles(helpersFile, projectFile)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "mod.toml", ctx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseProject(rootDir string) (*project.Project, error) {
	pakkDir := filepath.Join(rootDir, ".pakk")
	outDir := filepath.Join(pakkDir, "out")

	values, err := parseValues(pakkDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project values: %w", err)
	}

	ctx := &render.ProjectContext{
		Out:    outDir,
		Path:   rootDir,
		Values: values,
	}

	projectBytes, err := renderProjectFile(pakkDir, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to render project file: %w", err)
	}

	project, err := project.Parse(projectBytes, ctx)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func parseConfig(opts *parseConfigOptions) (*project.Project, error) {
	if len(opts.ExcludeDirs) == 0 {
		opts.ExcludeDirs = defaultExcludeDirs
	}

	project, err := parseProject(opts.RootDir)
	if err != nil {
		return nil, err
	}

	if err := findModules(opts.RootDir, opts, project); err != nil {
		return nil, err
	}

	return project, nil
}

func findModules(rootDir string, opts *parseConfigOptions, project *project.Project) error {
	excludeDirs := make([]string, len(opts.ExcludeDirs))

	for idx, dir := range opts.ExcludeDirs {
		if !strings.HasSuffix(dir, "/") {
			excludeDirs[idx] = dir + "/"
		} else {
			excludeDirs[idx] = dir
		}
	}

	pakkDirs, err := getMatchingPaths(rootDir, func(path string, d fs.DirEntry) bool {
		for _, dir := range opts.ExcludeDirs {
			if isExcluded(path, dir) {
				return false
			}
		}

		return d.IsDir() && d.Name() == ".pakk"
	})
	if err != nil {
		return err
	}

	// We have at most this many mod files
	modDirs := make([]string, 0, len(pakkDirs))

	for _, dir := range pakkDirs {
		modFile := filepath.Join(dir, "mod.toml")
		if _, err := os.Stat(modFile); err == nil {
			modDirs = append(modDirs, filepath.Dir(dir))
		}
	}

	for _, dir := range modDirs {
		pakkDir := filepath.Join(dir, ".pakk")

		values, err := parseValues(pakkDir)
		if err != nil {
			return err
		}

		ctx := &render.RenderContext{
			Mod: render.ModContext{
				Path: dir,
			},
			Project: project.Ctx(),
			Values:  values,
		}

		modBytes, err := renderModFile(pakkDir, ctx)
		if err != nil {
			return err
		}

		module, err := mod.Parse(modBytes, ctx)
		if err != nil {
			return err
		}

		project.AddMod(module)
	}

	return nil
}

func getMatchingPaths(rootDir string, fn func(path string, d fs.DirEntry) bool) ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if fn(path, d) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func isExcluded(path string, dir string) bool {
	pathList := filepath.SplitList(path)
	dirList := filepath.SplitList(dir)

	cleanedPathList := collections.NewStack[string]()
	cleanedDirList := collections.NewStack[string]()

	for _, s := range pathList {
		if s == ".." && !cleanedPathList.Empty() {
			cleanedPathList.Pop()
		}
		if s != "." {
			cleanedPathList.Push(s)
		}
	}

	for _, s := range dirList {
		if s == ".." && !cleanedDirList.Empty() {
			cleanedDirList.Pop()
		}
		if s != "." {
			cleanedDirList.Push(s)
		}
	}

	pathList = cleanedPathList.ToList()
	dirList = cleanedDirList.ToList()

	for idx, s := range dirList {
		if s != "*" && s != pathList[idx] {
			return false
		}
	}

	return true
}
