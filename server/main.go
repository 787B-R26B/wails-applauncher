package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// APIエンドポイントのサンプル
	http.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running")
	})

	// TODO: 管理画面用のWebページを配信するハンドラを追加
	// http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("./admin-ui/dist"))))

	// TODO: アプリケーションロジック（zip化、コンパイル等）のハンドラを追加

	port := "8080"
	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
