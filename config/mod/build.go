package mod

import (
	"errors"
	"fmt"
)

func (m *Mod) Build() error {
	if m.Bin != nil {
		if err := m.Bin.Build(m); err != nil {
			return err
		}
	}

	if m.Pkg != nil {
		if err := m.Pkg.Build(m); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mod) error(wraps error, format string, vals ...any) error {
	s := fmt.Sprintf(format, vals...)

	if wraps != nil {
		m.logger.Error(s, "error", wraps)
		return fmt.Errorf("%s: %w", s, wraps)
	}

	m.logger.Error(s)
	return errors.New(s)
}

func (b *Bin) Build(mod *Mod) error {
	if err := b.Cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (p *Pkg) Build(mod *Mod) error {
	for _, pre := range p.Pre {
		if pre.Each == nil {
			if err := pre.Cmd.Run(); err != nil {
				return mod.error(err, "failed to run command %s", pre.Cmd.String())
			}
		} else {
			if err := pre.Cmd.RunEach(pre.Each.Expanded()); err != nil {
				return mod.error(err, "failed to run command %s", pre.Cmd.String())
			}
		}
	}

	return nil
}
