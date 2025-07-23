package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, payload any) {
	var dat []byte
	statuscode := 200
	w.Header().Add("Content-Type", "application/json")

	dat, _ = json.Marshal(payload)

	w.WriteHeader(statuscode)
	w.Write(dat)

}
