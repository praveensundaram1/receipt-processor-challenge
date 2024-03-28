package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/praveensundaram1/receipt-processor-challenge/handlers"
	"github.com/praveensundaram1/receipt-processor-challenge/internal/routes"
)

// Load from .env file and set up logging
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	f, err := os.OpenFile("receiptprocessor.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func main() {

	Port := os.Getenv("PORT")
	receiptStore := handlers.NewReceiptStore()
	router := routes.NewRouter(receiptStore) // Create a new router, and sets up the routes
	addr := "localhost" + Port
	fmt.Println("Listening on", addr)
	err := http.ListenAndServe(Port, router)
	if err != nil {
		log.Println("Error starting server:", err)
		panic(err)
	}
}
