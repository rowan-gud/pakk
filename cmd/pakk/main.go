package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/rowan-gud/pakk/cmd/pakk/commands"
	"github.com/rowan-gud/pakk/cmd/pakk/shared"
)

var (
	rootCmd = &cobra.Command{
		Use:              "pakk",
		Short:            "A build tool for any language",
		PersistentPreRun: prepare,
	}
)

func init() {
	rootCmd.AddCommand(commands.QueryCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func prepare(cmd *cobra.Command, args []string) {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working dir", err)
	}

	_, err = shared.ParsePackageConfig(rootDir)
	if err != nil {
		log.Fatal("Could not parse package config", err)
	}

	builds, err := shared.ParseBuildConfigs(rootDir)
	if err != nil {
		log.Fatal("Could not parse build config files", err)
	}

	tree := shared.BuildFileTree(rootDir, builds)
	tree.Print()
}
