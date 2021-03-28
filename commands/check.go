package commands

import (
	"fmt"
	"github.com/RuscalWorld/FabricModManager/core"
	"github.com/RuscalWorld/FabricModManager/log"
	"github.com/urfave/cli/v2"
)

func CheckMods(_ *cli.Context) error {
	GetMinecraftInfo()
	mods, err := core.GetAllMods()
	if err != nil {
		return err
	}

	errors, warnings := 0, 0

	for _, mod := range *mods {
		for _, breaking := range *mod.GetBreaks(mods) {
			log.Error(log.Danger(breaking.Name), "is incompatible with", log.Highlight(mod.Name))
			errors++
		}

		for _, conflict := range *mod.GetConflicts(mods) {
			log.Warn(log.Warning(conflict.Name), "is conflicting with", log.Highlight(mod.Name))
			warnings++
		}

		for id, version := range *mod.GetMissingRecommends(mods) {
			if recommend, _ := mod.ResolveDependency(id, version, mods); recommend != nil {
				log.Info(fmt.Sprintf("%s %s is recommended to be installed with %s, but currently installed version (%s) doesn't satisfy this recommendation",
					log.Highlight(id), log.Good(version), log.Highlight(mod.Name), log.Warning(recommend.Version)))
			} else {
				log.Info(log.Warning(id), log.Warning(version), "is recommended to be installed with", log.Highlight(mod.Name))
			}
		}

		for id, version := range *mod.GetMissingDependencies(mods) {
			if dependency, _ := mod.ResolveDependency(id, version, mods); dependency != nil {
				log.Error(fmt.Sprintf("%s %s must be installed with %s, but currently installed version (%s) doesn't satisfy this requirement",
					log.Highlight(id), log.Warning(version), log.Highlight(mod.Name), log.Danger(dependency.Version)))
			} else {
				log.Error(log.Danger(id), log.Danger(version), "must be installed with", log.Highlight(mod.Name))
			}

			errors++
		}
	}

	log.Info(errors, log.Danger("errors"), "and", warnings, log.Warning("warnings"), "found while checking your mod list")
	return nil
}
