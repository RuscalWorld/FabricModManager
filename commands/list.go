package commands

import (
	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/RuscalWorld/FabricModManager/core"
	"github.com/RuscalWorld/FabricModManager/log"
	"github.com/urfave/cli/v2"
	"path"
	"strings"
)

func ListMods(_ *cli.Context) error {
	mods, err := core.GetMods(path.Join(config.Global.WorkDir, "mods"), false)
	if err != nil {
		return err
	}

	log.Info("Mod list")
	for _, mod := range *mods {
		PrintModInfo(mod, 0, false)
	}

	return nil
}

func PrintModInfo(mod core.FabricMod, depth int, last bool) {
	output := strings.Repeat(" ║", depth)

	if last {
		output += " ╚ "
	} else {
		output += " ╠ "
	}

	output += mod.GetName() + " " + log.Highlight(mod.Version)
	log.Info(output)
	for i, nested := range mod.Nested {
		PrintModInfo(nested, depth+1, i == len(mod.Nested)-1)
	}
}
