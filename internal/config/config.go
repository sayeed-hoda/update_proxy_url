package config

import (
	"fmt"
	"strings"
)

type Config struct {
	MongoURI          string
	MongoDataBase     string
	CollectionName    string
	CollectionKeyName []string
	Domain            string
}

func (c *Config) Validate() error {
	if c.Domain == "" {
		return fmt.Errorf("DOMAIN environment variable is not set")
	}
	if c.MongoDataBase == "" {
		return fmt.Errorf("MONGO_DATABASE environment variable is not set")
	}
	if c.MongoURI == "" {
		return fmt.Errorf("MONGO_URI environment variable is not set")
	}
	if c.CollectionName == "" {
		return fmt.Errorf("COLLECTION_NAME environment variable is not set")
	}
	if len(c.CollectionKeyName) == 0 {
		return fmt.Errorf("COLLECTION_KEY_NAME environment variable is not set or is empty")
	}

	return nil
}

func LoadEnv(getenv func(string) string) (*Config, error) {

	config := Config{}

	config.Domain = getEnvString(getenv, "DOMAIN", "")

	config.MongoURI = getEnvString(getenv, "MONGO_URI", "")
	config.MongoDataBase = getEnvString(getenv, "MONGO_DATABASE", "")
	config.CollectionName = getEnvString(getenv, "COLLECTION_NAME", "")
	config.CollectionKeyName = getEnvSlice(getenv, "COLLECTION_KEY_NAME", []string{})

	err := config.Validate()

	return &config, err

}

func getEnvSlice(getenv func(string) string, key string, fallback []string) []string {
	value := getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

func getEnvString(getenv func(string) string, key string, fallback string) string {
	value := getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
