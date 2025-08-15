package config

type Bin struct {
	Lang    string   `toml:"lang"`
	Imports []string `toml:"imports"`
	Sources []string `toml:"sources"`
}
