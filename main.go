package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func index_html(w http.ResponseWriter) {
	index, err := os.Open("./frontend-dist/index.html")
	if err != nil {
		http.Error(w, "Not Found", 404)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, index)
}

type MimeType struct {
	Name      string
	Mime_type string
	Ext       string
	Details   string
}

func main() {
	data, err := os.ReadFile("mime_types.json")
	if err != nil {
		log.Fatalf("failed to read mime_types.json: %v", err)
	}
	mime_types := []MimeType{}
	err = json.Unmarshal(data, &mime_types)
	if err != nil {
		log.Fatalf("failed to unmarshal mime_types.json: %v", err)
	}

	mime_map := make(map[string]string)
	for _, mime := range mime_types {
		mime_map[mime.Ext] = mime.Mime_type
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request-URL: %s\n", r.URL.String())
		now := float64(time.Now().UnixNano())
		defer func() { fmt.Printf(" (%f µs) \n", (float64(time.Now().UnixNano())-now)/1000.0) }()

		if r.URL.Path != "/" && r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
			return
		}

		if r.URL.Path == "/" {
			index_html(w)
			return
		}

		assetPath := "./frontend-dist" + r.URL.Path
		if _, err := os.Stat(assetPath); err == nil {
			ext := filepath.Ext(r.URL.Path)
			mime_type := mime_map[ext]
			if mime_type == "" {
				mime_type = "text/plain"
			}
			w.Header().Set("Content-Type", mime_type)
			http.ServeFile(w, r, assetPath)
			return
		}

		index_html(w)
	}))

	srv := &http.Server{
		Addr:    ":5173",
		Handler: nil,
	}

	log.Println("Serving on port 5173")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
