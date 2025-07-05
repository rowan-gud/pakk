package pkg

import (
	"fmt"
)

const (
	ProvideKindImport ProvideKind = "import"
)

type ProvideKind string

type Pkg struct {
	Sources  []string   `toml:"sources"`
	Pre      []Pre      `toml:"pre"`
	Provides []Provides `toml:"provides"`
}

type Pre struct {
	Sources   []string `toml:"sources"`
	RunPkg    []string `toml:"run_pkg"`
	RunFile   []string `toml:"run_file"`
	Generates []string `toml:"generates"`
}

type ProvidesUnion struct {
	Kind ProvideKind `toml:"kind"`
}

type Provides interface {
	Kind() ProvideKind
}

type ProvidesImport struct {
	Path []string `toml:"path"`
}

func (p ProvidesImport) Kind() ProvideKind {
	return ProvideKindImport
}

func (k *ProvideKind) UnmarshalText(text []byte) error {
	switch string(text) {
	case "import":
		*k = ProvideKindImport
	default:
		return fmt.Errorf("failed to parse %s as `ProvideKind`", text)
	}

	return nil
}
