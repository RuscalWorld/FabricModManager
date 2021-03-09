package main

import (
	"fmt"
	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
)

var WorkDir string

func main() {
	dir, err := os.Getwd()
	if err != nil || !config.IsMinecraftDirectory(dir) {
		WorkDir = config.GetMinecraftDirectory()
	} else {
		WorkDir = dir
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "list",
				Aliases:     []string{"l"},
				Description: "Shows info about installed mods",
				Action: func(context *cli.Context) error {
					mods, err := GetMods(path.Join(WorkDir, "mods"), false)
					if err != nil {
						return err
					}

					for i, mod := range *mods {
						fmt.Println(i+1, "|", mod.Name, "|", mod.Description, "|", len(mod.Nested), "nested mods")
					}

					return nil
				},
			},
			{
				Name:        "check",
				Aliases:     []string{"c"},
				Description: "Checks your mod list for conflicts and unmet dependencies",
				Action: func(context *cli.Context) error {
					modMap, err := GetFullModMap(nil)
					if err != nil {
						return err
					}

					errors, warnings := 0, 0

					for _, mod := range modMap {
						for id, breakVer := range mod.Breaks {
							if value, ok := modMap[id]; ok && CheckVersions(value.Version, breakVer) {
								fmt.Println("!!!", value.Name, "is incompatible with", mod.Name)
								errors++
							}
						}

						for id, conflictVer := range mod.Conflicts {
							if value, ok := modMap[id]; ok && CheckVersions(value.Version, conflictVer) {
								fmt.Println(value.Name, "is conflicting with", mod.Name)
								warnings++
							}
						}

						for id, recommendedVer := range mod.Recommends {
							if value, ok := modMap[id]; !ok || !CheckVersions(value.Version, recommendedVer) {
								fmt.Println(id, recommendedVer, "is recommended to be installed with", mod.Name)
								warnings++
							}
						}

						for id, dependVer := range mod.Depends {
							if value, ok := modMap[id]; !ok || !CheckVersions(value.Version, dependVer) {
								fmt.Println("!!!", id, dependVer, "is required for", mod.Name)
								if ok {
									fmt.Println(value.Version, "is installed, version", dependVer, "is required")
								}
								errors++
							}
						}
					}

					fmt.Println(errors, "errors and", warnings, "warnings found while checking your mod list")
					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
