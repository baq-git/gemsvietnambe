package httputils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func ReadIDParams(r *http.Request) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("id")
	if idStr == "" {
		return uuid.Nil, errors.New("missing id parameter")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID parameter: %w", err)
	}
	return id, nil
}
