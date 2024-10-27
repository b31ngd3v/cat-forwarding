package client

type config struct {
	ServerAddr string
	Version    string
}

func GetConfig() *config {
	return &config{
		ServerAddr: "cf.b31ngd3v.xyz:1337",
		Version:    "0.1",
	}
}
