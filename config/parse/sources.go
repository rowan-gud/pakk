package parse

import "github.com/BurntSushi/toml"

type Sources struct {
	parsed []string
	raw    any
}

func (s Sources) MarshalTOML() ([]byte, error) {
	return toml.Marshal(s.raw)
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
