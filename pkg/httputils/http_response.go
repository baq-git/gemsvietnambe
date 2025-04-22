package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	localErrs "gemsvietnambe/internal/errors"
	"gemsvietnambe/pkg/logger"
	"gemsvietnambe/pkg/validator"
	"net/http"
	"time"
)

// SuccessfulResponse standardizes responses with 200-299 status code
type SuccessfulResponse struct {
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// ClientErrorResponse standardizes responses with 400-499 status code
type ClientErrorResponse struct {
	StatusCode int           `json:"status_code"`
	Error      localErrs.Err `json:"error"`
	TimeStamp  time.Time     `json:"timestamp"`
}

// ServerErrorResponse standardizes responses with 500-599 status code
type ServerErrorResponse struct {
	StatusCode int           `json:"status_code"`
	Error      localErrs.Err `json:"error"`
	TimeStamp  time.Time     `json:"time_stamp"`
}

// FailedValidatesRepsone standardizes response with valid failed
type FailedValidatesRepsone struct {
	StatusCode int               `json:"status_code"`
	Errors     map[string]string `json:"error"`
	TimeStamp  time.Time         `json:"time_stamp"`
}

type Responser struct{}

func New() *Responser {
	return &Responser{}
}

func (r *Responser) FailedValidates(w http.ResponseWriter, statusCode int, v *validator.Validator) {
	w.Header().Set("Content-Type", "application/json")

	errReponse := localErrs.Err{Messages: v.Errors}.Errors()
	logger.Error("Validate failed", errors.New(localErrs.ErrValidateFailed))
	response := FailedValidatesRepsone{
		StatusCode: statusCode,
		Errors:     errReponse,
		TimeStamp:  time.Now().UTC(),
	}
	marshallingJSON(w, statusCode, response)
}

func (r *Responser) Response(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	codeClass := statusCode / 100

	if codeClass == 2 {
		response := SuccessfulResponse{Data: v, Timestamp: time.Now().UTC()}
		marshallingJSON(w, statusCode, response)
		logger.Response("Response: ", response)
		return nil
	}

	if codeClass == 4 {
		errResponse := localErrs.Err{Message: fmt.Sprint(v)}
		if errInfo, ok := v.(localErrs.Err); ok {
			errResponse.Message = errInfo.Error()
			errResponse.Data = errInfo.Data
		}
		response := ClientErrorResponse{StatusCode: statusCode, Error: errResponse, TimeStamp: time.Now().UTC()}
		logger.Error("Client Error Response: ", errResponse)
		marshallingJSON(w, statusCode, response)
		return errResponse
	}

	if codeClass == 5 {
		errResponse := localErrs.Err{Message: fmt.Sprint(v)}

		if errInfo, ok := v.(localErrs.Err); ok {
			errResponse.Message = errInfo.Error()
			errResponse.Data = errInfo.Data
		}
		response := ServerErrorResponse{StatusCode: statusCode, Error: errResponse, TimeStamp: time.Now().UTC()}
		logger.Error("Server Error Response: ", errResponse)
		marshallingJSON(w, statusCode, response)
		return errResponse
	}

	marshallingJSON(w, statusCode, v)
	return nil
}

func marshallingJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	dat, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		logger.Error("Error marshalling JSON: %s", err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
	return
}
