package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Karl1b/go-iban-bic/pkg/ibanbic"
	"github.com/Karl1b/go-iban-bic/pkg/settings"
	"github.com/Karl1b/go-iban-bic/pkg/utils"
)

func main() {
	// Register the health endpoint with CORS middleware
	http.HandleFunc("/health", enableCORS(healthHandler))
	http.HandleFunc("/iban", enableCORS(ibanHandler))

	log.Printf("runs on: %s", settings.Settings.Port)
	log.Fatal(http.ListenAndServe(":"+settings.Settings.Port, nil))
}

// CORS middleware to enable CORS for all origins
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// Health endpoint handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, struct{}{})

}

func ibanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type ResponseT struct {
		IBAN         string `json:"iban"`
		IsValid      bool   `json:"is_valid"`
		Bic          string `json:"bic"`
		Bezeichnung  string `json:"bezeichnung"`
		Ort          string `json:"ort"`
		Bankleitzahl string `json:"blz"`
	}

	// Parse the JSON request body
	var request struct {
		IBAN string `json:"iban"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isValid := ibanbic.ValidateIBAN(request.IBAN)

	var bic ibanbic.BicInfoDetailT
	if isValid {
		bic = ibanbic.GetBic(request.IBAN)
	}

	utils.RespondWithJSON(w, ResponseT{
		IBAN:         request.IBAN,
		IsValid:      isValid,
		Bic:          bic.BIC,
		Bezeichnung:  bic.Bezeichnung,
		Ort:          bic.Ort,
		Bankleitzahl: bic.Bankleitzahl,
	})
}
