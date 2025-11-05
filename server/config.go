package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultConfigFilename = "config.json"

type Config struct {
	ServerAddress string `json:"server_address"`
}

func (a *App) GetServerAddress() string {
	return a.config.ServerAddress
}

func (a *App) SetServerAddress(address string) error {
	a.config.ServerAddress = address
	return a.saveConfig()
}

func (a *App) loadConfig() error {
	data, err := os.ReadFile(a.configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &a.config)
}

func (a *App) saveConfig() error {
	data, err := json.MarshalIndent(a.config, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(a.configPath), 0750); err != nil {
		return err
	}
	return os.WriteFile(a.configPath, data, 0644)
}
