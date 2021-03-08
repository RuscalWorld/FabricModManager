package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "path",
				Aliases:     []string{"p"},
				Description: "Manage list of paths watched by FMM",
				Subcommands: []*cli.Command{
					{
						Name:        "list",
						Aliases:     []string{"l"},
						Description: "Shows list of watched paths",
						Action: func(context *cli.Context) error {
							paths := config.GetWatchedPaths()
							fmt.Println("There are", len(paths), "paths registered:")
							for i, path := range paths {
								fmt.Println(i+1, ":", path)
							}
							return nil
						},
					},
					{
						Name:        "add",
						Aliases:     []string{"a"},
						Description: "Adds path to watched list",
						Action: func(context *cli.Context) error {
							if context.NArg() == 0 {
								return errors.New("You should specify the path you want to add ")
							}

							path := context.Args().Get(0)
							if config.IsWatchingPath(path) {
								return errors.New("This path is already added to watched list ")
							}

							err := config.AddWatchedPath(path)
							if err != nil {
								return err
							}

							fmt.Println("Successfully added", path, "to watched path list")
							return nil
						},
					},
					{
						Name:        "remove",
						Aliases:     []string{"r"},
						Description: "Remove path from watched list",
						Action: func(context *cli.Context) error {
							if context.NArg() == 0 {
								return errors.New("You should specify the path you want to remove ")
							}

							path := context.Args().Get(0)
							if !config.IsWatchingPath(path) {
								return errors.New("This path isn't watched ")
							}

							err := config.RemoveWatchedPath(path)
							if err != nil {
								return err
							}

							fmt.Println("Successfully removed", path, "from watched list")
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
