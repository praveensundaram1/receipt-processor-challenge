package models

// Receipt is a struct that represents a receipt.
type Receipt struct {
	Retailer     string `json:"retailer"`     //ex. "M&M Corner Market"
	PurchaseDate string `json:"purchaseDate"` //ex. "2022-01-01"
	PurchaseTime string `json:"purchaseTime"` //ex. "13:01"
	Items        []Item `json:"items"`
	Total        string `json:"total"`        //ex. "6.49"
	Points       int    `json:"pointsEarned"` //ex. 100
}

// Item is a struct that represents an item on a receipt. It contains a short description and a price.
type Item struct {
	ShortDescription string `json:"shortDescription"` //ex. "Mountain Dew 12PK"
	Price            string `json:"price"`            //ex. "6.49"
}

// PointsResponse is a struct that represents the response to a request for points. It contains the number of points.
type PointsResponse struct {
	Points int `json:"points"`
}

// ReceiptResponse is a struct that represents the response to a request to process a receipt. It contains the unique identifier of the receipt.
type ReceiptResponse struct {
	Id string `json:"id"`
}
