package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"receipt-processor-challenge/models"
	"regexp"
	"strconv"
	"strings"
	"time"
	//"reflect"

	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
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
		if !priceRegex.MatchString(item.Price){
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

func computeReceiptPoints(receipt *models.Receipt) int {
	points := computePointsFromRetailerName(receipt.Retailer)
	points += computeBonusForTotal(receipt.Total)
	points += computeItemPoints(receipt.Items)
	points += computeDateBonus(receipt.PurchaseDate)
	points += computeTimeBonus(receipt.PurchaseTime)
	return points
}

func computePointsFromRetailerName(retailer string) int {
	points := 0
	for _, char := range retailer {
		if isAlphanumeric(char) {
			points++
		}
	}
	return points
}

func isAlphanumeric(char rune) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9')
}

func computeBonusForTotal(total string) int {
	points := 0
	// Remove the decimal point and convert to cents. We know from the regex that the total is in the correct format.
	receiptTotalCents, err := strconv.ParseInt(strings.ReplaceAll(total, ".", ""), 10, 64)
	if err != nil {
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

func computeItemPoints(items []models.Item) int {
	points := 5 * (len(items) / 2)
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			fmt.Println("yes", item.ShortDescription)
			itemPriceCents, err := strconv.ParseFloat(strings.ReplaceAll(item.Price, ".", ""), 64)
			if err != nil {
				continue
			}
			points += int(math.Ceil(itemPriceCents * 0.2 / 100.0))
		}
	}
	return points
}

func computeDateBonus(purchaseDateStr string) int {
	purchaseDate, err := time.Parse(time.DateOnly, purchaseDateStr)
	if err != nil {
		return 0
	}
	if purchaseDate.Day()%2 == 1 {
		return pointsForOddDay
	}
	return 0
}

func computeTimeBonus(purchaseTimeStr string) int {
	//Go uses Mon Jan 2 15:04:05 MST 2006 as the reference time for parsing dates and times.
	purchaseTime, err := time.Parse("15:04", purchaseTimeStr)

	if err != nil {
		return 0
	}
	twoPM, _ := time.Parse("15:04", "14:00")
	fourPM, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(twoPM) && purchaseTime.Before(fourPM) {
		return pointsForTimeBetweenTwoAndFourPM
	}
	return 0
}

func (receiptStore *ReceiptStore) generateAndStoreReceipt(receipt *models.Receipt) (string, error) {
	receiptStore.lock.Lock()
	defer receiptStore.lock.Unlock()

	receiptHash, err := hashstructure.Hash(receipt, nil)
	if err != nil {
		return "", fmt.Errorf("error hashing receipt")
	}

	receiptID := strconv.FormatUint(receiptHash, 10)
	if _, exists := receiptStore.receipts[receiptID]; exists {
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
		return err
	}
	writeJSONResponse(w, http.StatusOK, data)
	return nil
}
