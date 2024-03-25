package main

import "sync"
/*
ReceiptService is a struct that represents the receipt service.
It contains a map of receipts and a mutex for synchronizing access to shared resources.
*/
type ReceiptService struct {
	// receipts is a map that stores receipts by their unique identifier.
	// The key is a string representing the identifier, and the value is an instance of model.Receipt.
	receipts map[string]model.Receipt  //maybe change to pointer
	// used for synchronizing access to shared resources.
	lock     sync.RWMutex // RWMutex is a reader/writer mutex that allows multiple readers or a single writer.
}

func NewReceiptService() *ReceiptService {
	return &ReceiptService{
		receipts: make(map[string]model.Receipt),
	}
}