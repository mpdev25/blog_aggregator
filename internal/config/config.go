package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil

}

func Read() (*Config, error) {
	//homeDir, err := os.UserHomeDir()
	path, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find home directory: %w", err)
	}
	//configPath := filepath.Join(homeDir, ".gatorconfig.json")
	//data, err := os.ReadFile(configPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %w", path, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", configFileName, err)
	}
	return &config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(c)
}
