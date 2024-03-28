package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"github.com/praveensundaram1/receipt-processor-challenge/models"
)

// Precompiling regular expressions for efficiency, readability, and reusability
var (
	retailerRegex    = regexp.MustCompile(`\S+`)
	dateRegex        = regexp.MustCompile(`^[1-2]\d{3}-[0-1]\d-[0-3]\d$`)
	timeRegex        = regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`)
	totalRegex       = regexp.MustCompile(`^\d+\.\d{2}$`)
	descriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
	priceRegex       = regexp.MustCompile(`^\d+\.\d{2}$`)
)

/*
*
This function checks the validity of the receipt data in the request body.
*
*/
func checkReceiptValidity(r *http.Request) (*models.Receipt, error) {
	var parsedReceipt models.Receipt
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "checkReceiptValidity: reading body failed")
	}

	if err := json.Unmarshal(requestBody, &parsedReceipt); err != nil {
		return nil, errors.Wrap(err, "checkReceiptValidity: unmarshaling failed")
	}

	if err := validateReceiptData(parsedReceipt); err != nil {
		return nil, err
	}

	return &parsedReceipt, nil
}

/*
*
Helper function to validate Retailer, Date, Time, Total, and Item Price & Description
*
*/
func validateReceiptData(receipt models.Receipt) error {

	validations := []struct {
		Regex   *regexp.Regexp
		Field   string
		Message string
	}{
		{retailerRegex, receipt.Retailer, "retailer validation failed"},
		{dateRegex, receipt.PurchaseDate, "date validation failed"},
		{timeRegex, receipt.PurchaseTime, "time validation failed"},
		{totalRegex, receipt.Total, "total validation failed"},
	}

	for _, validation := range validations {
		if !validation.Regex.MatchString(validation.Field) {
			return errors.New("validateReceiptData: " + validation.Message)
		}
	}

	for _, item := range receipt.Items {
		if !priceRegex.MatchString(item.Price) {
			return errors.New("validateReceiptData: price validation failed")
		}
		if !descriptionRegex.MatchString(item.ShortDescription) {
			return errors.New("validateReceiptData: description validation failed")
		}
	}

	return nil
}

const (
	pointsForRoundDollarTotal         = 50
	pointsForTotalInCentsMultipleOf25 = 25
	pointsForOddDay                   = 6
	pointsForTimeBetweenTwoAndFourPM  = 10
)

/*
*
This function computes the points for a receipt based on the retailer name, total, items, purchase date, and purchase time.
*
*/
func computeReceiptPoints(receipt *models.Receipt) int {
	points := computePointsFromRetailerName(receipt.Retailer)
	points += computeBonusForTotal(receipt.Total)
	points += computeItemPoints(receipt.Items)
	points += computeDateBonus(receipt.PurchaseDate)
	points += computeTimeBonus(receipt.PurchaseTime)
	return points
}

/*
*
Helper function to compute points based on the retailer name.
*
*/
func computePointsFromRetailerName(retailer string) int {
	points := 0
	for _, char := range retailer {
		if isAlphanumeric(char) {
			points++
		}
	}
	return points
}

/*
*
Helper function to check if a character is alphanumeric.
*
*/
func isAlphanumeric(char rune) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9')
}

/*
*
Helper function to compute bonus points based on the total.
*
*/
func computeBonusForTotal(total string) int {
	points := 0
	// Remove the decimal point and convert to cents. We know from the regex that the total is in the correct format.
	receiptTotalCents, err := strconv.ParseInt(strings.ReplaceAll(total, ".", ""), 10, 64)
	if err != nil {
		log.Println("Error parsing total")
		return 0
	}

	if receiptTotalCents%100 == 0 {
		points += pointsForRoundDollarTotal
	}
	if receiptTotalCents%25 == 0 {
		points += pointsForTotalInCentsMultipleOf25
	}
	return points
}

/*
*
Helper function to compute points based on the items in the receipt.
*
*/
func computeItemPoints(items []models.Item) int {
	points := 5 * (len(items) / 2)
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			itemPriceCents, err := strconv.ParseFloat(strings.ReplaceAll(item.Price, ".", ""), 64)
			if err != nil {
				log.Println("Error parsing item price")
				continue
			}
			points += int(math.Ceil(itemPriceCents * 0.2 / 100.0))
		}
	}
	return points
}

/*
*
Helper function to compute bonus points based on the purchase date
*
*/
func computeDateBonus(purchaseDateStr string) int {
	purchaseDate, err := time.Parse(time.DateOnly, purchaseDateStr)
	if err != nil {
		log.Println("Error parsing purchase date")
		return 0
	}
	if purchaseDate.Day()%2 == 1 {
		return pointsForOddDay
	}
	return 0
}

/*
*
Helper function to compute bonus points based on the purchase time
*
*/
func computeTimeBonus(purchaseTimeStr string) int {
	//Go uses Mon Jan 2 15:04:05 MST 2006 as the reference time for parsing dates and times.
	purchaseTime, err := time.Parse("15:04", purchaseTimeStr)

	if err != nil {
		log.Println("Error parsing purchase time")
		return 0
	}
	twoPM, _ := time.Parse("15:04", "14:00")
	fourPM, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(twoPM) && purchaseTime.Before(fourPM) {
		return pointsForTimeBetweenTwoAndFourPM
	}
	return 0
}

/*
*
This function generates a receipt ID and stores the receipt in the receipt store.
*
*/
func (receiptStore *ReceiptStore) generateAndStoreReceipt(receipt *models.Receipt) (string, error) {
	receiptStore.lock.Lock()
	defer receiptStore.lock.Unlock()

	receiptHash, err := hashstructure.Hash(receipt, nil) //Hashing receipt to generate a unique receipt ID
	if err != nil {
		log.Println("Error hashing receipt")
		return "", fmt.Errorf("error hashing receipt")
	}

	receiptID := strconv.FormatUint(receiptHash, 10) //Converting the hash to a string
	if _, exists := receiptStore.receipts[receiptID]; exists {
		log.Println("Duplicate receipt submission")
		return "", fmt.Errorf("ProcessReceipt: Duplicate receipt submission")
	}

	receipt.Points = computeReceiptPoints(receipt)
	receiptStore.receipts[receiptID] = *receipt

	return receiptID, nil
}

func sendReceiptResponse(w http.ResponseWriter, receiptID string) error {
	response := models.ReceiptResponse{Id: receiptID}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling receipt response")
		return err
	}
	writeJSONResponse(w, http.StatusOK, data)
	return nil
}
