package remote

type Source interface {
	GetModInfo(name string) (Mod, error)
}

type Mod interface {
	GetName() string
	GetDescription() string
	GetVersions() (*[]ModVersion, error)
}

type ModVersion interface {
	GetName() string
	GetFileHash() string
	GetDownloadURL() string
}
