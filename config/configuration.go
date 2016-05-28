package config

import (
	"os"
)

type Configuration struct {
	DatabaseURI string
	Secret      string
}

var C = Configuration{}

func getEnv(variable string) string {
	prefix := "KREDITOR_"
	return os.Getenv(prefix + variable)
}

func init() {

	C.DatabaseURI = getEnv("DATABASE_URI")
	C.Secret = getEnv("SECRET")

}
