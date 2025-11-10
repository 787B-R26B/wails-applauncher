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

type App struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ArtifactType string `json:"artifact_type"`
	BuildCommand string `json:"build_command"`
	ArtifactPath string `json:"artifact_path"`
	RunCommand   string `json:"run_command"`
}

const (
	serverFilesDir = "./server-files"
	manifestPath   = serverFilesDir + "/manifest.json"
)

var (
	appList []App
)

// --- Handlers ---

func listApplications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appList)
}

func getArtifact(w http.ResponseWriter, r *http.Request) {
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

	var targetApp *App
	for i := range appList {
		if appList[i].Name == appName {
			targetApp = &appList[i]
			break
		}
	}
	if targetApp == nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	log.Printf("Building '%s' with: %s", targetApp.Name, targetApp.BuildCommand)
	cmd := exec.Command("bash", "-c", targetApp.BuildCommand)
	cmd.Dir = serverFilesDir // Set working directory for the build command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Build failed for '%s': %s\n%s", targetApp.Name, err, string(output))
		http.Error(w, fmt.Sprintf("Build failed: %s", string(output)), http.StatusInternalServerError)
		return
	}
	log.Printf("Build successful for '%s'", targetApp.Name)

	artifactAbsPath := filepath.Join(serverFilesDir, targetApp.ArtifactPath)
	fileName := filepath.Base(artifactAbsPath)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	http.ServeFile(w, r, artifactAbsPath)
}

// --- CORS Middleware ---

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// --- Main ---

func main() {
	// Load manifest at startup
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Fatalf("Failed to read manifest file at %s: %s", manifestPath, err)
	}
	if err := json.Unmarshal(content, &appList); err != nil {
		log.Fatalf("Failed to parse manifest file: %s", err)
	}
	log.Printf("Loaded %d applications from manifest.", len(appList))

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

	httpHandler := corsMiddleware(mux)

	port := "8080"
	log.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, httpHandler); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}
