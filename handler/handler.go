package handler

import (
	"fmt"
	"net/http"

	"github.com/coderjojo/gourlshortner/shortner"
	"github.com/gorilla/mux"
)

// URL Shortner Handler

func HandleShortnerURL(w http.ResponseWriter, r *http.Request, shortner *shortner.UrlShortner) {
	orignalURL := r.FormValue("url")
	if orignalURL == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	shortUrl, err := shortner.ShortenURL(orignalURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	shortner.Logger.Printf("URL shortned: %s", shortUrl)
	shortner.Logger.Println("shortner URL handler called")

	fmt.Fprintf(w, "Shortened URL: %s", shortUrl)

}

// Redirect to orignalURL
func HandleRedirect(w http.ResponseWriter, r *http.Request, shortner *shortner.UrlShortner) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	// Check if URL present or expired
	orignalURL, err := shortner.Redirect(shortURL)
	if err != nil {
		shortner.Logger.Printf("Error : %v", err)
		http.Error(w, fmt.Sprintf("Error : %v", err), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, orignalURL, http.StatusFound)

}

// Retrive the details around the URL stored
func HandleStats(w http.ResponseWriter, r *http.Request, shortner *shortner.UrlShortner) {
	stats := shortner.URLStats()
	shortner.Logger.Printf("Retrive stats : %v", stats)
	fmt.Fprintf(w, stats)
}
