package commands

import (
	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/RuscalWorld/FabricModManager/core"
	"github.com/RuscalWorld/FabricModManager/log"
)

func GetMinecraftInfo() {
	var err error
	if config.Global.CheckMinecraft {
		if config.Global.MinecraftVersion == "" {
			config.Global.MinecraftVersion, err = core.GetCurrentMinecraftVersion()

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
}
