package commands

import (
	"fmt"
	"strings"

	"github.com/rowan-gud/pakk/cmd/pakk/shared"
	"github.com/spf13/cobra"
)

var (
	QueryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query a module in the build graph",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Querying the module", args[0])

			fileTree := shared.FileTree()

			query := args[0]

			if !strings.HasPrefix(query, "//") {
				return fmt.Errorf("query must start with //")
			}

			query = strings.TrimPrefix(query, "//")

			parts := strings.Split(query, "/")

			if parts[0] == ".." {
				return fmt.Errorf("query cannot start with ..")
			}

			if parts[0] == "." {
				parts = parts[1:]
			}

			if len(parts) == 0 {
				parts = []string{"/"}
			}

			selector := strings.Split(parts[len(parts)-1], ":")

			parts[len(parts)-1] = selector[0]
			selector = selector[1:]

			build, err := fileTree.Get(parts)
			if err != nil {
				return fmt.Errorf("could not get build %s: %w", args[0], err)
			}

			fmt.Println("Build", build)

			switch selector[0] {
			case "bin":
				bin, ok := build.Bin[selector[1]]
				if !ok {
					return fmt.Errorf("bin %s not found", selector[1])
				}

				fmt.Println("Bin", bin)
			case "lib":
				lib, ok := build.Lib[selector[1]]
				if !ok {
					return fmt.Errorf("lib %s not found", selector[1])
				}

				fmt.Println("Lib", lib)
			default:
				return fmt.Errorf("invalid selector %s", selector[0])
			}

			return nil
		},
	}
)
