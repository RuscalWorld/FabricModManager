package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetWatchedPaths() []string {
	file, err := os.Open(GetConfigPath("watched.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			err = AddWatchedPath(GetMinecraftDirectory())
			if err != nil {
				log.Fatalln("Unable to save default watched.txt:", err)
				return nil
			}
		} else {
			log.Fatalln("Unable to read watched.txt:", err)
			return nil
		}
	}

	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln("Unable to read from watched.txt:", err)
		return nil
	}

	paths := strings.Split(string(bytes), "\n")
	if len(paths) < 1 {
		return []string{GetMinecraftDirectory()}
	}

	return paths
}

func AddWatchedPath(path string) error {
	path = strings.ReplaceAll(path, "\\", "/")
	file, err := os.Open(GetConfigPath("watched.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(GetConfigPath("watched.txt")), os.ModePerm)
			if err != nil {
				return err
			}

			file, err = os.Create(GetConfigPath("watched.txt"))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err == nil && len(string(bytes)) > 0 {
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString(path)
	if err != nil {
		return err
	}

	return nil
}

func IsWatchingPath(path string) bool {
	path = strings.ReplaceAll(path, "\\", "/")
	paths := GetWatchedPaths()

	for _, cpath := range paths {
		cpath = strings.ReplaceAll(cpath, "\\", "/")
		if cpath == path {
			return true
		}
	}

	return false
}

func RemoveWatchedPath(path string) error {
	file, err := os.Open(GetConfigPath("watched.txt"))
	if err != nil {
		return err
	}

	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	path = strings.ReplaceAll(path, "\\", "/")
	paths := strings.Split(string(bytes), "\n")
	newPaths := make([]string, 0)

	for _, cpath := range paths {
		cpath = strings.ReplaceAll(cpath, "\\", "/")
		if cpath != path {
			newPaths = append(newPaths, cpath)
		}
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(newPaths, "\n"))
	if err != nil {
		return err
	}

	return nil
}

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
