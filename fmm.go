package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/urfave/cli/v2"
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
						fmt.Println(i+1, "|", mod.Name, "|", mod.Description)
					}

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
