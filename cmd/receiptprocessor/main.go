package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"receipt-processor-challenge/handlers"
	"receipt-processor-challenge/internal/routes"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {
	Port := os.Getenv("PORT")
	receiptStore := handlers.NewReceiptStore()
	router := routes.NewRouter(receiptStore) // Create a new router, and sets up the routes
	addr := "localhost" + Port
	fmt.Println("Listening on", addr)
	err := http.ListenAndServe(Port, router)
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}
}
