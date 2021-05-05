package config

type Config struct {
	WorkDir                 string
	MinecraftVersion        string
	CheckMinecraft          bool
	Force                   bool
	IgnoreIncompatibilities bool
	NoDependencyChecks      bool
}

var Global = Config{}
