package main

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"github.com/praveensundaram1/receipt-processor-challenge/model"
	
	//"example.com/project/model"
)

const (
	dateFormat = "2006-01-02" // ISO 8601 format
	timeFormat = "13:45"
	// afternoonStart = "14:00"
	// afternoonEnd   = "16:00"
)

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	json "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"github.com/praveensundaram1/receipt-processor-challenge/model"
	"github.com/praveensundaram1/receipt-processor-challenge/service"
)

const (
	dateFormat = "2006-01-02" // ISO 8601 format
	timeFormat = "13:45"
	// afternoonStart = "14:00"
	// afternoonEnd   = "16:00"
)

func (service *service.ReceiptService) SetupRoutes(router *httprouter.Router) {
	router.POST("/receipts/process", service.ProcessReceipt)
	router.GET("/receipts/:id/points", service.FetchPoints)
}

func (service *ReceiptService) FetchPoints(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	receiptID := params.ByName("id")
	if receiptID == "" {
		log.Println("FetchPoints: No receipt ID provided")
		respondWithJSON(w, http.StatusBadRequest, nil)
		return
	}

	service.lock.Lock()
	defer service.lock.Unlock()

	receipt, found := service.receipts[receiptID]
	if !found {
		log.Println("FetchPoints: Receipt not found")
		respondWithJSON(w, http.StatusNotFound, nil)
		return
	}

	points := receipt.Points
	response := model.PointsResponse{Points: points}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling points response:", err)
		respondWithJSON(w, http.StatusInternalServerError, nil)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

func (service *ReceiptService) ProcessReceipt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	receipt, err := validateReceipt(r)
	if err != nil {
		log.Println("ProcessReceipt validation error:", err)
		respondWithJSON(w, http.StatusBadRequest, nil)
		return
	}

	service.lock.Lock()
	defer service.lock.Unlock()

	receiptHash, err := hashstructure.Hash(receipt, nil)
	if err != nil {
		log.Println("Error hashing receipt:", err)
		respondWithJSON(w, http.StatusInternalServerError, nil)
		return
	}

	receiptID := strconv.FormatUint(receiptHash, 10)
	if _, exists := service.receipts[receiptID]; exists {
		log.Println("ProcessReceipt: Duplicate receipt submission")
		respondWithJSON(w, http.StatusConflict, nil)
		return
	}

	receipt.Points = calculateReceiptPoints(receipt)
	service.receipts[receiptID] = *receipt

	response := model.ReceiptResponse{ID: receiptID}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling receipt response:", err)
		respondWithJSON(w, http.StatusInternalServerError, nil)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}


