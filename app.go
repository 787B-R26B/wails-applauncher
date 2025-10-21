package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const serverBaseURL = "http://localhost:8080/"

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
	resp, err := http.Get(serverBaseURL + "manifest.json")
	if err != nil {
		return "", fmt.Errorf("failed to fetch manifest: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read manifest body: %w", err)
	}
	return string(body), nil
}

func buildLinuxTerminalCmd(workDir, commandToRun string) *exec.Cmd {
	// The command to run, with a pause at the end to keep the window open.
	fullCommand := fmt.Sprintf("cd %s && %s; echo; echo 'Script finished. Press Enter to close.'; read", workDir, commandToRun)

	// Terminal emulators and their command-line argument formats.
	terminals := []struct {
		path string
		args func(string) []string
	}{
		{"gnome-terminal", func(c string) []string { return []string{"--", "bash", "-c", c} }},
		{"konsole", func(c string) []string { return []string{"-e", "bash", "-c", c} }},
		{"xfce4-terminal", func(c string) []string { return []string{"--command", fmt.Sprintf("bash -c '%s'", c)} }},
		{"xterm", func(c string) []string { return []string{"-e", "bash", "-c", c} }},
	}

	for _, t := range terminals {
		if path, err := exec.LookPath(t.path); err == nil {
			args := t.args(fullCommand)
			return exec.Command(path, args...)
		}
	}
	return nil // No supported terminal found
}

func downloadScript(filename string) (string, error) {
	resp, err := http.Get(serverBaseURL + filename)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	tempFile, err := os.CreateTemp("", "launcher-*"+filepath.Ext(filename))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func (a *App) ExecuteScriptInTerminal(language string, filename string) error {
	var commandToRun string

	switch language {
	case "shell":
		commandToRun = filename
	case "python", "ruby":
		tempFilePath, err := downloadScript(filename)
		if err != nil {
			return fmt.Errorf("failed to download script: %w", err)
		}
		if language == "python" {
			commandToRun = fmt.Sprintf("python3 -u %s", tempFilePath)
		} else {
			commandToRun = fmt.Sprintf("ruby %s", tempFilePath)
		}
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}

	var cmd *exec.Cmd
	workDir, _ := os.Getwd()

	switch runtime.GOOS {
	case "darwin":
		appleScript := fmt.Sprintf(`tell app "Terminal" to do script "cd %s && %s"`, workDir, commandToRun)
		cmd = exec.Command("osascript", "-e", appleScript)
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", "start", "cmd.exe", "/K", commandToRun)
	case "linux":
		cmd = buildLinuxTerminalCmd(workDir, commandToRun)
		if cmd == nil {
			return fmt.Errorf("could not find a supported terminal on Linux (gnome-terminal, konsole, xfce4-terminal, xterm)")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
