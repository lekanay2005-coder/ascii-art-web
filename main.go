package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error reading form", http.StatusBadRequest)
		return
	}
	text := r.FormValue("text")
	banner := r.FormValue("banner")
	if text == "" || banner == "" {
		http.Error(w, "Missing text or banner", http.StatusBadRequest)
		return
	}
	bannerFile := fmt.Sprintf("banners/%s.txt", banner)
	data, err := os.ReadFile(bannerFile)
	if err != nil {
		http.Error(w, "Banner file not found", http.StatusNotFound)
		return
	}
	asciiArt := generateAsciiArt(text, string(data))

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	_ = tmpl.Execute(w, asciiArt) // Render ASCII art in template
}