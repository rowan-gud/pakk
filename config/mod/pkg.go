package mod

import (
	"os/exec"

	"github.com/rowan-gud/pakk/config/parser"
	"github.com/rowan-gud/pakk/config/renderctx"
)

type Pkg struct {
	Sources  []string
	Pre      []Pre
	Provides Provides
}

type Pre struct {
	Each      []string
	Run       []*exec.Cmd
	Generates []string
}

type Provides struct {
	Import *ProvidesImport
}

type ProvidesImport struct {
	Path []string
}

func parsePkg(data any, ctx *renderctx.RenderCtx) (*Pkg, error) {
	root, err := parser.ParseMap(data)
	if err != nil {
		return nil, parseRootError("Pkg", err)
	}

	sources, err := parseRenderSources(root, "Pkg", "sources", true, ctx)
	if err != nil {
		return nil, err
	}

	var pre []Pre

	if root["pre"] == nil {
		pre = []Pre{}
	} else {
		roots, err := parser.ParseMapArray(root["pre"])
		if err != nil {
			return nil, parseError("Pkg", "pre", err)
		}

		pre = make([]Pre, len(roots))

		for idx, root := range roots {
			eachStrs, err := parseRenderSources(root, "Pre", "each", true, ctx)
			if err != nil {
				return nil, err
			}

			each := make([]any, len(eachStrs))
			for idx, e := range eachStrs {
				each[idx] = e
			}

			runCmd, err := parser.ParseCommand(root["run"])
			if err != nil {
				return nil, parseError("Pre", "run", err)
			}

			var run []*exec.Cmd

			if len(each) == 0 {
				cmd, err := ctx.RenderCmd(runCmd)
				if err != nil {
					return nil, renderError("Pre", "run", err)
				}

				run = []*exec.Cmd{cmd}
			} else {
				cmds, err := ctx.RenderCmdItems(runCmd, each)
				if err != nil {
					return nil, renderError("Pre", "run", err)
				}

				run = cmds
			}

			generates, err := parseRenderSources(root, "Pre", "generates", false, ctx)
			if err != nil {
				return nil, err
			}

			pre[idx] = Pre{
				Each:      eachStrs,
				Run:       run,
				Generates: generates,
			}
		}
	}

	var provides Provides

	if root["provides"] != nil {
		providesRoot, err := parser.ParseMap(root["provides"])
		if err != nil {
			return nil, parseRootError("Provides", err)
		}

		if providesRoot["import"] != nil {
			providesImportRoot, err := parser.ParseMap(providesRoot["import"])
			if err != nil {
				return nil, parseRootError("ProvidesImport", err)
			}

			path, err := parser.ParseStringArray(providesImportRoot["path"])
			if err != nil {
				return nil, parseError("ProvidesImport", "path", err)
			}

			renderedPath, err := ctx.RenderStringArray(path)
			if err != nil {
				return nil, renderError("ProvidesImport", "path", err)
			}

			provides.Import = &ProvidesImport{
				Path: renderedPath,
			}
		}
	}

	return &Pkg{
		Sources:  sources,
		Pre:      pre,
		Provides: provides,
	}, nil
}
