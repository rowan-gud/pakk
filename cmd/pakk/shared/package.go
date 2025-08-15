package shared

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rowan-gud/pakk/config"
)

var (
	possiblePackageFiles = map[string]struct{}{
		"package.toml": {},
	}

	packageConfig *config.Package
)

func ParsePackageConfig(rootDir string) (*config.Package, error) {
	packageFile, err := findPackageFile(rootDir)
	if err != nil {
		return nil, err
	}

	cfg, err := config.ParsePackage(rootDir, packageFile)
	if err != nil {
		return nil, err
	}

	packageConfig = cfg

	return cfg, nil
}

func PackageConfig() *config.Package {
	return packageConfig
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
