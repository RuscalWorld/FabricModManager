package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/RuscalWorld/FabricModManager/config"
)

type Profile struct {
	Created       time.Time `json:"created"`
	LastUsed      time.Time `json:"lastUsed"`
	LastVersionID string    `json:"lastVersionId"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
}

type LauncherProfiles struct {
	Profiles map[string]Profile `json:"profiles"`
}

func (p LauncherProfiles) GetLastUsedFabricProfile() *Profile {
	var latest Profile
	first := true

	for _, profile := range p.Profiles {
		if !strings.Contains(profile.LastVersionID, "fabric-loader") {
			continue
		}

		if first || profile.LastUsed.After(latest.LastUsed) {
			first = false
			latest = profile
		}
	}

	return &latest
}

func GetLauncherProfiles() (*LauncherProfiles, error) {
	file, err := ioutil.ReadFile(path.Join(config.Global.WorkDir, "launcher_profiles.json"))
	if err != nil {
		return nil, err
	}

	launcherProfiles := &LauncherProfiles{}
	err = json.Unmarshal(file, launcherProfiles)
	return launcherProfiles, err
}

func GetCurrentMinecraftVersion() (string, error) {
	profiles, err := GetLauncherProfiles()
	if err != nil {
		return "", fmt.Errorf("cannot read launcher_profiles.json: %s", err)
	}

	profile := profiles.GetLastUsedFabricProfile()
	if profile == nil {
		return "", fmt.Errorf("cannot find fabric installation in launcher_profiles.json")
	}

	// It looks like fabric-loader-x.xx.x-x.xx.x, where latest part is Minecraft version
	parts := strings.Split(profile.LastVersionID, "-")
	if len(parts) < 4 {
		return "", fmt.Errorf("profile version id has incorrect name (%s)", profile.LastVersionID)
	}

	return parts[3], nil
}
