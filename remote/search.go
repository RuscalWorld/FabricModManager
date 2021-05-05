package remote

func FindMod(name string) (Mod, error) {
	modrinth := Modrinth{URL: "https://api.modrinth.com"}
	info, err := modrinth.GetModInfo(name)
	return info, err
}
