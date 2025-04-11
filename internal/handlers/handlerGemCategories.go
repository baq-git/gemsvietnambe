package handlers

import (
	"gemsvietnambe/pkg/logger"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type GemCategory struct {
	ID           uuid.UUID `json:"id"`
	CategoryName string    `json:"category_name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (h *Handlers) HandlerGemCategoriesRetrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	dbGemCategories, err := h.DB.GetAllGemCategories(r.Context())
	if err != nil {
		logger.Error("Couldn't retrive gem:", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	gemCategories := []GemCategory{}

	for _, dbGemCategory := range dbGemCategories {
		gemCategories = append(gemCategories, GemCategory{
			ID:           dbGemCategory.ID,
			CategoryName: dbGemCategory.CategoryName,
			Slug:         dbGemCategory.Slug,
			Description:  dbGemCategory.Description,
			CreatedAt:    dbGemCategory.CreatedAt,
			UpdatedAt:    dbGemCategory.CreatedAt,
		})
	}

	h.responser.Response(w, http.StatusOK, gemCategories)
	logger.Info("Retrieve gems successfully")
}
