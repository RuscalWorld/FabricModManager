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
		for id, breakVer := range mod.Breaks {
			if dep, exact := mod.ResolveDependency(id, breakVer, mods); dep != nil && exact {
				log.Error(log.Danger(dep.Name), "is incompatible with", log.Highlight(mod.Name))
				errors++
			}
		}

		for id, conflictVer := range mod.Conflicts {
			if dep, exact := mod.ResolveDependency(id, conflictVer, mods); dep != nil && exact {
				log.Warn(log.Warning(dep.Name), "is conflicting with", log.Highlight(mod.Name))
				warnings++
			}
		}

		for id, recommendedVer := range mod.Recommends {
			if dep, exact := mod.ResolveDependency(id, recommendedVer, mods); !exact {
				if dep != nil {
					log.Info(fmt.Sprintf("%s %s is recommended to be installed with %s, but currently installed version (%s) doesn't satisfy this recommendation",
						log.Highlight(id), log.Good(recommendedVer), log.Highlight(mod.Name), log.Warning(dep.Version)))
				} else {
					log.Info(log.Warning(id), log.Warning(recommendedVer), "is recommended to be installed with", log.Highlight(mod.Name))
				}
			}
		}

		for id, dependVer := range mod.Depends {
			if dep, exact := mod.ResolveDependency(id, dependVer, mods); !exact {
				if dep != nil {
					log.Error(fmt.Sprintf("%s %s is required to be installed with %s, but currently installed version (%s) doesn't satisfy this requirement",
						log.Highlight(id), log.Warning(dependVer), log.Highlight(mod.Name), log.Danger(dep.Version)))
				} else {
					log.Error(log.Danger(id), log.Danger(dependVer), "must be installed with", log.Highlight(mod.Name))
				}

				errors++
			}
		}
	}

	log.Info(errors, log.Danger("errors"), "and", warnings, log.Warning("warnings"), "found while checking your mod list")
	return nil
}
