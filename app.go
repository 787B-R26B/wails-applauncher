package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetScriptManifest() (string, error) {
	content, err := os.ReadFile("server-files/manifest.json")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (a *App) ExecuteScript(language string, filename string) (string, error) {
	var cmd *exec.Cmd

	switch language {
	case "shell":
		parts := strings.Split(filename, " ")
		head := parts[0]
		parts = parts[1:]
		cmd = exec.Command(head, parts...)
	case "python":
		ScriptPath, err := filepath.Abs(filepath.Join("server-files", filename))
		if err != nil {
			return "", fmt.Errorf("could not get absolute path for script: %w", err)
		}
		cmd = exec.Command("python3", ScriptPath)
	case "ruby":
		ScriptPath, err := filepath.Abs(filepath.Join("server-files", filename))
		if err != nil {
			return "", fmt.Errorf("could not get absolute path for script: %w", err)
		}
		cmd = exec.Command(language, ScriptPath)

	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
