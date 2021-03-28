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
					mods, err := GetAllMods()
					if err != nil {
						return err
					}

					errors, warnings := 0, 0

					for _, mod := range *mods {
						for id, breakVer := range mod.Breaks {
							if dep, exact := mod.ResolveDependency(id, breakVer, mods); dep != nil && exact {
								fmt.Println("!!!", dep.Name, "is incompatible with", mod.Name)
								errors++
							}
						}

						for id, conflictVer := range mod.Conflicts {
							if dep, exact := mod.ResolveDependency(id, conflictVer, mods); dep != nil && exact {
								fmt.Println(dep.Name, "is conflicting with", mod.Name)
								warnings++
							}
						}

						for id, recommendedVer := range mod.Recommends {
							if dep, exact := mod.ResolveDependency(id, recommendedVer, mods); !exact {
								if dep != nil {
									fmt.Println(fmt.Sprintf("%s %s is recommended to be installed with %s, but currently installed version (%s) doesn't satisfy this recommendation", id, recommendedVer, mod.Name, dep.Version))
								} else {
									fmt.Println(id, recommendedVer, "is recommended to be installed with", mod.Name)
								}

								warnings++
							}
						}

						for id, dependVer := range mod.Depends {
							if dep, exact := mod.ResolveDependency(id, dependVer, mods); !exact {
								if dep != nil {
									fmt.Println(fmt.Sprintf("!!! %s %s is required to be installed with %s, but currently installed version (%s) doesn't satisfy this requirement %t", id, dependVer, mod.Name, dep.Version, exact))
								} else {
									fmt.Println(id, dependVer, "must be installed with", mod.Name)
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
