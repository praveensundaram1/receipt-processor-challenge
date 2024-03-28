package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/praveensundaram1/receipt-processor-challenge/handlers"
)

/*
*
NewRouter creates a new router and sets up the routes.
*
*/
func NewRouter(receiptStore *handlers.ReceiptStore) *httprouter.Router {
	router := httprouter.New()
	SetUpRoutes(router, receiptStore)
	return router
}

/*
*
SetUpRoutes sets up the routes for the router.
*
*/
func SetUpRoutes(router *httprouter.Router, receiptStore *handlers.ReceiptStore) {
	router.POST("/receipts/process", receiptStore.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptStore.FetchPoints)
}
