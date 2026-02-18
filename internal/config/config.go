package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFile := homeDir + configFileName

	return configFile, nil
}

func write(cfg Config) error {
	jsonFile, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(jsonFile)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFile, err := getConfigFilePath()
	if err != nil {
		fmt.Println(err)
		return Config{}, err
	}

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	byteVal, _ := io.ReadAll(file)

	var cfg Config
	err = json.Unmarshal(byteVal, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.Current_user_name = name
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}
