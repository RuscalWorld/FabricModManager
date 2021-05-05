package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/RuscalWorld/FabricModManager/config"
	"github.com/RuscalWorld/FabricModManager/core"
	"github.com/RuscalWorld/FabricModManager/log"
	"github.com/RuscalWorld/FabricModManager/remote"
	"github.com/urfave/cli/v2"
)

func InstallMod(ctx *cli.Context) error {
	mods, err := core.GetAllMods()
	if err != nil {
		return fmt.Errorf("unable to retrieve list of installed mods: %s", err)
	}

	for _, name := range ctx.Args().Slice() {
		log.Info(fmt.Sprintf("Searching mod %s", log.Highlight(name)))
		err = DownloadMod(name, *mods)
		if err != nil {
			return fmt.Errorf("unable to install mod %s: %s", name, err)
		}
	}

	return nil
}

func DownloadMod(name string, mods []core.FabricMod) error {
	mod, err := remote.FindMod(name)
	if err != nil {
		return fmt.Errorf("unable to find mod: %s", err)
	}

	log.Fine(fmt.Sprintf("Found mod %s at %s", log.Highlight(mod.GetName()), log.Highlight(mod.GetRemote().GetDisplayName())))

	jarPath := path.Join(config.GetCachePath(), mod.GetName()+".jar")
	if _, err = os.Stat(jarPath); err == nil && !config.Global.Force {
		log.Warn(fmt.Sprintf("Destination file already exists. It looks like this mod was downloaded before. "+
			"If you want to download new mod anyway, run install command with %s flag.", log.Highlight("--force")))
	} else if os.IsNotExist(err) || config.Global.Force {
		log.Info("Retrieving mod versions")

		versions, err := mod.GetVersions()
		if err != nil {
			return fmt.Errorf("unable to retrieve mod versions: %s", err)
		}

		version := (*versions)[0]
		log.Fine(fmt.Sprintf("Found version %s, downloading it", version.GetName()))

		startTime := time.Now().UnixNano() / 1000000
		response, err := http.Get(version.GetDownloadURL())
		if err != nil {
			return fmt.Errorf("unable to download mod: %s", err)
		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Fatal(fmt.Sprintf("unable to close response body reader: %s", err))
			}
		}(response.Body)

		err = os.MkdirAll(config.GetCachePath(), 0660)
		if err != nil {
			return fmt.Errorf("unable to create directoies: %s", err)
		}

		file, err := os.Create(path.Join(config.GetCachePath(), mod.GetName()+".jar"))
		if err != nil {
			return fmt.Errorf("unable to create output file: %s", err)
		}

		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			return fmt.Errorf("unable to save data to output file: %s", err)
		}

		log.Fine(fmt.Sprintf("Successfully downloaded %d bytes in %dms", response.ContentLength, time.Now().UnixNano()/1000000-startTime))
	} else {
		return fmt.Errorf("unable to check existance of destination file: %x", err)
	}

	info, err := core.GetModInfo(jarPath)
	if err != nil {
		return fmt.Errorf("unable to read mod: %s", err)
	}

	log.Info(fmt.Sprintf("Checking compatibility of %s with your installed mods", log.Highlight(mod.GetName())))

	mods = append(mods, *info)
	incompatibilities := 0

	for _, installedMod := range mods {
		breaks := installedMod.GetBreaks(&mods)
		breakNames := make([]string, 0)
		for _, breakMod := range *breaks {
			breakNames = append(breakNames, log.Highlight(breakMod.Name))
		}

		if len(breakNames) > 0 {
			log.Error(fmt.Sprintf("%s is incompatible with some mods: %s", log.Danger(installedMod.GetName()), strings.Join(breakNames, ", ")))
			incompatibilities += len(breakNames)
		}
	}

	if incompatibilities > 0 {
		if config.Global.IgnoreIncompatibilities {
			log.Warn(fmt.Sprintf("%d incompatibilities were found while checking list of your mods, "+
				"but %s flag is set, so ignoring this issue", incompatibilities, log.Highlight("--ignore-incompatibilities")))
		} else {
			return fmt.Errorf("%d incompatibilities were found while checking list of your mods. "+
				"You should remove incompatible mods or not to install incompatible mods. "+
				"If you want to install mod anyway and ignore this error, rerun install command with %s flag",
				incompatibilities, log.Highlight("--ignore-incompatibilities"))
		}
	} else {
		log.Fine("No incompatibilities found")
	}

	if !config.Global.NoDependencyChecks {
		log.Info(fmt.Sprintf("Checking if any other mods required for %s to work properly", log.Highlight(mod.GetName())))
		dependencies := *info.GetMissingDependencies(&mods)

		if len(dependencies) > 0 {
			for dependency, version := range dependencies {
				log.Info(fmt.Sprintf("%s requires mod %s %s, but it is not installed, attempting to install",
					log.Highlight(mod.GetName()), log.Warning(dependency), log.Warning(version)))

				err = DownloadMod(fmt.Sprintf("%s", dependency), mods)
				if err != nil {
					return fmt.Errorf("unable to download dependency (%s): %s \n"+
						"If you want to skip dependency checks and install mod anyway, "+
						"rerun install command with %s flag", dependency, err, log.Highlight("--no-dependencies"))
				}
			}
		} else {
			log.Fine(fmt.Sprintf("%s has no unmet dependencies", log.Highlight(mod.GetName())))
		}
	}

	srcFile, err := os.Open(jarPath)
	if err != nil {
		return fmt.Errorf("unable to open cached mod file: %s", err)
	}

	dstFile, err := os.Create(path.Join(config.Global.WorkDir, "mods", filepath.Base(jarPath)))
	if err != nil {
		return fmt.Errorf("unable to create mod file in game directory: %s", err)
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("unable to copy cached file to the game directory: %s", err)
	}

	log.Fine(fmt.Sprintf("Successfully installed mod %s", log.Highlight(info.GetName())))

	return nil
}
