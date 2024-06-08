package config

import (
	"log"
	"os"
)

func GetAwsRegion() string {
	return getEnvironmentValue("AWS_REGION")
}

func GetAwsAccessKey() string {
	return getEnvironmentValue("AWS_ACCESS_KEY_ID")
}
func GetBucketName() string { return getEnvironmentValue("AWS_BUCKET_NAME") }
func GetAwsSecretAccessKey() string {
	return getEnvironmentValue("AWS_SECRET_ACCESS_KEY")
}
func GetDataSourceUrl() string { return getEnvironmentValue("DATA_SOURCE_URL") }
func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}
	return os.Getenv(key)
}
