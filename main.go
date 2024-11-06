package main

import (
	"backend/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/value", handlers.GetValueHandler)
	http.HandleFunc("/api/value/update", handlers.UpdateValueHandler)
	http.HandleFunc("/api/value/keep", handlers.KeepValueHandler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
