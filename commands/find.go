package commands

import (
	"github.com/RuscalWorld/FabricModManager/log"
	"github.com/RuscalWorld/FabricModManager/remote"
	"github.com/urfave/cli/v2"
)

func FindMod(ctx *cli.Context) error {
	name := ctx.Args().Get(0)
	info, err := remote.FindMod(name)
	if err != nil {
		return err
	}

	mod := (*info).(remote.ModrinthMod)
	log.Fine("Found a mod at modrinth: " + mod.GetName())
	log.Fine(mod.Description)

	return nil
}
