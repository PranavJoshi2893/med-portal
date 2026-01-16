package responses

import "github.com/PranavJoshi2893/med-portal/internal/model"

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorEnvelope struct {
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    int                    `json:"code"`
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Fields  model.ValidationErrors `json:"fields,omitempty"`
}
