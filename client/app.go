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

// SaveAndRunArtifact receives file data and a run command, saves the data, prepares, and executes the command.
func (a *App) SaveAndRunArtifact(isZip bool, fileData []byte, runCommand string) (string, error) {
	// Create a temporary directory for the artifact
	tempDir, err := os.MkdirTemp("", "wails-app-artifact-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// This will be the root directory containing the executable/script
	executionRoot := tempDir

	if isZip {
		zipReader, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))
		if err != nil {
			return "", fmt.Errorf("failed to read zip data: %w", err)
		}

		// Extract all files from the zip into the temp directory
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
				return "", fmt.Errorf("failed to copy file content: %w", err)
			}
		}
	} else {
		// It's a single binary file. We place it inside the tempDir.
		// The name of the binary inside the temp dir doesn't matter as much.
		executablePath := filepath.Join(tempDir, "executable")
		err = os.WriteFile(executablePath, fileData, 0755) // 0755 makes it executable
		if err != nil {
			return "", fmt.Errorf("failed to write executable file: %w", err)
		}
	}

	// Prepare the final command. We replace placeholders with the temp directory path.
	// We assume the server's `run_command` uses a placeholder like `{exec_root}`
	// For example: "python {exec_root}/hello.py" or "{exec_root}/c_hello"
	// This is a simplification. A better approach would be to replace the original filename.
	// Let's assume the run_command from the server is relative to the server-files dir.
	// e.g., "python ./server-files/hello.py" -> we need to run "python {tempDir}/hello.py"
	// This is getting complex. Let's try a simpler replacement for now.
	// We will assume the command is just the executable name, and we will prepend the path.

	// Let's simplify the logic: the `run_command` from the manifest is a template.
	// For a zip, the command might be `python hello.py`. We run this inside `tempDir`.
	// For a binary, the command might be `c_hello`. We run `{tempDir}/c_hello`.

	// A much simpler model:
	// The `run_command` is a template with the executable name.
	// e.g., `python hello.py` or `./c_hello`
	// The Go function will just run this command inside the temp directory.

	cmdArgs := strings.Fields(runCommand)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = executionRoot // Run the command inside the temp directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run command '%s': %w\nOutput: %s", runCommand, err, string(output))
	}

	return string(output), nil
}
