package commands

import (
	"fmt"
	"path"

	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/RuscalWorld/FabricModManager/core"
	"github.com/urfave/cli/v2"
)

func ListMods(_ *cli.Context) error {
	mods, err := core.GetMods(path.Join(config.Global.WorkDir, "mods"), false)
	if err != nil {
		return err
	}

	for i, mod := range *mods {
		fmt.Println(i+1, "|", mod.Name, "|", mod.Description, "|", len(mod.Nested), "nested mods")
	}

	return nil
}
