package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON takes in a boolean and some data and a pointer to a writer and writes to a client
func WriteJSON(success bool, data interface{}, w http.ResponseWriter) {
	js, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !success {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
