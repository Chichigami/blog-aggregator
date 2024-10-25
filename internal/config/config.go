package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func GetConfigFilePath() (string, error) {
	//var path string
	cfgPath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	//filePath += "/" + configFileName
	path := filepath.Join(cfgPath, configFileName) //prob safer
	return path, nil
}

func ReadAndParse(filePath string) Config {
	body, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	var dbConfig Config
	err = json.Unmarshal(body, &dbConfig)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(dbConfig)
	return dbConfig
}

func write(path string, cfg *Config) error {
	body, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) SetUser(name string) {
	cfg.CurrentUserName = name
	path, err := GetConfigFilePath()
	if err != nil {
		fmt.Println(err)
	}
	write(path, cfg)
}

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}
