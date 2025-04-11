package httputils

import (
	"encoding/json"
	"fmt"
	localErrs "gemsvietnambe/internal/errors"
	"gemsvietnambe/pkg/logger"
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

type Responser struct{}

func (r *Responser) Response(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	codeClass := statusCode / 100

	if codeClass == 2 {
		response := SuccessfulResponse{Data: v, Timestamp: time.Now().UTC()}
		marshallingJSON(w, statusCode, response)
		return nil
	}

	if codeClass == 4 {
		errResponse := localErrs.Err{Message: fmt.Sprint(v)}
		if errInfo, ok := v.(localErrs.Err); ok {
			errResponse.Message = errInfo.Error()
			errResponse.Data = errInfo.Data
		}
		response := ClientErrorResponse{StatusCode: statusCode, Error: errResponse, TimeStamp: time.Now().UTC()}

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
		marshallingJSON(w, statusCode, response)
		return errResponse
	}

	marshallingJSON(w, statusCode, v)
	return nil
}

func marshallingJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error marshalling JSON: %s", err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
	return
}
