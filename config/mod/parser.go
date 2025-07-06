package mod

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/rowan-gud/pakk/config/parser"
	"github.com/rowan-gud/pakk/config/renderctx"
)

func parseRootError(parent string, wrap error) error {
	return fmt.Errorf("failed to parse `%s`: %w", parent, wrap)
}

func parseError(parent, field string, wrap error) error {
	return fmt.Errorf("failed to parse `%s.%s`: %w", parent, field, wrap)
}

func parseErrorNil(parent, field string) error {
	return parseError(parent, field, errors.New("unexpected nil value"))
}

func renderError(parent, field string, wrap error) error {
	return fmt.Errorf("failed to render `%s.%s`: %w", parent, field, wrap)
}

func parseRenderSources(data map[string]any, parent, field string, allowNil bool, ctx *renderctx.RenderCtx) ([]string, error) {
	if data[field] == nil && !allowNil {
		return nil, parseErrorNil(parent, field)
	}

	if data[field] == nil {
		return []string{}, nil
	}

	sources, err := parser.ParseStringArray(data[field])
	if err != nil {
		return nil, parseError(parent, field, err)
	}

	renderedSources, err := ctx.RenderSources(sources)
	if err != nil {
		return nil, renderError(parent, field, err)
	}

	return renderedSources, err
}

func parseRenderCommand(data map[string]any, parent, field string, allowNil bool, ctx *renderctx.RenderCtx) (*exec.Cmd, error) {
	if data[field] == nil && !allowNil {
		return nil, parseErrorNil(parent, field)
	}

	if data[field] == nil {
		return nil, nil
	}

	cmd, err := parser.ParseCommand(data[field])
	if err != nil {
		return nil, parseError(parent, field, err)
	}

	fmt.Println("cmd", cmd)

	renderedCmd, err := ctx.RenderCmd(cmd)
	if err != nil {
		return nil, renderError(parent, field, err)
	}

	return renderedCmd, nil
}
