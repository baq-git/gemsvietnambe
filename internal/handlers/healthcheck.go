package handlers

import (
	"gemsvietnambe/pkg/httputils"
	"gemsvietnambe/pkg/logger"
	"net/http"
	"os"
)

func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := httputils.Envelope{
		"status":  "available",
		"message": "Application is running",
		"system_info": map[string]string{
			"environment": os.Getenv("ENV"),
		},
	}

	err := h.responser.Response(w, http.StatusOK, data)
	if err != nil {
		logger.Error("Health Check fail!", err)
	}
}
