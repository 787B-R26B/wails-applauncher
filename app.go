package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const defaultConfigFilename = "config.json"

type Config struct {
	ServerAddress string `json:"server_address"`
}

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

func buildLinuxTerminalCmd(workDir, commandToRun string) *exec.Cmd {
	// The command to run, with a pause at the end to keep the window open.
	fullCommand := fmt.Sprintf("cd %s && %s; echo; echo 'Script finished. Press Enter to close.'; read", workDir, commandToRun)

	// Terminal emulators and their command-line argument formats.
	terminals := []struct {
		path string
		args func(string) []string
	}{
		{"gnome-terminal", func(c string) []string { return []string{"--", "bash", "-c", c} }},
		{"konsole", func(c string) []string { return []string{"-", "e", "bash", "-c", c} }},
		{"xfce4-terminal", func(c string) []string { return []string{"--command", fmt.Sprintf("bash -c '%s'", c)} }},
		{"xterm", func(c string) []string { return []string{"-", "e", "bash", "-c", c} }},
	}

	for _, t := range terminals {
		if path, err := exec.LookPath(t.path); err == nil {
			args := t.args(fullCommand)
			return exec.Command(path, args...)
		}
	}
	return nil // No supported terminal found
}

func (a *App) downloadScript(filename string) (string, error) {
	resp, err := http.Get(a.config.ServerAddress + filename)
	if err != nil {
		return "", fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status from server: %s", resp.Status)
	}

	tempFile, err := os.CreateTemp("", "launcher-*"+filepath.Ext(filename))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	bytesCopied, err := io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	if bytesCopied == 0 {
		return "", fmt.Errorf("downloaded file is empty")
	}

	return tempFile.Name(), nil
}

func (a *App) ExecuteScriptInTerminal(language string, filename string) error {
	var commandToRun string

	switch language {
	case "shell":
		commandToRun = filename
	case "python", "ruby":
		tempFilePath, err := a.downloadScript(filename)
		if err != nil {
			return fmt.Errorf("failed to download script: %w", err)
		}
		if language == "python" {
			commandToRun = fmt.Sprintf("python3 -u %s", tempFilePath)
		} else {
			commandToRun = fmt.Sprintf("ruby %s", tempFilePath)
		}
	case "c":
		if _, err := exec.LookPath("gcc"); err != nil {
			return fmt.Errorf("gcc not found. Please install gcc to compile C code")
		}
		sourcePath, err := a.downloadScript(filename)
		if err != nil {
			return fmt.Errorf("failed to download c source: %w", err)
		}
		// Assume gcc is installed.
		outputPath := sourcePath + ".out"
		if runtime.GOOS == "windows" {
			outputPath = sourcePath + ".exe"
		}

		cmd := exec.Command("gcc", sourcePath, "-o", outputPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to compile c code: %s, err: %w", string(output), err)
		}
		if err := os.Chmod(outputPath, 0755); err != nil {
			return fmt.Errorf("failed to make c output executable: %w", err)
		}
		commandToRun = outputPath
	case "binary":
		tempFilePath, err := a.downloadScript(filename)
		if err != nil {
			return fmt.Errorf("failed to download binary: %w", err)
		}
		if err := os.Chmod(tempFilePath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}
		commandToRun = tempFilePath
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
