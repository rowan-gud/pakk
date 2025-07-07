package render

type RenderContext struct {
	Mod     ModContext      `toml:"mod"`
	Project *ProjectContext `toml:"project"`
	Values  map[string]any  `toml:"values"`
}

type ProjectContext struct {
	Out    string         `toml:"out"`
	Path   string         `toml:"path"`
	Values map[string]any `toml:"values"`
}

type ModContext struct {
	Path string `toml:"path"`
}
