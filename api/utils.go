package main

import (
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		_, err := w.Write(data)
		if err != nil {
			log.Println("Error writing response:", err)
		}
	}
}