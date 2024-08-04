package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
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

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request-URL %s", r.URL.String())
		now := float64(time.Now().UnixNano())
		defer func() { fmt.Printf(" (%f µs) \n", (float64(time.Now().UnixNano())-now)/1000.0) }()

		re := regexp.MustCompile(`.*(/assets/.*)`)
		match := re.FindStringSubmatch(r.URL.Path)
		if len(match) > 1 {
			file, err := os.Open("./frontend-dist" + match[1])
			if err != nil {
				index_html(w)
				return
			}
			if strings.Contains(match[1], ".js") {
				w.Header().Set("Content-Type", "text/javascript")
			} else if strings.Contains(match[1], ".css") {
				w.Header().Set("Content-Type", "text/css")
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
	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
