package pkg

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

const (
	ProvidesKindImport ProvidesKind = "import"
)

type ProvidesKind string

type Pkg struct {
	Sources  []string        `toml:"sources"`
	Pre      []Pre           `toml:"pre"`
	Provides []ProvidesUnion `toml:"provides"`
}

type Pre struct {
	Sources   []string `toml:"sources"`
	RunPkg    []string `toml:"run_pkg"`
	RunFile   []string `toml:"run_file"`
	Generates []string `toml:"generates"`
}

type ProvidesUnion struct {
	Provides Provides `toml:"-"`
}

type Provides interface {
	Kind() ProvidesKind
}

type ProvidesImport struct {
	Path []string `toml:"path"`
}

func (p ProvidesImport) Kind() ProvidesKind {
	return ProvidesKindImport
}

func (p *ProvidesUnion) UnmarshalTOML(data any) error {
	m, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("expected type %T received %T", m, data)
	}

	kind, ok := m["kind"].(string)
	if !ok {
		return fmt.Errorf("kind expected type %T received %T", m, data)
	}

	marshalled, err := toml.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %w", err)
	}

	switch kind {
	case string(ProvidesKindImport):
		return toml.Unmarshal(marshalled, &p.Provides)
	}

	return fmt.Errorf("invalid kind %s", kind)
}

func (k *ProvidesKind) UnmarshalText(text []byte) error {
	switch string(text) {
	case "import":
		*k = ProvidesKindImport
	default:
		return fmt.Errorf("failed to parse %s as `ProvideKind`", text)
	}

	return nil
}
