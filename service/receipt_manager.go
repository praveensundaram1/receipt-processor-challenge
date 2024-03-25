package service

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	// Add the missing import statement for the "http" package
	//"github.com/pkg/errors"
)

func validateReceipt(r *http.Request) (*model.Receipt, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "validateReceipt: reading body failed")
	}

	var receipt model.Receipt
	if err = json.Unmarshal(body, &receipt); err != nil {
		return nil, errors.Wrap(err, "validateReceipt: unmarshaling failed")
	}

	if !regexp.MustCompile(`\S+`).MatchString(receipt.Retailer) {
		return nil, errors.New("validateReceipt: retailer validation failed")
	}

	if !regexp.MustCompile(`^[1-2]\d{3}-[0-1]\d-[0-3]\d$`).MatchString(receipt.PurchaseDate) {
		return nil, errors.New("validateReceipt: date validation failed")
	}

	if !regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`).MatchString(receipt.PurchaseTime) {
		return nil, errors.New("validateReceipt: time validation failed")
	}

	if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(receipt.Total) {
		return nil, errors.New("validateReceipt: total validation failed")
	}

	for _, item := range receipt.Items {
		if !regexp.MustCompile(`^[\w\s\-]+$`).MatchString(item.ShortDescription) || !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(item.Price) {
			return nil, errors.New("validateReceipt: item validation failed")
		}
	}

	return &receipt, nil
}

func calculateReceiptPoints(receipt *model.Receipt) int {
	points := 0
	for _, char := range receipt.Retailer {
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') {
			points++
		}
	}

	totalInCents, _ := strconv.ParseInt(strings.ReplaceAll(receipt.Total, ".", ""), 10, 64)
	if totalInCents%100 == 0 {
		points += 50
	}

	if totalInCents%25 == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			itemPriceInCents, _ := strconv.ParseFloat(strings.ReplaceAll(item.Price, ".", ""), 64)
			points += int(math.Ceil(itemPriceInCents * 0.2 / 100.0))
		}
	}

	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	twoPM, _ := time.Parse("15:04", "14:00")
	fourPM, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(twoPM) && purchaseTime.Before(fourPM) {
		points += 10
	}

	return points
}