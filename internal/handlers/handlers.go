package handlers

import (
	"gemsvietnambe/internal/model/database"
	"gemsvietnambe/pkg/httputils"
)

type Handlers struct {
	DB        *database.Queries
	responser *httputils.Responser
}
