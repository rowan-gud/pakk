package mod

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/rowan-gud/pakk/config/parse"
	"github.com/rowan-gud/pakk/config/render"
)

const (
	DepKindInvalid DepKind = iota
	DepKindBin
	DepKindPkg
)

type DepKind int

type Mod struct {
	Name string   `toml:"name"`
	Deps []string `toml:"deps,omitempty"`

	*Bin `toml:"bin,omitempty"`
	*Pkg `toml:"pkg,omitempty"`

	ctx      *render.RenderContext
	buildDir string
	logger   *slog.Logger
}

type Dep struct {
	Kind DepKind `toml:"kind"`
	Path string  `toml:"path"`
}

type Bin struct {
	Artifacts []string      `toml:"artifacts"`
	Cmd       parse.Command `toml:"cmd"`
	Sources   parse.Sources `toml:"sources"`
}

type Pkg struct {
	Sources *parse.Sources `toml:"sources"`
	Pre     []PkgPre       `toml:"pre,omitempty"`
}

type PkgPre struct {
	Each      *parse.Sources `toml:"each,omitempty"`
	Cmd       parse.Command  `toml:"cmd"`
	Generates *parse.Sources `toml:"generates,omitempty"`
}

func Parse(bytes []byte, ctx *render.RenderContext, logFile *os.File) (*Mod, error) {
	var mod Mod

	if err := toml.Unmarshal(bytes, &mod); err != nil {
		return nil, fmt.Errorf("failed to parse module: %w", err)
	}

	mod.ctx = ctx

	if mod.Bin != nil {
		mod.Bin.Sources.Expand(ctx.Mod.Path)
	}

	if mod.Pkg != nil {
		if mod.Pkg.Sources != nil {
			mod.Pkg.Sources.Expand(ctx.Mod.Path)
		}

		for _, pre := range mod.Pre {
			if pre.Each != nil {
				pre.Each.Expand(ctx.Mod.Path)
			}
			if pre.Generates != nil {
				pre.Generates.Expand(ctx.Mod.Path)
			}
		}
	}

	mod.buildDir = filepath.Join(ctx.Project.Out, ctx.Mod.Path)

	_ = os.MkdirAll(mod.buildDir, 0755)

	mod.logger = slog.New(slog.NewTextHandler(logFile, nil)).With(
		slog.String("module", ctx.Mod.Path),
	)

	return &mod, nil
}

func (m *Mod) Path() string {
	return m.ctx.Mod.Path
}

func (m *Mod) Sum() (string, error) {
	h, err := m.WriteSum(nil)
	if err != nil {
		return "", err
	}

	sum := hex.EncodeToString(h.Sum([]byte{}))
	return sum, nil
}

func (m *Mod) WriteSum(h hash.Hash) (hash.Hash, error) {
	if h == nil {
		h = sha256.New()
	}

	var err error

	if m.Bin != nil {
		h, err = m.Bin.WriteSum(h)
		if err != nil {
			return nil, err
		}
	}
	if m.Pkg != nil {
		h, err = m.Pkg.WriteSum(h)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

func (b *Bin) WriteSum(h hash.Hash) (hash.Hash, error) {
	if h == nil {
		h = sha256.New()
	}

	var err error

	h, err = b.Sources.WriteSum(h)
	if err != nil {
		return nil, err
	}

	if len(b.Artifacts) > 0 {
		for _, artifact := range b.Artifacts {
			var bytes []byte

			info, err := os.Stat(artifact)
			if os.IsNotExist(err) {
				// If the artifact doesn't exist we should always rebuild so write the current timestamp
				bytes = []byte(time.Now().Format(time.StampMicro))
			} else if err != nil {
				// If it's a different error then fail
				return nil, err
			} else {
				// If no error we acknowledge that the artifact exists by writing the file name
				bytes = []byte(info.Name())
			}

			_, err = h.Write(bytes)
			if err != nil {
				return nil, err
			}
		}
	}

	return h, nil
}

func (p *Pkg) WriteSum(h hash.Hash) (hash.Hash, error) {
	if h == nil {
		h = sha256.New()
	}

	var err error

	if p.Sources != nil {
		h, err = p.Sources.WriteSum(h)
		if err != nil {
			return nil, err
		}
	}

	for _, pre := range p.Pre {
		h, err = pre.WriteSum(h)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

func (p *PkgPre) WriteSum(h hash.Hash) (hash.Hash, error) {
	if h == nil {
		h = sha256.New()
	}

	var err error

	if p.Each != nil {
		h, err = p.Each.WriteSum(h)
		if err != nil {
			return nil, err
		}
	}

	if p.Generates != nil {
		newH, err := p.Generates.WriteSum(h)
		if err == nil {
			h = newH
		}
	}

	return h, nil
}

func (k *DepKind) UnmarshalText(b []byte) error {
	switch string(b) {
	case "bin":
		*k = DepKindBin
	case "pkg":
		*k = DepKindPkg
	default:
		return fmt.Errorf("invalid value %s for type `DepKind`", b)
	}

	return nil
}

func (k DepKind) MarshalText() ([]byte, error) {
	switch k {
	case DepKindBin:
		return []byte("bin"), nil
	case DepKindPkg:
		return []byte("pkg"), nil
	default:
		return nil, fmt.Errorf("invalid value %d for type `DepKind`", int(k))
	}
}
