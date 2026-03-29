package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string //`json:"db_url"`
	CurrentUserName string //`json:"current_user_name"`
}

type State struct {
	Config *Config
}

const configFileName = ".gatorconfig.json"

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil

}

func Read() (Config, error) {

	fullPath, err := GetConfigPath()
	if err != nil {
		return Config{}, err

	}
	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	c := Config{}
	err = decoder.Decode(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil

}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return write(*c)

}

func write(c Config) error {
	fullPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}
	return nil
}
