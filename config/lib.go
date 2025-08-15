package config

type Lib struct {
	Lang    string   `toml:"lang"`
	Imports []string `toml:"imports"`
	Sources []string `toml:"sources"`
}
