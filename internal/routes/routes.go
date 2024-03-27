package routes

import (
	"receipt-processor-challenge/handlers"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(receiptStore *handlers.ReceiptStore) *httprouter.Router {
	router := httprouter.New()
	//receiptService := handlers.NewReceiptService()
	SetUpRoutes(router, receiptStore)
	return router
}

func SetUpRoutes(router *httprouter.Router, receiptStore *handlers.ReceiptStore) {
	router.POST("/receipts/process", receiptStore.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptStore.FetchPoints)
}
