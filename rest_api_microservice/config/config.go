package config

import (
	"log"
	"os"
)

func GetJaegerUrl() string     { return getEnvironmentValue("DATA_JAEGER_URL") }
func GetDataSourceUrl() string { return getEnvironmentValue("DATA_SOURCE_URL") }
func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}
	return os.Getenv(key)
}
