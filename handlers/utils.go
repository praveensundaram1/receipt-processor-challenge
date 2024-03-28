package handlers

import (
	"log"
	"net/http"
)

/*
*
Helper function to write JSON response with status code.
*
*/
func writeJSONResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(data) == 0 {
		return
	}
	if _, err := w.Write(data); err != nil {
		log.Printf("writeJSONResponse: error writing response: %v\n", err)
	}
}

/*
*
Helper function to handle errors.
*
*/
func handleErr(w http.ResponseWriter, err error, errorMessage string, statusCode int) {
	if err != nil {
		log.Printf("%s: %v", errorMessage, err)
	} else {
		log.Println(errorMessage)
	}
	writeJSONResponse(w, statusCode, nil)
}
