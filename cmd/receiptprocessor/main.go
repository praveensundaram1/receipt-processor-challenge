package main

import (
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"

)

const (
	Port = ":8080"
)


func main() {
	router := httprouter.New() 
	service := NewReceiptService()  // receipts is a map that stores receipts by their unique identifier.
	// The key is a string representing the identifier, and the value is an instance of model.Receipt.
	service.SetupRoutes(router)
	log.Println("Listening on localhost:", Port)
	log.Fatal(http.ListenAndServe(Port, router))

}