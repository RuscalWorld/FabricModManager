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
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
