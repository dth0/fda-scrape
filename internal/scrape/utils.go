package scrape

import (
	"encoding/json"
	"net/http"
)

func response(w http.ResponseWriter, payload any, statusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode == http.StatusOK {
		json.NewEncoder(w).Encode(payload)
	}

}
