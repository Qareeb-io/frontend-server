package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
		fmt.Printf("Request-URL %s", r.URL.String())
		now := float64(time.Now().UnixNano())
		defer func() { fmt.Printf(" (%f µs) \n", (float64(time.Now().UnixNano())-now)/1000.0) }()
		if r.URL.String() == "/" {
			index_html(w)
			return
		}

		re := regexp.MustCompile(`.*(/assets/.*)`)
		match := re.FindStringSubmatch(r.URL.Path)
		if len(match) > 1 {
			file, err := os.Open("./frontend-dist" + match[1])
			if err != nil {
				index_html(w)
				return
			}

			mime_type, ok := mime_map[filepath.Ext(match[1])]
			if ok {
				w.Header().Set("Content-Type", mime_type)
			} else {
				w.Header().Set("Content-Type", "text/plain")
			}
			io.Copy(w, file)
		} else {
			re := regexp.MustCompile(`.*(/.*)`)
			match := re.FindStringSubmatch(r.URL.Path)
			if len(match) > 1 {
				file, err := os.Open("./frontend-dist" + match[1])
				if err != nil {
					index_html(w)
					return
				}
				mime_type, ok := mime_map[filepath.Ext(match[1])]
				if ok {
					w.Header().Set("Content-Type", mime_type)
				} else {
					w.Header().Set("Content-Type", "text/plain")
				}
				io.Copy(w, file)
				return
			} else {
				index_html(w)
			}
		}
	}))

	srv := &http.Server{
		Addr:    ":5173",
		Handler: nil,
	}

	slog.Info("Serving on port 5173")
	err = srv.ListenAndServe()
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
