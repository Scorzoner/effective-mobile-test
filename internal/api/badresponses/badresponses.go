package badresponses

import (
	"fmt"
	"net/http"

	"github.com/Scorzoner/effective-mobile-test/internal/api/jsonutil"
	"github.com/Scorzoner/effective-mobile-test/internal/logger"
)

func SendBadResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	errors := map[string]any{"errors": message}
	finalMessage := fmt.Errorf("encountered errors while processing %s %s request: %v",
		r.Method, r.URL.String(), message)
	switch status {
	case http.StatusInternalServerError:
		logger.Zap.Error(finalMessage)
	default:
		logger.Zap.Debug(finalMessage)
	}

	err := jsonutil.WriteJSON(w, status, errors, nil)
	if err != nil {
		logger.Zap.Error(fmt.Errorf("failed to write bad response: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, message any) {
	SendBadResponse(w, r, http.StatusBadRequest, fmt.Sprintf("%v", message))
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "resource not found"
	SendBadResponse(w, r, http.StatusNotFound, message)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	SendBadResponse(w, r, http.StatusUnprocessableEntity, fmt.Sprintf("%v", errors))
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	SendBadResponse(w, r, http.StatusMethodNotAllowed, message)
}

func InternalServerErrorResponse(w http.ResponseWriter, r *http.Request, message any) {
	SendBadResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("%v", message))
}
