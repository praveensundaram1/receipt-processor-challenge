package handlers

import (
	"sync"

	"github.com/praveensundaram1/receipt-processor-challenge/models"
)

/*
ReceiptStore is a struct that represents the receipt store.
It contains a map of receipts and a mutex for synchronizing access to shared resources.
*/

type ReceiptStore struct {
	// receipts is a map that stores receipts by their unique identifier.
	// The key is a string representing the identifier, and the value is an instance of model.Receipt.
	receipts map[string]models.Receipt //maybe change to pointer
	// used for synchronizing access to shared resources.
	lock sync.RWMutex // RWMutex is a reader/writer mutex that allows multiple readers or a single writer.
}

func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{
		receipts: make(map[string]models.Receipt),
	}
}
