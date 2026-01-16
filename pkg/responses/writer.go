package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteSuccess(
	w http.ResponseWriter,
	status int,
	message string,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	}); err != nil {
		log.Printf("error encoding success response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, errResp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.Code)

	if err := json.NewEncoder(w).Encode(ErrorEnvelope{
		Error: errResp,
	}); err != nil {
		log.Printf("error encoding error response: %v", err)
	}
}
