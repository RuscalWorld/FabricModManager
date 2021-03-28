package main

import (
	"fmt"
	"os"
	"path"

	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/RuscalWorld/FabricModManager/log"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetupLogging()

	dir, err := os.Getwd()
	if err != nil || !config.IsMinecraftDirectory(dir) {
		config.Global.WorkDir = config.GetMinecraftDirectory()
	} else {
		config.Global.WorkDir = dir
	}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "check-mc",
				Usage:       "Should FMM check your Minecraft version?",
				Destination: &config.Global.CheckMinecraft,
			},
			&cli.StringFlag{
				Name:        "mc-version",
				Usage:       "Explicitly specify Minecraft version",
				Destination: &config.Global.MinecraftVersion,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "list",
				Aliases:     []string{"l"},
				Description: "Shows info about installed mods",
				Action: func(context *cli.Context) error {
					mods, err := GetMods(path.Join(config.Global.WorkDir, "mods"), false)
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
					if config.Global.CheckMinecraft {
						if config.Global.MinecraftVersion == "" {
							config.Global.MinecraftVersion, err = GetCurrentMinecraftVersion()

							if err != nil {
								log.Warn("Unable to determine current Minecraft version:", err)
								log.Warn("It looks like you are either using non-official launcher or have no Fabric installed")
								log.Warn("You can provide Minecraft version explicitly using --mc-version flag")
							}
						}
					}

					if config.Global.MinecraftVersion != "" && config.Global.CheckMinecraft {
						log.Info("Assuming that you're using Minecraft", config.Global.MinecraftVersion)
					}

					mods, err := GetAllMods()
					if err != nil {
						return err
					}

					errors, warnings := 0, 0

					for _, mod := range *mods {
						for id, breakVer := range mod.Breaks {
							if dep, exact := mod.ResolveDependency(id, breakVer, mods); dep != nil && exact {
								log.Error(log.Danger(dep.Name), "is incompatible with", log.Highlight(mod.Name))
								errors++
							}
						}

						for id, conflictVer := range mod.Conflicts {
							if dep, exact := mod.ResolveDependency(id, conflictVer, mods); dep != nil && exact {
								log.Warn(log.Warning(dep.Name), "is conflicting with", log.Highlight(mod.Name))
								warnings++
							}
						}

						for id, recommendedVer := range mod.Recommends {
							if dep, exact := mod.ResolveDependency(id, recommendedVer, mods); !exact {
								if dep != nil {
									log.Info(fmt.Sprintf("%s %s is recommended to be installed with %s, but currently installed version (%s) doesn't satisfy this recommendation",
										log.Highlight(id), log.Good(recommendedVer), log.Highlight(mod.Name), log.Warning(dep.Version)))
								} else {
									log.Info(log.Warning(id), log.Warning(recommendedVer), "is recommended to be installed with", log.Highlight(mod.Name))
								}

								warnings++
							}
						}

						for id, dependVer := range mod.Depends {
							if dep, exact := mod.ResolveDependency(id, dependVer, mods); !exact {
								if dep != nil {
									log.Error(fmt.Sprintf("%s %s is required to be installed with %s, but currently installed version (%s) doesn't satisfy this requirement",
										log.Highlight(id), log.Warning(dependVer), log.Highlight(mod.Name), log.Danger(dep.Version)))
								} else {
									log.Error(log.Danger(id), log.Danger(dependVer), "must be installed with", log.Highlight(mod.Name))
								}

								errors++
							}
						}
					}

					log.Info(errors, log.Danger("errors"), "and", warnings, log.Warning("warnings"), "found while checking your mod list")
					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
