package project

import "github.com/BurntSushi/toml"

type Project struct {
	Name string `toml:"name"`

	ctx ProjectContext `toml:"-"`
}

func Parse(filePath string, outDir string) (*Project, error) {
	var proj Project

	if _, err := toml.DecodeFile(filePath, &proj); err != nil {
		return nil, err
	}

	proj.ctx = ProjectContext{
		Name:   proj.Name,
		Path:   filePath,
		OutDir: outDir,
	}

	return &proj, nil
}
