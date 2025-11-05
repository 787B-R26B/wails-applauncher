package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// listApplicationsは manifest.json を読み込み、その内容をJSONとして返します。
func listApplications(w http.ResponseWriter, r *http.Request) {
	// CORSを許可（開発中にクライアントからアクセスするため）
	w.Header().Set("Access-Control-Allow-Origin", "*")

	manifestPath := "./server-files/manifest.json"
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		http.Error(w, "Failed to read manifest.json", http.StatusInternalServerError)
		log.Printf("Error reading manifest file %s: %v", manifestPath, err)
		return
	}

	// 内容が有効なJSONかだけ確認し、そのままクライアントに返す
	var temp interface{}
	if err := json.Unmarshal(content, &temp); err != nil {
		http.Error(w, "Failed to parse manifest.json", http.StatusInternalServerError)
		log.Printf("Error parsing manifest file: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func main() {
	http.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `Server is running`)
	})

	// アプリケーション一覧を返す新しいエンドポイント
	http.HandleFunc("/api/v1/applications", listApplications)

	port := "8080"
	log.Printf("Starting server on http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %s\n", err)
	}
}
