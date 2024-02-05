package utils

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}, err interface{}) {
	responseStruct := struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  interface{} `json:"error,omitempty"`
	}{
		Status: status,
		Data:   data,
		Error:  err,
	}

	response, err := json.Marshal(responseStruct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error processing response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
