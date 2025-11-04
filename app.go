package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx        context.Context
	config     Config
	configPath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to current directory
		configDir = "."
	}
	a.configPath = filepath.Join(configDir, "wails-applauncher", defaultConfigFilename)

	if err := a.loadConfig(); err != nil {
		// If config doesn't exist, create it with default
		a.config.ServerAddress = "http://localhost:8080/"
		if err := a.saveConfig(); err != nil {
			// Log this error, but the app can still run with the default.
			fmt.Fprintf(os.Stderr, "Could not save default config: %v\n", err)
		}
	}
}

func (a *App) GetScriptManifest() (string, error) {
	resp, err := http.Get(a.config.ServerAddress + "manifest.json")
	if err != nil {
		return "", fmt.Errorf("failed to fetch manifest: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read manifest body: %w", err)
	}
	return string(body), nil
}
