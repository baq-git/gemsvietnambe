package handlers

import (
	"database/sql"
	"encoding/json"
	"gemsvietnambe/internal/model/database"
	"gemsvietnambe/pkg/httputils"
	"gemsvietnambe/pkg/logger"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Gem struct {
	ID            uuid.UUID `json:"id"`
	GemName       string    `json:"gemName"`
	Description   string    `json:"description"`
	Instruction   string    `json:"instruction"`
	GemCategoryID uuid.UUID `json:"gemCategoryId"`
	Coordinates   []float64 `json:"coordinates"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (h *Handlers) HandlerGemCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GemName       string    `json:"gemName"`
		Description   string    `json:"description"`
		Instruction   string    `json:"instruction"`
		GemCategoryID uuid.UUID `json:"gemCategoryId"`
		Coordinates   []float64 `json:"coordinates"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.Error("Couldn't decode parameters", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	gem, err := h.DB.CreateGem(r.Context(), database.CreateGemParams{
		GemName:       params.GemName,
		Description:   params.Description,
		Instruction:   params.Instruction,
		GemCategoryID: params.GemCategoryID,
		Coordinates:   params.Coordinates,
	})
	if err != nil {
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("Created gem")
	h.responser.Response(w, http.StatusCreated, gem)
}

func (h *Handlers) HandlerGemsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbGems, err := h.DB.GetAllGems(r.Context())
	if err != nil {
		logger.Error("Couldn't retrive gem:", err)
		h.responser.Response(w, http.StatusBadRequest, err)
	}

	gems := []Gem{}

	for _, dbGem := range dbGems {
		gems = append(gems, Gem{
			ID:            dbGem.ID,
			GemName:       dbGem.GemName,
			Description:   dbGem.Description,
			Instruction:   dbGem.Instruction,
			GemCategoryID: dbGem.GemCategoryID,
			Coordinates:   dbGem.Coordinates,
			CreatedAt:     dbGem.CreatedAt,
			UpdatedAt:     dbGem.CreatedAt,
		})
	}

	h.responser.Response(w, http.StatusOK, gems)
	logger.Info("Retrieved gems")
}

func (h *Handlers) HandlerGetGem(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ReadIDParams(r)
	if err != nil {
		logger.Error("Couldn't read id params", err)
		h.responser.Response(w, http.StatusBadRequest, err)
		return
	}

	gem, err := h.DB.GetGem(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Couldn't found resource", err)
			h.responser.Response(w, http.StatusNotFound, "Gem not found")
			return
		}
		logger.Error("Couldn't found resource", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("gem deleted")
	h.responser.Response(w, http.StatusOK, gem)
}

func (h *Handlers) HandlerGemDelete(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ReadIDParams(r)
	if err != nil {
		logger.Error("Couldn't read id params", err)
		h.responser.Response(w, http.StatusBadRequest, err)
		return
	}

	gem, err := h.DB.GetGem(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Couldn't found resource", err)
			h.responser.Response(w, http.StatusNotFound, "Gem not found")
			return
		}
		logger.Error("Couldn't found resource", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	err = h.DB.DeleteGem(r.Context(), gem.ID)
	if err != nil {
		logger.Error("Couldn't delete resource", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("gem deleted")
	h.responser.Response(w, http.StatusNoContent, map[string]string{
		"status": "success", "message": "Gem delete successfully",
	})
}

func (h *Handlers) HandlerGemUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GemName       string    `json:"gemName"`
		Description   string    `json:"description"`
		Instruction   string    `json:"instruction"`
		GemCategoryID uuid.UUID `json:"gemCategoryId"`
		Coordinates   []float64 `json:"coordinates"`
	}

	id, err := httputils.ReadIDParams(r)
	if err != nil {
		h.responser.Response(w, http.StatusBadRequest, err)
		return
	}

	gem, err := h.DB.GetGem(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.responser.Response(w, http.StatusNotFound, "Gem not found")
			return
		}
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		logger.Error("Couldn't decode parameters", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	if params.GemCategoryID != uuid.Nil {
		if _, err := h.DB.GetGemCategory(r.Context(), params.GemCategoryID); err != nil {
			if err == sql.ErrNoRows {
				h.responser.Response(w, http.StatusNotFound, "Gem category not found!")
				return
			}
			logger.Error("Couldn't retrieve gem category", err)
			h.responser.Response(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		params.GemCategoryID = gem.GemCategoryID
	}

	if params.Coordinates == nil {
		params.Coordinates = gem.Coordinates
	}

	newGem, err := h.DB.UpdateGem(r.Context(), database.UpdateGemParams{
		ID:            gem.ID,
		GemCategoryID: params.GemCategoryID,
		GemName:       fallback(params.GemName, gem.GemName),
		Description:   fallback(params.Description, gem.Description),
		Instruction:   fallback(params.Instruction, gem.Instruction),
		Coordinates:   params.Coordinates,
	})
	if err != nil {
		logger.Error("Couldn't update resource:", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("gem updated")
	h.responser.Response(w, http.StatusOK, newGem)
}

func fallback(provided, defaultValue string) string {
	if provided == "" {
		return defaultValue
	}
	return provided
}
