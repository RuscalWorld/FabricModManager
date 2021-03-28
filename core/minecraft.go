package core

import "github.com/RuscalWorld/FabricModManager/config"

func GetMinecraftDependency() *FabricMod {
	return &FabricMod{
		FabricModInfo: FabricModInfo{
			ID:      "minecraft",
			Version: config.Global.MinecraftVersion,
		},
	}
}
