package shared

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/rowan-gud/pakk/collections"
	"github.com/rowan-gud/pakk/config"
)

var (
	buildFileTree *collections.Trie[string, config.Build]
)

func BuildFileTree(rootDir string, builds map[string]*config.Build) *collections.Trie[string, config.Build] {
	tree := collections.NewTrie[string, config.Build]()

	for path, build := range builds {
		rel, err := filepath.Rel(rootDir, path)
		if err != nil {
			log.Fatal("Could not get relative path", err)
		}

		parts := strings.Split(rel, string(filepath.Separator))

		if parts[0] == "." {
			parts = parts[1:]
		}

		if len(parts) == 0 {
			parts = []string{"/"}
		}

		tree.Add(parts, build)
	}

	buildFileTree = tree

	return tree
}

func FileTree() *collections.Trie[string, config.Build] {
	return buildFileTree
}
