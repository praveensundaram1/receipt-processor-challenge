package handlers

import (
	"bytes"
	"net/http"
	"testing"

	json "github.com/json-iterator/go"
	"github.com/praveensundaram1/receipt-processor-challenge/models"
)

// GetSampleReceipt returns a sample receipt for testing.
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

// SimulateReceiptPostRequest creates a new POST request with the given receipt body.
func SimulateReceiptPostRequest(rBody models.Receipt) (*http.Request, error) {
	bodyBytes, err := json.Marshal(rBody)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPost, "/dummy-url", bytes.NewBuffer(bodyBytes))
}

// TestNewReceiptStore tests the NewReceiptStore function.
func TestNewReceiptStore(t *testing.T) {
	receiptStore := NewReceiptStore()
	if receiptStore == nil {
		t.Fatal("NewReceiptStore() returned nil")
	}
}

// TestProcessReceipt tests the ProcessReceipt function.
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

// TestComputeReceiptPoints tests the computeReceiptPoints function.
func TestComputeReceiptPoints(t *testing.T) {
	testCases := []struct {
		name           string
		modifyReceipt  func(models.Receipt) models.Receipt
		expectedPoints int
	}{
		{
			name: "Default receipt",
			modifyReceipt: func(receipt models.Receipt) models.Receipt {
				return receipt // No modification needed, use sample receipt as-is
			},
			expectedPoints: 109,
		},
		{
			name: "Modified receipt with different purchase date and items",
			modifyReceipt: func(receipt models.Receipt) models.Receipt {
				receipt.PurchaseDate = "2022-03-21"
				receipt.Items = []models.Item{
					{ShortDescription: "Gatorades", Price: "2.25"},
					{ShortDescription: "Gatorades", Price: "2.25"},
					{ShortDescription: "Gatorades", Price: "2.25"},
					{ShortDescription: "Gatorades", Price: "2.25"},
				}
				return receipt
			},
			expectedPoints: 119,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			receipt := GetSampleReceipt()
			modifiedReceipt := tc.modifyReceipt(receipt)
			points := computeReceiptPoints(&modifiedReceipt)
			if points != tc.expectedPoints {
				t.Errorf("%s: computeReceiptPoints() got = %d, expected %d", tc.name, points, tc.expectedPoints)
			}
		})
	}
}
