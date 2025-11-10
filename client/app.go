package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SaveAndRunArtifact saves the artifact, and then executes it in a new terminal window.
func (a *App) SaveAndRunArtifact(isZip bool, fileData []byte, runCommand string) (string, error) {
	tempDir, err := os.MkdirTemp("", "wails-app-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	// Not cleaning up tempDir so the user can see the files if needed.

	var commandToRunInTerminal string
	var workingDir string

	if isZip {
		// Unzip logic
		zipReader, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))
		if err != nil {
			return "", fmt.Errorf("failed to read zip data: %w", err)
		}
		for _, file := range zipReader.File {
			extractedFilePath := filepath.Join(tempDir, file.Name)
			if file.FileInfo().IsDir() {
				os.MkdirAll(extractedFilePath, file.Mode())
				continue
			}
			fileReader, err := file.Open()
			if err != nil {
				return "", err
			}
			defer fileReader.Close()
			targetFile, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return "", err
			}
			defer targetFile.Close()
			if _, err := io.Copy(targetFile, fileReader); err != nil {
				return "", err
			}
		}
		workingDir = tempDir
		commandToRunInTerminal = runCommand
	} else {
		// Raw binary logic
		executablePath := filepath.Join(tempDir, runCommand)
		if runtime.GOOS == "windows" {
			executablePath += ".exe"
		}
		err = os.WriteFile(executablePath, fileData, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to write executable: %w", err)
		}
		workingDir = ""
		commandToRunInTerminal = executablePath
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		var script string
		if workingDir != "" {
			script = fmt.Sprintf("cd '%s' && %s", workingDir, commandToRunInTerminal)
		} else {
			script = commandToRunInTerminal
		}
		// Safer way to build the AppleScript command to avoid escape sequence errors.
		escapedScript := strings.ReplaceAll(script, `"`, `\\\"`)
		appleScript := `tell application "Terminal" to do script "` + escapedScript + `"`
		cmd = exec.Command("osascript", "-e", appleScript)

	case "windows":
		cmd = exec.Command("cmd", "/C", "start", "cmd.exe", "/K", commandToRunInTerminal)
		if workingDir != "" {
			cmd.Dir = workingDir
		}

	case "linux":
		var script string
		// Append '; exec bash' to keep the terminal open after the command finishes.
		if workingDir != "" {
			script = fmt.Sprintf("cd '%s' && %s; exec bash", workingDir, commandToRunInTerminal)
		} else {
			script = fmt.Sprintf("%s; exec bash", commandToRunInTerminal)
		}
		cmd = exec.Command("x-terminal-emulator", "-e", "bash", "-c", script)

		if workingDir != "" {
			cmd.Dir = workingDir
		}

	default:
		return "", fmt.Errorf("unsupported OS for terminal execution: %s", runtime.GOOS)
	}

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to launch terminal: %w", err)
	}

	return fmt.Sprintf("Launched '%s' in a new terminal.", runCommand), nil
}
