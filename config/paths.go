package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/RuscalWorld/FabricModManager/log"
)

func GetRootDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(dir, ".fmm")
}

func GetConfigPath(filename string) string {
	return path.Join(GetRootDir(), filename)
}

func GetCachePath() string {
	return path.Join(GetRootDir(), "mods")
}

func IsMinecraftDirectory(path string) bool {
	modsPath := filepath.Join(path, "mods")
	_, err := ioutil.ReadDir(modsPath)
	if err != nil {
		return false
	}
	return true
}

func GetMinecraftDirectory() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(dir)
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		return path.Join(dir, "AppData", "Roaming", ".minecraft")
	case "linux":
		return path.Join(dir, ".minecraft")
	default:
		log.Fatal("Unable to determine default directory of your Minecraft for your OS. Please add your .minecraft to", GetConfigPath("watched.txt"))
		return ""
	}
}
