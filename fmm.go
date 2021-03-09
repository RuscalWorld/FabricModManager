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
					mods, err := GetMods(path.Join(WorkDir, "mods"))
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

					modMap["minecraft"] = FabricMod{}
					modMap["fabricloader"] = FabricMod{}
					errors, warnings := 0, 0

					for _, mod := range modMap {
						for id := range mod.Breaks {
							if value, ok := modMap[id]; ok {
								fmt.Println("!!!", value.Name, "is incompatible with", mod.Name)
								errors++
							}
						}

						for id := range mod.Conflicts {
							if value, ok := modMap[id]; ok {
								fmt.Println(value.Name, "is conflicting with", mod.Name)
								warnings++
							}
						}

						for id, version := range mod.Recommends {
							if _, ok := modMap[id]; !ok {
								fmt.Println(id, version, "is recommended to be installed with", mod.Name)
								warnings++
							}
						}

						for id, version := range mod.Depends {
							if _, ok := modMap[id]; !ok {
								fmt.Println("!!!", id, version, "is required for", mod.Name)
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
