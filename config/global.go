package config

type Config struct {
	WorkDir          string
	MinecraftVersion string
	CheckMinecraft   bool
	Force            bool
}

var Global = Config{}
