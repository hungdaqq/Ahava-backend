package response

import (
	"ahava/pkg/utils/models"
	"errors"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Error      interface{} `json:"error"`
}

func ClientErrorResponse(message string, data interface{}, err error) Response {
	var status_code int
	switch e := err.(type) {
	case error:
		switch {
		case errors.Is(e, models.ErrUnauthorized):
			status_code = http.StatusUnauthorized
		case errors.Is(e, models.ErrEntityNotFound):
			status_code = http.StatusNotFound
		case errors.Is(e, models.ErrBadRequest):
			status_code = http.StatusBadRequest
		case errors.Is(e, models.ErrConflict):
			status_code = http.StatusConflict
		case errors.Is(e, models.ErrForbidden):
			status_code = http.StatusForbidden
		default:
			status_code = http.StatusBadRequest
		}
	default:
		status_code = http.StatusBadRequest
	}

	return Response{
		StatusCode: status_code,
		Message:    message,
		Data:       data,
		Error:      err.Error(),
	}
}

func ClientResponse(status_code int, message string, data interface{}, err interface{}) Response {
	return Response{
		StatusCode: status_code,
		Message:    message,
		Data:       data,
		Error:      err,
	}
}

type WebhookResponse struct {
	Success bool `json:"success"`
}

func ClientWebhookResponse(status bool) WebhookResponse {
	return WebhookResponse{
		Success: status,
	}
}
