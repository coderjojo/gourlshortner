package main

import (
	"github.com/coderjojo/gourlshortner/handler"
	"github.com/coderjojo/gourlshortner/shortner"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	logFile, err := os.Create("shortner.log")
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	defer logFile.Close()

	logger := log.New(logFile, "URLShortner: ", log.Ldate|log.Ltime|log.Lshortfile)

	shortner := shortner.NewUrlShorter(time.Hour*24, logger)
	router := mux.NewRouter()

	router.HandleFunc("/shortner", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleShortnerURL(w, r, shortner)
	}).Methods("POST")

	router.HandleFunc("/url/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleRedirect(w, r, shortner)
	}).Methods("GET")

	router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleStats(w, r, shortner)
	}).Methods("GET")

	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))

}
