package config

type Config struct {
	WorkDir          string
	MinecraftVersion string
	CheckMinecraft   bool
}

var Global = Config{}
