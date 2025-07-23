package ibanbic

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type BicInfoDetailT struct {
	BIC          string
	Bezeichnung  string
	Ort          string
	Bankleitzahl string
}

type BicInfoT map[string]BicInfoDetailT // key is Bankleitzahl (bank code)

var BicInfo BicInfoT

func init() {

	readBundesBankCSV()

}

// GetBic returns BIC depending on IBAN using the map BicInfo
func GetBic(iban string) BicInfoDetailT {
	// Remove spaces and convert to uppercase
	cleanIban := strings.ToUpper(strings.ReplaceAll(iban, " ", ""))

	// Validate IBAN format (basic check)
	ibanRegex := regexp.MustCompile(`^[A-Z]{2}\d{2}[A-Z0-9]+$`)
	if !ibanRegex.MatchString(cleanIban) {
		return BicInfoDetailT{}
	}

	// For German IBAN, extract bank code (Bankleitzahl)
	// German IBAN format: DE + 2 check digits + 8 digit bank code + 10 digit account number
	if strings.HasPrefix(cleanIban, "DE") && len(cleanIban) == 22 {
		bankCode := cleanIban[4:12] // Extract 8-digit bank code

		if bicInfo, exists := BicInfo[bankCode]; exists {
			return bicInfo
		}
	}

	return BicInfoDetailT{} // BIC not found or invalid IBAN
}

// GetBicInfo returns full BIC information for an IBAN
func GetBicInfo(iban string) (BicInfoDetailT, bool) {
	cleanIban := strings.ToUpper(strings.ReplaceAll(iban, " ", ""))

	if strings.HasPrefix(cleanIban, "DE") && len(cleanIban) == 22 {
		bankCode := cleanIban[4:12]

		if bicInfo, exists := BicInfo[bankCode]; exists {
			return bicInfo, true
		}
	}

	return BicInfoDetailT{}, false
}

func ValidateIBAN(iban string) bool {
	// Remove spaces and convert to uppercase
	iban = strings.ReplaceAll(strings.ToUpper(iban), " ", "")

	// Check if length is between 15 and 34 characters
	if len(iban) < 15 || len(iban) > 34 {
		return false
	}

	// Check if first two characters are letters (country code)
	if len(iban) < 2 || !isLetter(iban[0]) || !isLetter(iban[1]) {
		return false
	}

	// Check if characters 3-4 are digits (check digits)
	if len(iban) < 4 || !isDigit(iban[2]) || !isDigit(iban[3]) {
		return false
	}

	// Rearrange: move first 4 characters to the end
	rearranged := iban[4:] + iban[:4]

	// Convert letters to numbers (A=10, B=11, ..., Z=35)
	var numericString string
	for _, char := range rearranged {
		if isLetter(byte(char)) {
			numericString += fmt.Sprintf("%d", int(char-'A'+10))
		} else if isDigit(byte(char)) {
			numericString += string(char)
		} else {
			return false // Invalid character
		}
	}

	// Perform mod 97 calculation
	remainder := 0
	for _, digit := range numericString {
		remainder = (remainder*10 + int(digit-'0')) % 97
	}

	return remainder == 1
}

func isLetter(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func readBundesBankCSV() {

	BicInfo = make(BicInfoT)

	// Open CSV file
	file, err := os.Open("blz-aktuell-csv-data.csv")
	if err != nil {
		// Handle error gracefully - in production you might want to log this
		return
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';' // ;-)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	// Skip header row and parse data
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		if len(record) >= 8 {
			bankleitzahl := strings.Trim(record[0], `"`)
			bezeichnung := strings.Trim(record[2], `"`)
			ort := strings.Trim(record[4], `"`)
			bic := strings.Trim(record[7], `"`)

			BicInfo[bankleitzahl] = BicInfoDetailT{
				BIC:          bic,
				Bezeichnung:  bezeichnung,
				Ort:          ort,
				Bankleitzahl: bankleitzahl,
			}
		}
	}
}
