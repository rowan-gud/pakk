package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rowan-gud/pakk/config"
)

var (
	possibleBuildFiles = map[string]struct{}{
		"build.toml": {},
	}
	possiblePackageFiles = map[string]struct{}{
		"package.toml": {},
	}
)

func parsePackageConfig(rootDir string) (*config.Package, error) {
	packageFile, err := findPackageFile(rootDir)
	if err != nil {
		return nil, err
	}

	packageConfig, err := config.ParsePackage(rootDir, packageFile)
	if err != nil {
		return nil, err
	}

	return packageConfig, nil
}

func parseBuildConfigs(rootDir string) (map[string]*config.Build, error) {
	res := make(map[string]*config.Build)

	if err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		filePath, err := findBuildFile(path)
		if err != nil || filePath == "" {
			return err
		}

		build, err := config.ParseBuild(rootDir, filePath)
		if err != nil {
			return err
		}

		res[path] = build

		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

func findPackageFile(rootDir string) (string, error) {
	dirFiles, err := os.ReadDir(rootDir)
	if err != nil {
		return "", fmt.Errorf("failed to read dir %s: %w", rootDir, err)
	}

	for _, dirFile := range dirFiles {
		if dirFile.IsDir() {
			continue
		}

		if _, ok := possiblePackageFiles[dirFile.Name()]; ok {
			return filepath.Join(rootDir, dirFile.Name()), nil
		}
	}

	return "", nil
}

func findBuildFile(rootDir string) (string, error) {
	dirFiles, err := os.ReadDir(rootDir)
	if err != nil {
		return "", fmt.Errorf("failed to read dir %s: %w", rootDir, err)
	}

	for _, dirFile := range dirFiles {
		if dirFile.IsDir() {
			continue
		}

		if _, ok := possibleBuildFiles[dirFile.Name()]; ok {
			return filepath.Join(rootDir, dirFile.Name()), nil
		}
	}

	return "", nil
}
