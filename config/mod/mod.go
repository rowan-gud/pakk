package mod

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	ModKindBin ModKind = "bin"
	ModKindPkg ModKind = "pkg"
)

type ModKind string

type Mod struct {
	Name string  `toml:"name"`
	Kind ModKind `toml:"kind"`
}

func Parse(filePath string) (*Mod, error) {
	var mod Mod

	if _, err := toml.DecodeFile(filePath, &mod); err != nil {
		return nil, err
	}

	return &mod, nil
}

func (k *ModKind) UnmarshalText(text []byte) error {
	switch string(text) {
	case "bin", "binary":
		*k = ModKindBin
	case "pkg", "package":
		*k = ModKindPkg
	default:
		return fmt.Errorf("failed to parse %s as `ModKind`", text)
	}

	return nil
}
