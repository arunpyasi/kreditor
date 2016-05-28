package config

import (
	"os"
)

type Configuration struct {
	DatabaseURI   string
	ListenAddress string
	Secret        string
	Debug         bool
}

var C = Configuration{}

func getEnv(variable string) string {
	prefix := "KREDITOR_"
	return os.Getenv(prefix + variable)
}

func init() {
	C.DatabaseURI = getEnv("DATABASE_URI")
	C.Secret = getEnv("SECRET")
	C.ListenAddress = getEnv("LISTEN_ADDRESS")

	if getEnv("DEBUG") != "" {
		C.Debug = false
	}

}
