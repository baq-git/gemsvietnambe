package handlers

import (
	cache "gemsvietnambe/internal/cache/app"
	"gemsvietnambe/internal/model/database"
	"gemsvietnambe/pkg/httputils"
	"time"
)

var handlerCache = cache.NewCache(15*time.Minute, 1*time.Minute)

type Handlers struct {
	DB        *database.Queries
	responser *httputils.Responser
}
