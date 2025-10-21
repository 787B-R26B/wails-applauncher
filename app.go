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
		linuxCommand := fmt.Sprintf("%s; echo; echo \"Script finished. Press Enter to close.\"; read", commandToRun)
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", linuxCommand)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
