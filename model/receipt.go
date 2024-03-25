package model

type Receipt struct {
	Retailer     string    `json:"retailer"`	 //ex. "M&M Corner Market"
	PurchaseDate string    `json:"purchaseDate"` //ex. "2022-01-01"
	PurchaseTime string    `json:"purchaseTime"` //ex. "13:01"
	Items        []Item    `json:"items"`		
	Total        string    `json:"total"`	     //ex. "6.49"
	PointsEarned int       `json:"pointsEarned"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"` //ex. "Mountain Dew 12PK"
	Price            string `json:"price"`			  //ex. "6.49"
}