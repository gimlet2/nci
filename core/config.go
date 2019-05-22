package core

import (
	"os"
)

type Config struct {
	Hostname string
	SSLCacheDir string
	ClientId string
	ClientSecret string
}

func GetConfig() *Config {
	return &Config {
		Hostname: os.Getenv("NCI_HOSTNAME"),
		SSLCacheDir: os.Getenv("NCI_SSL_CACHE_DIR"),
		ClientId: os.Getenv("NCI_CLIENT_ID"),
		ClientSecret: os.Getenv("NCI_CLIENT_SECRET"),
	}
}