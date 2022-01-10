package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rmil/mc-cloud-server-manager/cloud"
	"github.com/rmil/mc-cloud-server-manager/ssh"
)

// Config represents the global application configuration.
type Config struct {
	Cloud cloud.Config `json:"cloud"`
	SSH   ssh.Config   `json:"ssh"`
}

var defaultConfig = Config{
	Cloud: cloud.Config{
		Token: "Hetzner API token",
	},
	SSH: ssh.Config{
		Hostname:    "remote server hostname or IP address",
		Username:    "remote server user name",
		KeyFilePath: "SSH key file-path",
	},
}

// LoadConfig loads the global configuration from a path.
func LoadConfig(filePath string) (Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}

	config := Config{}
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	config.Cloud.AppName = "mc-cloud-server-manager"
	config.Cloud.AppVersion = "v0.1.0"
	return config, nil
}

// GenerateDefaultconfigFile generates a config file with the required fields.
func GenerateDefaultConfigFile(filePath string) error {
	jsonString, _ := json.MarshalIndent(defaultConfig, "", "    ")
	err := os.WriteFile(filePath, jsonString, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}
	return nil
}
