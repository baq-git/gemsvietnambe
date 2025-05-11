package handlers

import (
	"context"
	"gemsvietnambe/internal/model/database"
	"gemsvietnambe/pkg/logger"
	"gemsvietnambe/pkg/validator"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

type gem struct {
	GemID               uuid.UUID `json:"gemId"`
	GemName             string    `json:"gemName"`
	Instruction         string    `json:"instruction"`
	Description         string    `json:"description"`
	GemCategoryID       uuid.UUID `json:"gemCategoryId"`
	CategoryName        string    `json:"categoryName"`
	CategoryDescription string    `json:"categoryDescription"`
	CategorySlug        string    `json:"categorySlug"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

func (h *Handlers) HandlerGemsList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		validator.Query
	}

	type response struct {
		Gems     []gem    `json:"gems"`
		Metadata Metadata `json:"metadata"`
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	var gemCategories []database.GemCategory
	var err error

	cacheKey := "gem_category_all"
	cachedGemCats, ok := handlerCache.Get(cacheKey)
	if ok {
		logger.Info("gem cats hitted")
		gemCategories = cachedGemCats.([]database.GemCategory)
	}
	if !ok {
		logger.Info("gem cats missed")
		gemCategories, err = h.DB.GetAllGemCategories(r.Context())
		if err != nil {
			h.responser.Response(w, http.StatusInternalServerError, err)
			return
		}

		handlerCache.Set(cacheKey, gemCategories, 10*time.Minute)
	}

	var slugFilters []string
	for _, gc := range gemCategories {
		slugFilters = append(slugFilters, gc.Slug)
	}

	v := validator.New()
	qs := r.URL.Query()

	querySearch := readQueryString(qs, "search", "")
	filter := readQueryString(qs, "filter", "")
	limit := readQueryInt(qs, "page_size", 10, v)
	offset := readQueryInt(qs, "page", 1, v)

	input.Query.Search = querySearch
	input.Query.Filter = filter
	input.Query.PageSize = limit
	input.Query.Page = offset
	input.Query.Filters = slugFilters

	if validator.ValidatorQuery(v, input.Query); !v.Valid() {
		logger.Error("validate query fail", err)
		h.responser.FailedValidates(w, http.StatusBadRequest, v)
		return
	}

	dbGems, err := h.DB.ListingGemsComplex(ctx, database.ListingGemsComplexParams{
		PlaintoTsquery: input.Query.Search,
		Limit:          int32(input.Query.Limit()),
		Offset:         int32(input.Query.Offset()),
		Slug:           input.Query.Filter,
	})
	if err != nil {
		h.responser.Response(w, http.StatusBadRequest, err)
		return
	}

	var gems []gem
	for _, g := range dbGems {
		gems = append(gems, gem{
			GemID:               g.GemID,
			GemName:             g.GemName,
			CategoryName:        g.CategoryName,
			Description:         g.Description,
			Instruction:         g.Instruction,
			GemCategoryID:       g.GemCategoryID,
			CategoryDescription: g.CategoryDescription,
			CategorySlug:        g.CategorySlug,
			CreatedAt:           g.CreatedAt,
			UpdatedAt:           g.UpdatedAt,
		})
	}

	logger.Info("Gems retrieves")
	h.responser.Response(w, http.StatusOK, response{
		Gems:     gems,
		Metadata: calcMetadata(len(gems), input.Query.Page, input.Query.PageSize),
	})
}

func readQueryString(qs url.Values, key, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func readQueryInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be int value")
		return defaultValue
	}

	return i
}

func calcMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
