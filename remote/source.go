package remote

type Source interface {
	GetModInfo(name string) (Mod, error)
	GetDisplayName() string
}

type Mod interface {
	GetName() string
	GetDescription() string
	GetVersions() (*[]ModVersion, error)
	GetRemote() Source
}

type ModVersion interface {
	GetName() string
	GetFileHash() string
	GetDownloadURL() string
}
