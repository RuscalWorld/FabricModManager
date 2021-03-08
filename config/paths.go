package config

import (
	"log"
	"os"
	"path"
	"runtime"
)

func GetConfigPath(filename string) string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	return path.Join(dir, ".fmm", filename)
}

func GetMinecraftDirectory() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(dir)
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		return path.Join(dir, "AppData", "Roaming", ".minecraft")
	case "linux":
		return path.Join(dir, ".minecraft")
	default:
		log.Fatalln("Unable to determine default directory of your Minecraft for your OS. Please add your .minecraft to", GetConfigPath("watched.txt"))
		return ""
	}
}
