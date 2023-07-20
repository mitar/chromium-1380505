package main

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	gddo "github.com/golang/gddo/httputil"
)

func handleHTML(w http.ResponseWriter, req *http.Request) {
	log.Printf("Serving HTML at %s", req.URL.Path)
	w.Header().Set("Content-Type", "text/html")
	data, err := os.ReadFile("main.html")
	if err != nil {
		panic(err)
	}
	hash := sha256.New()
	_, _ = hash.Write(data)
	etag := `"` + base64.RawURLEncoding.EncodeToString(hash.Sum(nil)) + `"`
	w.Header().Set("ETag", etag)
	http.ServeFile(w, req, "main.html")
}

func handleJSON(w http.ResponseWriter, req *http.Request) {
	log.Printf("Serving JSON at %s", req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	data, err := os.ReadFile("main.json")
	if err != nil {
		panic(err)
	}
	hash := sha256.New()
	_, _ = hash.Write(data)
	etag := `"` + base64.RawURLEncoding.EncodeToString(hash.Sum(nil)) + `"`
	w.Header().Set("ETag", etag)
	http.ServeFile(w, req, "main.json")
}

func handle(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" && req.URL.Path != "/data.json" {
		http.NotFound(w, req)
		return
	}

	w.Header().Add("Vary", "Accept")
	w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	contentType := gddo.NegotiateContentType(req, []string{"text/html", "application/json"}, "text/html")
	switch contentType {
	case "text/html":
		handleHTML(w, req)
	case "application/json":
		handleJSON(w, req)
	}
}

func main() {
	log.Print("Listening on :8000")
	err := http.ListenAndServe(":8000", http.HandlerFunc(handle))
	log.Printf("%s", err)
}
