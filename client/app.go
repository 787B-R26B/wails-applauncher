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

// SaveAndRunArtifact receives file data and a run command, saves the data, prepares, and executes the command.
func (a *App) SaveAndRunArtifact(isZip bool, fileData []byte, runCommand string) (string, error) {
	tempDir, err := os.MkdirTemp("", "wails-app-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	var cmd *exec.Cmd

	if isZip {
		// Unzip all files into tempDir
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
				return "", fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer fileReader.Close()
			targetFile, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return "", fmt.Errorf("failed to create target file: %w", err)
			}
			defer targetFile.Close()
			if _, err := io.Copy(targetFile, fileReader); err != nil {
				return "", fmt.Errorf("failed to copy from zip: %w", err)
			}
		}

		// For zips, the runCommand is executed inside the tempDir
		cmdArgs := strings.Fields(runCommand)
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = tempDir
	} else {
		// It's a raw binary. The `runCommand` is the desired filename.
		executablePath := filepath.Join(tempDir, runCommand)
		if runtime.GOOS == "windows" {
			executablePath += ".exe"
		}
		err = os.WriteFile(executablePath, fileData, 0755) // 0755 = rwxr-xr-x
		if err != nil {
			return "", fmt.Errorf("failed to write executable: %w", err)
		}

		// The command is the direct path to the saved binary
		cmd = exec.Command(executablePath)
	}

	// Execute the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
