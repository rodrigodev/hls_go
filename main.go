package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	server := &http.Server{
		Handler: handlers(),
		Addr:    "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Starting server at localhost:8000")
	log.Fatal(server.ListenAndServe())
}

func handlers() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", mainHandler).Methods("GET")
	router.HandleFunc("/media/{mId:[0-9]+}/stream/", streamHandler).Methods("GET")
	router.HandleFunc("/media/{mId:[0-9]+}/stream/{segName:index[0-9]+.ts}", streamHandler).Methods("GET")

	router.HandleFunc("/upload", mainUpload).Methods("GET")
	router.HandleFunc("/video/upload/", uploadHandler).Methods("POST")

	return router
}

