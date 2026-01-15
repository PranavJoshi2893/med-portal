package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	}); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
