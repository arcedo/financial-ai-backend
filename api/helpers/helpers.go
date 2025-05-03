package helpers

import (
	"encoding/json"
	"net/http"

	db "github.com/arcedo/financial-ai-backend/database"
	"github.com/arcedo/financial-ai-backend/utils"
)

type ApiFunc func(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error

type APIResponse struct {
	Data    any    `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data any, apiErr *utils.APIError, message string) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
	}

	response := APIResponse{
		Data:    data,
		Error:   "",
		Message: message,
	}

	if apiErr != nil {
		response.Error = apiErr.Code
		response.Message = apiErr.Message
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write response"})
	}
}

func MakeHTTPHandleFunc(fn ApiFunc, s db.MongoStorage, allowedMethods []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is allowed
		if err := utils.ValidateMethods(r, allowedMethods); err != nil {
			errValue := utils.ErrorMap[err]
			WriteJSON(w, http.StatusMethodNotAllowed, nil, &errValue, "")
			return
		}

		// Call the handler function and handle errors
		if err := fn(w, r, s); err != nil {
			// Map the error to an appropriate HTTP status code and API error
			apiErr, statusCode := utils.MapErrorToAPIError(err)
			WriteJSON(w, statusCode, nil, apiErr, "")
			return
		}

	}
}
