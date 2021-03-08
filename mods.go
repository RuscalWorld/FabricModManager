package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type FabricMod struct {
	Environment string      `json:"environment"`
	Depends     interface{} `json:"depends"`
	Recommends  interface{} `json:"recommends"`
	Suggests    interface{} `json:"suggests"`
	Breaks      interface{} `json:"breaks"`
	Conflicts   interface{} `json:"conflicts"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
}

func (m FabricMod) GetDepends() []string {
	return StringOrArray(m.Depends)
}

func (m FabricMod) GetRecommends() []string {
	return StringOrArray(m.Recommends)
}

func (m FabricMod) GetSuggests() []string {
	return StringOrArray(m.Suggests)
}

func (m FabricMod) GetBreaks() []string {
	return StringOrArray(m.Breaks)
}

func (m FabricMod) GetConflicts() []string {
	return StringOrArray(m.Conflicts)
}

func StringOrArray(input interface{}) []string {
	str, ok := input.(string)
	if ok {
		return []string{str}
	}
	return input.([]string)
}

func GetModInfo(path string) (*FabricMod, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	for _, file := range reader.File {
		if file.Name == "fabric.mod.json" {
			file, err := file.Open()
			if err != nil {
				return nil, err
			}

			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}

			info := &FabricMod{}
			err = json.Unmarshal(bytes, info)
			if err != nil {
				return nil, err
			}

			return info, nil
		}
	}

	return nil, errors.New("Input JAR file wasn't a Fabric mod ")
}

func GetMods(dirname string) (*[]FabricMod, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	mods := make([]FabricMod, 0)
	for _, file := range files {
		path, err := filepath.Abs("mods/" + file.Name())
		if err != nil {
			continue
		}

		mod, err := GetModInfo(path)
		if err != nil {
			fmt.Println("Unable to read", file.Name()+":", err)
			continue
		}

		mods = append(mods, *mod)
	}

	return &mods, nil
}
