package httpjson

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ShareCampus/RecRem/backend/pkg/utils/logger"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Response[T any] struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
	Code    int    `json:"code,omitempty"`
	Result  *T     `json:"result,omitempty"`
}

func ReturnError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var resp = Response[struct{}]{
		Message: message,
		Success: false,
		Code:    code,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Errorf("Failed to encode response: %s", err.Error())
		logger.Debugf("Response: %+v", resp)
	}
}

func WriteJson[T any](w http.ResponseWriter, result T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var resp = Response[T]{
		Success: true,
		Result:  &result,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Errorf("Failed to encode response: %s", err.Error())
		logger.Debugf("Response: %+v", resp)
	}
	return err
}

func BindJson[T any](stream io.Reader) (*T, error) {
	var result = new(T)

	err := json.NewDecoder(stream).Decode(&result)
	if err != nil {
		logger.Debugf("Failed to decode request: %s", err.Error())
		return nil, err
	}

	err = validate.Struct(result)
	if err != nil {
		logger.Debugf("Failed to validate request: %s", err.Error())
		return nil, err
	}

	return result, nil
}
