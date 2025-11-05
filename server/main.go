package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Applicationはmanifest.json内の各アプリケーションの構造を定義します
type Application struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ArtifactType string `json:"artifact_type"`
	BuildCommand string `json:"build_command"`
	ArtifactPath string `json:"artifact_path"`
}

func listApplications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	manifestPath := "./server-files/manifest.json"
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		http.Error(w, "Failed to read manifest.json", http.StatusInternalServerError)
		return
	}
	var temp []Application
	if err := json.Unmarshal(content, &temp); err != nil {
		http.Error(w, "Failed to parse manifest.json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func getArtifact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 6 {
		http.NotFound(w, r)
		return
	}
	encodedAppName := pathParts[4]
	appName, err := url.PathUnescape(encodedAppName)
	if err != nil {
		http.Error(w, "Invalid app name", http.StatusBadRequest)
		return
	}

	content, err := os.ReadFile("./server-files/manifest.json")
	if err != nil {
		http.Error(w, "Failed to read manifest", http.StatusInternalServerError)
		return
	}
	var apps []Application
	if err := json.Unmarshal(content, &apps); err != nil {
		http.Error(w, "Failed to parse manifest", http.StatusInternalServerError)
		return
	}

	var targetApp *Application
	for i := range apps {
		if apps[i].Name == appName {
			targetApp = &apps[i]
			break
		}
	}
	if targetApp == nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	log.Printf("Building '%s' with: %s", targetApp.Name, targetApp.BuildCommand)
	cmd := exec.Command("bash", "-c", targetApp.BuildCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Build failed for '%s': %s\n%s", targetApp.Name, err, string(output))
		http.Error(w, fmt.Sprintf("Build failed: %s", string(output)), http.StatusInternalServerError)
		return
	}
	log.Printf("Build successful for '%s'", targetApp.Name)

	// Set the proper filename for download
	fileName := filepath.Base(targetApp.ArtifactPath)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	// 成果物ファイルをクライアントに返す
	http.ServeFile(w, r, targetApp.ArtifactPath)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `Server is running`)
	})
	mux.HandleFunc("/api/v1/applications", listApplications)

	mux.HandleFunc("/api/v1/applications/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/artifact") {
			getArtifact(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	port := "8080"
	log.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}
