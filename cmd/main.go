package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out := filepath.Join(root, ".pakk")

	_ = os.MkdirAll(out, 0755)

	project, err := parseConfig(&parseConfigOptions{
		RootDir: root,
	})
	if err != nil {
		log.Fatal(err)
	}

	marshalled, err := toml.Marshal(project)
	if err != nil {
		slog.Error("failed to marshal project",
			slog.Any("error", err),
		)
	}

	log.Println(string(marshalled))

	marshalled, err = toml.Marshal(map[string]any{"modules": project.Modules()})
	if err != nil {
		slog.Error("failed to marshal modules",
			slog.Any("error", err),
		)
	}

	log.Println(string(marshalled))
}
