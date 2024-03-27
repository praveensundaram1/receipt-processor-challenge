package handlers

import (
	"net/http"
	"receipt-processor-challenge/models"
	"strings"

	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
)

func (receiptStore *ReceiptStore) ProcessReceipt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	receipt, err := checkReceiptValidity(r)
	if err != nil {
		handleErr(w, err, "ProcessReceipt validation error", http.StatusBadRequest)
		return
	}

	receiptID, err := receiptStore.generateAndStoreReceipt(receipt)
	if err != nil {
		handleErr(w, err, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendReceiptResponse(w, receiptID); err != nil {
		handleErr(w, err, "Error marshaling receipt response", http.StatusInternalServerError)
	}
}

func (receiptStore *ReceiptStore) FetchPoints(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	receiptID := strings.TrimSpace(params.ByName("id"))
	if receiptID == "" {
		handleErr(w, nil, "FetchPoints: No receipt ID provided", http.StatusBadRequest)
		return
	}
	receiptStore.lock.Lock()
	defer receiptStore.lock.Unlock()
	receipt, found := receiptStore.receipts[receiptID]
	if !found {
		handleErr(w, nil, "FetchPoints: Receipt not found", http.StatusNotFound)
		return
	}
	response := models.PointsResponse{Points: receipt.Points}
	data, err := json.Marshal(response)
	if err != nil {
		handleErr(w, err, "Error marshaling points response", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, data)
}
