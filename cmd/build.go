package main

import "github.com/rowan-gud/pakk/config/proj"

func buildProject(project *proj.Project) error {
	for _, mod := range project.Modules {
		if mod.Pkg != nil {
			for _, pre := range mod.Pkg.Pre {
				for _, cmd := range pre.Run {
					if err := cmd.Run(); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
