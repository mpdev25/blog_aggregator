package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
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
		//return nil, fmt.Errorf("failed to find home directory: %w", err)
	}
	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	//data, err := os.ReadFile(path)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read config file at %s: %w", path, err)
	//}
	decoder := json.NewDecoder(file)
	c := Config{}
	err = decoder.Decode(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
	//var config Config
	//err = json.Unmarshal(data, &config)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", configFileName, err)
	//}
	//return &config, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return write(*c)
	//file, err := os.Create(configFileName)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//encoder := json.NewEncoder(file)
	//encoder.SetIndent("", " ")
	//return encoder.Encode(c)
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
