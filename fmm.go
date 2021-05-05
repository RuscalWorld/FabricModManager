package main

import (
	"os"

	"github.com/RuscalWorld/FabricModManager/commands"
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
				Action:      commands.ListMods,
			},
			{
				Name:        "check",
				Aliases:     []string{"c"},
				Description: "Checks your mod list for conflicts and unmet dependencies",
				Action:      commands.CheckMods,
			},
			{
				Name:        "find",
				Description: "Finds mod at available data sources",
				Action:      commands.FindMod,
				ArgsUsage:   "[mod name]",
			},
			{
				Name:        "install",
				Description: "Installs mod with given name",
				Action:      commands.InstallMod,
				ArgsUsage:   "[mod name]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "force",
						Usage:       "Download mod anyway, remove old mod if it exists",
						Aliases:     []string{"f"},
						Destination: &config.Global.Force,
					},
					&cli.BoolFlag{
						Name:        "ignore-incompatibilities",
						Usage:       "Forcibly install mod that is incompatible with already installed mods",
						Aliases:     []string{"i"},
						Destination: &config.Global.IgnoreIncompatibilities,
					},
					&cli.BoolFlag{
						Name:        "no-dependencies",
						Usage:       "Do not install dependencies of mod",
						Destination: &config.Global.NoDependencyChecks,
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
