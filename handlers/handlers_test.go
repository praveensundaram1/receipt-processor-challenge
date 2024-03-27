package handlers

import (
	"bytes"
	"net/http"
	"receipt-processor-challenge/models"
	"testing"

	json "github.com/json-iterator/go"
)


func GetSampleReceipt() models.Receipt {
	return models.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
		Total:  "9.00",
		Points: 109,
	}
}

func SimulateReceiptPostRequest(rBody models.Receipt) (*http.Request, error) {
	bodyBytes, err := json.Marshal(rBody)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPost, "/dummy-url", bytes.NewBuffer(bodyBytes))
}

func TestNewReceiptStore(t *testing.T) {
	receiptStore := NewReceiptStore()
	if receiptStore == nil {
		t.Fatal("NewReceiptStore() returned nil")
	}
}

func TestValidateReceipt(t *testing.T) {
	var tests = []struct {
		name      string
		modify    func(models.Receipt) models.Receipt
		expectErr bool
	}{
		{"Empty receipt", func(_ models.Receipt) models.Receipt { return models.Receipt{} }, true},
		{"Invalid Retailer", func(r models.Receipt) models.Receipt { r.Retailer = ""; return r }, true},
		{"Invalid Purchase Date", func(r models.Receipt) models.Receipt { r.PurchaseDate = "invalid-date"; return r }, true},
		{"Invalid Purchase Time", func(r models.Receipt) models.Receipt { r.PurchaseTime = "invalid-time"; return r }, true},
		{"Invalid Short Description", func(r models.Receipt) models.Receipt { r.Items[0].ShortDescription = "%@!,~"; return r }, true},
		{"Invalid Item Price", func(r models.Receipt) models.Receipt { r.Items[0].Price = "invalid-price"; return r }, true},
		{"Invalid Total", func(r models.Receipt) models.Receipt { r.Total = "invalid-total"; return r }, true},
		{"Valid Receipt", func(r models.Receipt) models.Receipt { return r }, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			receipt := GetSampleReceipt()
			request, err := SimulateReceiptPostRequest(tc.modify(receipt))
			if err != nil {
				t.Fatalf("SimulateReceiptPostRequest failed: %v", err)
			}

			_, err = checkReceiptValidity(request)
			if (err != nil) != tc.expectErr {
				t.Errorf("checkReceiptValidity() for %s: expected error %v, got %v", tc.name, tc.expectErr, err != nil)
			}
		})
	}
}

func TestComputeReceiptPoints(t *testing.T) {
	receipt := GetSampleReceipt()
	points := computeReceiptPoints(&receipt)
	expectedPoints := 109 // Update this value based on your business logic

	if points != expectedPoints {
		t.Errorf("computeReceiptPoints() got = %d, want %d", points, expectedPoints)
	}
}




