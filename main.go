package main

import (
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("", "Request-URL", r.URL.String())
		if strings.Contains(r.URL.Path, "favicon.ico") {
			icon, err := os.ReadFile("./frontend-dist/favicon.ico")
			if err != nil {
				http.Error(w, "Not Found", 404)
				return
			}
			w.Write(icon)
			return
		} else if strings.Contains(r.URL.Path, "Marker.svg") {
			file, err := os.ReadFile("./frontend-dist/Marker.svg")
			if err != nil {
				http.Error(w, "Not Found", 404)
				return
			}
			w.Write(file)
			return
		} else {
			re := regexp.MustCompile(`.*(/assets/.*)`)
			match := re.FindStringSubmatch(r.URL.Path)
			if len(match) > 1 {
				index, err := os.ReadFile("./frontend-dist" + match[1])
				if err != nil {
					http.Error(w, "Not Found", 404)
					return
				}
				if strings.Contains(match[1], ".js") {
					w.Header().Set("Content-Type", "text/javascript")
				} else if strings.Contains(match[1], ".css") {
					w.Header().Set("Content-Type", "text/css")
				}
				w.Write(index)
			} else {
				index, err := os.ReadFile("./frontend-dist/index.html")
				if err != nil {
					http.Error(w, "Not Found", 404)
					return
				}
				w.Header().Set("Content-Type", "text/html")
				w.Write(index)
			}
		}
	}))

	srv := &http.Server{
		Addr:    ":5173",
		Handler: nil,
	}

	slog.Info("Serving on port 5173")
	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
