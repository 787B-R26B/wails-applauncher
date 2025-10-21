package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func (a *App) ExecuteScriptInTerminal(language string, filename string) error {
	var commandToRun string

	switch language {
	case "shell":
		commandToRun = filename
	case "python":
		ScriptPath, _ := filepath.Abs(filepath.Join("server-files", filename))
		commandToRun = fmt.Sprintf("python3 -u %s", ScriptPath)
	case "ruby":
		ScriptPath, _ := filepath.Abs(filepath.Join("server-files", filename))
		commandToRun = fmt.Sprintf("ruby %s", ScriptPath)
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
