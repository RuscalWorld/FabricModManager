package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
)

type FabricMod struct {
	ID          string                 `json:"id"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	Depends     map[string]interface{} `json:"depends"`
	Recommends  map[string]interface{} `json:"recommends"`
	Suggests    map[string]interface{} `json:"suggests"`
	Breaks      map[string]interface{} `json:"breaks"`
	Conflicts   map[string]interface{} `json:"conflicts"`
	JARs        *[]NestedJAR           `json:"jars"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Nested      []FabricMod
}

type NestedJAR struct {
	File string `json:"file"`
}

func (m *FabricMod) ResolveNestedDependency(id string) *FabricMod {
	if m.ID == id {
		return m
	}

	for _, mod := range m.Nested {
		if mod.ID == id {
			return &mod
		}

		dependency := mod.ResolveNestedDependency(id)
		if dependency != nil {
			return dependency
		}
	}

	return nil
}

func (m FabricMod) ResolveDependency(id string, version interface{}, mods *[]FabricMod) (*FabricMod, bool) {
	if mod := m.ResolveNestedDependency(id); mod != nil {
		return mod, true
	}

	if mods == nil {
		return nil, false
	}

	var best *FabricMod
	for _, mod := range *mods {
		dependency := mod.ResolveNestedDependency(id)
		if dependency != nil {
			satisfy := CheckVersions(dependency.Version, version)
			if satisfy {
				return dependency, true
			}

			best = dependency
		}
	}

	return best, false
}

func CheckVersions(ver string, constraint interface{}) bool {
	if versions, ok := constraint.([]interface{}); ok {
		for _, required := range versions {
			if required == ver {
				return true
			}
		}

		return false
	}

	required := constraint.(string)
	if required == "*" {
		return true
	}

	ver = strings.TrimPrefix(ver, "v")
	required = strings.TrimPrefix(required, "v")

	vConstraint, err := version.NewConstraint(required)
	if err != nil {
		return false
	}

	vVer, err := version.NewVersion(ver)
	if err != nil {
		return false
	}

	return vConstraint.Check(vVer)
}

func GetModInfo(path string) (*FabricMod, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ReadModInfo(file)
}

func ReadModInfo(input []byte) (*FabricMod, error) {
	reader, err := zip.NewReader(bytes.NewReader(input), int64(len(input)))
	if err != nil {
		return nil, err
	}

	var mod *FabricMod
	nested := make([]FabricMod, 0)
	nestedPaths := make([]string, 0)

	for _, file := range reader.File {
		if file.Name == "fabric.mod.json" {
			file, err := file.Open()
			if err != nil {
				return nil, err
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}

			info := &FabricMod{}
			err = json.Unmarshal(data, info)
			if err != nil {
				return nil, err
			}

			mod = info

			if mod.JARs != nil {
				for _, jar := range *mod.JARs {
					nestedPaths = append(nestedPaths, jar.File)
				}
			}

			break
		}
	}

	for _, file := range reader.File {
		isNestedJar := false
		for _, nestedPath := range nestedPaths {
			if nestedPath == file.Name {
				isNestedJar = true
			}
		}

		if !isNestedJar {
			continue
		}

		file, err := file.Open()
		if err != nil {
			continue
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			continue
		}

		nestedMod, err := ReadModInfo(data)
		if err != nil {
			continue
		}

		nested = append(nested, *nestedMod)
	}

	if mod == nil {
		return nil, errors.New("Input JAR file wasn't a Fabric mod ")
	} else {
		mod.Nested = nested
		return mod, nil
	}
}

func GetMods(dirname string, recursive bool) (*[]FabricMod, error) {
	paths := make([]string, 0)
	var err error

	if recursive {
		err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, ".jar") {
				paths = append(paths, path)
			}

			return nil
		})
	} else {
		files, err := ioutil.ReadDir(dirname)
		if err == nil {
			for _, file := range files {
				paths = append(paths, path.Join(dirname, file.Name()))
			}
		}
	}

	if err != nil {
		return nil, err
	}

	mods := make([]FabricMod, 0)
	for _, modPath := range paths {
		mod, err := GetModInfo(modPath)
		if err != nil {
			continue
		}

		mods = append(mods, *mod)
	}

	return &mods, nil
}

func GetAllMods() (*[]FabricMod, error) {
	mods, err := GetMods(path.Join(WorkDir, "mods"), false)
	if err != nil {
		return nil, err
	}

	libraries, err := GetMods(path.Join(WorkDir, "libraries"), true)
	if err == nil {
		for _, library := range *libraries {
			*mods = append(*mods, library)
		}
	}

	return mods, nil
}

func GetFullModMap(mods *[]FabricMod) (map[string]FabricMod, error) {
	if mods == nil {
		var err error
		mods, err = GetAllMods()

		if err != nil {
			return nil, err
		}
	}

	modMap := make(map[string]FabricMod)
	for _, mod := range *mods {
		modMap[mod.ID] = mod
		if len(mod.Nested) > 0 {
			nestedModMap, err := GetFullModMap(&mod.Nested)
			if err != nil {
				continue
			}

			for id, nestedMod := range nestedModMap {
				modMap[id] = nestedMod
			}
		}
	}

	return modMap, nil
}
