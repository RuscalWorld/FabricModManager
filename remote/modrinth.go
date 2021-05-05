package remote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Modrinth struct {
	URL string
}

func (m Modrinth) GetDisplayName() string {
	return "Modrinth"
}

func (m Modrinth) Fetch(endpoint string, dest interface{}) error {
	response, err := http.Get(m.URL + "/api/v1" + endpoint)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("server returned non-success status (%s)", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, dest)
}

func (m Modrinth) GetModInfo(name string) (Mod, error) {
	mod := ModrinthMod{}
	err := m.Fetch("/mod/"+name, &mod)
	if err != nil {
		return nil, err
	}

	mod.Modrinth = m
	return mod, err
}

type ModrinthMod struct {
	ID           string   `json:"id"`
	Slug         string   `json:"slug"`
	Team         string   `json:"team"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Downloads    int      `json:"downloads"`
	Followers    int      `json:"followers"`
	Categories   []string `json:"categories"`
	Versions     []string `json:"versions"`
	IconURL      string   `json:"icon_url"`
	IssuesURL    string   `json:"issues_url"`
	SourceURL    string   `json:"source_url"`
	WikiURL      string   `json:"wiki_url"`
	DiscordURL   string   `json:"discord_url"`
	DonationURLs []string `json:"donation_urls"`
	Modrinth     Modrinth `json:"-"`
}

func (m ModrinthMod) GetRemote() Source {
	return m.Modrinth
}

func (m ModrinthMod) GetDescription() string {
	return m.Description
}

func (m ModrinthMod) GetVersions() (*[]ModVersion, error) {
	version := ModrinthModVersion{}
	err := m.Modrinth.Fetch("/version/"+m.Versions[0], &version)
	if err != nil {
		return nil, err
	}

	versions := make([]ModVersion, 0)
	versions = append(versions, version)
	return &versions, nil
}

func (m ModrinthMod) GetName() string {
	return m.Title
}

type ModrinthModVersion struct {
	ID            string         `json:"id"`
	ModID         string         `json:"mod_id"`
	AuthorID      string         `json:"author_id"`
	Featuerd      bool           `json:"featuerd"`
	Name          string         `json:"name"`
	VersionNumber string         `json:"version_number"`
	Changelog     string         `json:"changelog"`
	DatePublished time.Time      `json:"date_published"`
	Downloads     int            `json:"downloads"`
	VersionType   string         `json:"version_type"`
	Files         []ModrinthFile `json:"files"`
	GameVersions  []string       `json:"game_versions"`
	Loaders       []string       `json:"loaders"`
}

func (m ModrinthModVersion) GetName() string {
	return m.Name
}

func (m ModrinthModVersion) GetFileHash() string {
	return m.Files[0].Hashes.SHA512
}

func (m ModrinthModVersion) GetDownloadURL() string {
	return m.Files[0].URL
}

type ModrinthFile struct {
	Hashes   ModrinthHashes `json:"hashes"`
	URL      string         `json:"url"`
	FileName string         `json:"filename"`
	Primary  bool           `json:"primary"`
}

type ModrinthHashes struct {
	SHA1   string `json:"sha_1"`
	SHA512 string `json:"sha_512"`
}
