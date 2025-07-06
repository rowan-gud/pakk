package mod

import (
	"os/exec"

	"github.com/rowan-gud/pakk/config/parser"
	"github.com/rowan-gud/pakk/config/renderctx"
)

type Bin struct {
	Artifacts []string
	Cmd       *exec.Cmd
}

func parseBin(data any, ctx *renderctx.RenderCtx) (*Bin, error) {
	root, err := parser.ParseMap(data)
	if err != nil {
		return nil, parseRootError("Bin", err)
	}

	artifacts, err := parseRenderSources(root, "Bin", "artifacts", false, ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := parseRenderCommand(root, "Bin", "cmd", false, ctx)
	if err != nil {
		return nil, err
	}

	return &Bin{
		Artifacts: artifacts,
		Cmd:       cmd,
	}, nil
}
