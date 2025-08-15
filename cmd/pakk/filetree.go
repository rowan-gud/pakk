package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/rowan-gud/pakk/collections"
	"github.com/rowan-gud/pakk/config"
)

func buildFileTree(rootDir string, builds map[string]*config.Build) *collections.Trie[string, config.Build] {
	tree := collections.NewTrie[string, config.Build]()

	for path, build := range builds {
		rel, err := filepath.Rel(rootDir, path)
		if err != nil {
			log.Fatal("Could not get relative path", err)
		}

		fmt.Println("Rel", rel)
		parts := strings.Split(rel, string(filepath.Separator))
		fmt.Println("Parts", parts)

		if parts[0] == "." {
			parts = parts[1:]
		}

		parts = append([]string{"/"}, parts...)

		fmt.Println("Adding", parts, build)

		tree.Add(parts, build)
	}

	return tree
}
