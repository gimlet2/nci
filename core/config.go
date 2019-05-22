package core

import (
	"os"
)

type Config struct {
	Hostname string
	ClientId string
	ClientSecret string
}

func GetConfig() *Config {
	return &Config {
		Hostname: os.Getenv("HOSTNAME"),
		ClientId: os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}