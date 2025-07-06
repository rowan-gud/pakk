package renderctx

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type RenderCtx struct {
	Mod  *ModContext
	Proj *ProjectContext

	additionalContext map[string]any
	cached            map[string]any
}

func New(proj *ProjectContext, mod *ModContext) *RenderCtx {
	return &RenderCtx{
		Proj: proj,
		Mod:  mod,
	}
}

func NewWithContext(proj *ProjectContext, mod *ModContext, additionalContext map[string]any) *RenderCtx {
	return &RenderCtx{
		Proj:              proj,
		Mod:               mod,
		additionalContext: additionalContext,
	}
}

func (r *RenderCtx) Ctx() map[string]any {
	if r.cached != nil {
		return r.cached
	}

	r.cached = merge(make(map[string]any), r.additionalContext)

	r.cached["Proj"] = r.Proj
	r.cached["Mod"] = r.Mod

	return r.cached
}

func (r *RenderCtx) RenderCmd(cmdStr []string) (*exec.Cmd, error) {
	rendered := make([]string, len(cmdStr))

	for idx, cmd := range cmdStr {
		r, err := r.RenderString(cmd)
		if err != nil {
			return nil, err
		}

		rendered[idx] = r
	}

	return exec.Command(rendered[0], rendered[1:]...), nil
}

func (r *RenderCtx) RenderCmdItems(cmdStr []string, iItems any) ([]*exec.Cmd, error) {
	items := iItems.([]any)

	tpls := make([]*template.Template, len(cmdStr))

	for idx, cmd := range cmdStr {
		tpl, err := template.New("tpl").Parse(cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to construct template %d: %w", idx, err)
		}

		tpls[idx] = tpl
	}

	res := make([]*exec.Cmd, len(items))
	ctx := r.Ctx()

	for idx, item := range items {
		ctx["Item"] = item
		rendered := make([]string, len(tpls))

		for tplIdx, tpl := range tpls {
			var buf bytes.Buffer
			err := tpl.Execute(&buf, ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to execute template <%d, %d>: %w", idx, tplIdx, err)
			}

			rendered[tplIdx] = buf.String()
		}

		res[idx] = exec.Command(rendered[0], rendered[1:]...)
	}

	return res, nil
}

func (r *RenderCtx) RenderItems(templateStr string, iItems any) ([]string, error) {
	items := iItems.([]any)

	tpl, err := template.New("tpl").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to construct template: %w", err)
	}

	res := make([]string, len(items))
	ctx := r.Ctx()

	for idx, item := range items {
		ctx["Item"] = item

		var buf bytes.Buffer
		err := tpl.Execute(&buf, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to execute template %d: %w", idx, err)
		}

		res[idx] = buf.String()
	}

	return res, nil
}

func (r *RenderCtx) RenderSources(sources []string) ([]string, error) {
	res := []string{}

	for _, source := range sources {
		rendered, err := r.RenderString(source)
		if err != nil {
			return nil, err
		}

		if !filepath.IsAbs(rendered) {
			rendered = filepath.Join(r.Mod.Dir, rendered)
		}

		if strings.Contains(rendered, "**") {
			return nil, errors.New("glob ** pattern not supported")
		}

		if strings.Contains(rendered, "*") {
			glob, err := filepath.Glob(rendered)
			if err != nil {
				return nil, err
			}

			res = append(res, glob...)
		} else {
			res = append(res, rendered)
		}
	}

	return res, nil
}

func (r *RenderCtx) RenderStringArray(arr []string) ([]string, error) {
	res := make([]string, len(arr))

	for idx, s := range arr {
		rendered, err := r.RenderString(s)
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", idx, err)
		}

		res[idx] = rendered
	}

	return res, nil
}

func (r *RenderCtx) RenderString(s string) (string, error) {
	tpl, err := template.New("tpl").Parse(s)
	if err != nil {
		return "", fmt.Errorf("failed to construct template: %w", err)
	}

	var buf bytes.Buffer

	if err := tpl.Execute(&buf, r.Ctx()); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (r *RenderCtx) With(additionalContext map[string]any) *RenderCtx {
	additional := merge(make(map[string]any), r.additionalContext, additionalContext)

	return &RenderCtx{
		Proj:              r.Proj,
		Mod:               r.Mod,
		additionalContext: additional,
	}
}

func merge(target map[string]any, sources ...map[string]any) map[string]any {
	for _, source := range sources {
		if source == nil {
			continue
		}

		maps.Copy(target, source)
	}

	return target
}
