package models

type Receipt struct {
	Retailer     string `json:"retailer"`     //ex. "M&M Corner Market"
	PurchaseDate string `json:"purchaseDate"` //ex. "2022-01-01"
	PurchaseTime string `json:"purchaseTime"` //ex. "13:01"
	Items        []Item `json:"items"`
	Total        string `json:"total"`        //ex. "6.49"
	Points       int    `json:"pointsEarned"` //ex. 100
}

type Item struct {
	ShortDescription string `json:"shortDescription"` //ex. "Mountain Dew 12PK"
	Price            string `json:"price"`            //ex. "6.49"
}

type PointsResponse struct {
	Points int `json:"points"`
}

type ReceiptResponse struct {
	Id string `json:"id"`
}
