package parse

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Sources struct {
	expanded []string
	parsed   []string
	raw      any
}

func (s *Sources) Expand(base string) ([]string, error) {
	if len(s.expanded) != 0 {
		return s.expanded, nil
	}

	for _, source := range s.parsed {
		if !filepath.IsAbs(source) && base != "" {
			source = filepath.Join(base, source)
		}

		if strings.Contains(source, "**") {
			return nil, errors.New("double star (`**`) glob pattern not supported")
		}

		if strings.Contains(source, "*") {
			glob, err := filepath.Glob(source)
			if err != nil {
				return nil, fmt.Errorf("failed to expand path: %w", err)
			}

			s.expanded = append(s.expanded, glob...)
		} else {
			s.expanded = append(s.expanded, source)
		}
	}

	return s.expanded, nil
}

func (s *Sources) Expanded() []string {
	return s.expanded
}

func (s Sources) MarshalTOML() ([]byte, error) {
	if len(s.expanded) > 0 {
		return toml.Marshal(s.expanded)
	}

	return toml.Marshal(s.raw)
}

func (s *Sources) WriteSum(h hash.Hash) (hash.Hash, error) {
	if h == nil {
		h = sha256.New()
	}

	for _, source := range s.expanded {
		f, err := os.ReadFile(source)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", source, err)
		}

		_, err = h.Write(f)
		if err != nil {
			return nil, fmt.Errorf("failed to write hash: %w", err)
		}
	}

	return h, nil
}

func (s *Sources) UnmarshalTOML(data any) error {
	if data == nil {
		return nilErr("Sources")
	}

	var sources []string
	var err error

	switch d := data.(type) {
	case string:
		sources = []string{d}
	case []string:
		sources = d
	case []any:
		sources, err = parseStringArrayFromAnyArray(d)
	default:
		return typeErr("Sources", d)
	}

	if err != nil {
		return err
	}

	s.raw = data
	s.parsed = sources

	return nil
}
