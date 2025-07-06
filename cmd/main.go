package main

import (
	"log"
	"os"
	"path/filepath"
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
		OutDir:  out,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := buildProject(project); err != nil {
		log.Fatal(err)
	}
}
