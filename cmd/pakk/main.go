package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working dir", err)
	}

	packageConfig, err := parsePackageConfig(rootDir)
	if err != nil {
		log.Fatal("Could not parse package config", err)
	}

	builds, err := parseBuildConfigs(rootDir)
	if err != nil {
		log.Fatal("Could not parse build config files", err)
	}

	tree := buildFileTree(rootDir, builds)
	tree.Print()

	parts, err := builds["/home/rowan/source/pakk/utils"].Resolve("//config")
	if err != nil {
		log.Fatal("Could not resolve path", err)
	}
	fmt.Println("Resolve //config:lib:config", parts)

	paths := [][]string{
		{"/", "cmd", "pakk"},
		{"/", "collections"},
		{"/", "config"},
		{"/", "utils"},
	}

	t, _ := toml.Marshal(packageConfig)
	log.Println("Package", string(t))

	for _, path := range paths {
		v, err := tree.Get(path)
		if err != nil {
			log.Fatal("Could not get build ", err)
		}
		t, _ := toml.Marshal(v)
		log.Println("build", path)
		fmt.Println(string(t))
	}
}
