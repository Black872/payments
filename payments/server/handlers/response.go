package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
	log.Println(err)
}
