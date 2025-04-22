package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gemsvietnambe/internal/middleware"
	"gemsvietnambe/internal/model/database"
	"gemsvietnambe/pkg/auth"
	"gemsvietnambe/pkg/logger"
	"gemsvietnambe/pkg/validator"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Token     string    `json:"token,omitempty"`
}

const (
	tokenTime   = time.Minute * 15
	refreshTime = time.Minute * 30
)

func (h *Handlers) HandlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		validator.UserValidatorParams
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.Error("Couldn't decode parameters", err)
		h.responser.Response(w, http.StatusBadRequest, err)
		return
	}

	email := params.Email
	username := params.Username
	password := params.Password

	v := validator.New()

	if validator.ValidateUser(v, &validator.UserValidatorParams{
		Email:    email,
		Username: username,
		Password: password,
	}); !v.Valid() {
		h.responser.FailedValidates(w, http.StatusBadRequest, v)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		logger.Error("Couldn't hash password", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		logger.Error("Couldn't create user", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("Success created user: ")
	h.responser.Response(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func (h *Handlers) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		validator.UserValidatorParams
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		h.responser.Response(w, http.StatusBadRequest, err)
		logger.Error("Couldn't decode parameters", err)
		return
	}

	v := validator.New()
	if validator.ValidateUser(v, &validator.UserValidatorParams{
		Email:    params.Email,
		Password: params.Password,
	}); !v.Valid() {
		h.responser.FailedValidates(w, http.StatusBadRequest, v)
		return
	}

	user, err := h.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Couldn't found resource", err)
			h.responser.Response(w, http.StatusNotFound, "Invalid credentials")
			return
		}
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	err = auth.VerifyPassword(params.Password, user.PasswordHash)
	if err != nil {
		logger.Error("Invalid credentials", err)
		h.responser.Response(w, http.StatusUnauthorized, err)
		return
	}

	secretkey := r.Context().Value(middleware.ContextValuesK).(middleware.Values).Get(string(middleware.SecretkeyContextK))
	token, err := auth.GenerateJWTToken(secretkey, user.ID, tokenTime)
	if err != nil {
		logger.Error("Something wrong when generating jwt token", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	refreshkey := r.Context().Value(middleware.ContextValuesK).(middleware.Values).Get(string(middleware.RefreshkeyContextK))
	refreshToken, err := auth.GenerateRefreshToken(refreshkey, user.ID, refreshTime)
	if err != nil {
		logger.Error("Something wrong when generating refresh token", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	_, err = h.DB.SaveRefreshToken(r.Context(), database.SaveRefreshTokenParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		logger.Error("Can't save refresh token", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	h.responser.Response(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *Handlers) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		logger.Error("Couldn't decode parameters", err)
		h.responser.Response(w, http.StatusBadRequest, err)
		return
	}

	refreshkey := r.Context().Value(middleware.ContextValuesK).(middleware.Values).Get(string(middleware.RefreshkeyContextK))
	userID, err := auth.ValidateJWT(req.RefreshToken, refreshkey)
	if err != nil {
		logger.Error("Invalid refresh token", err)
		h.responser.Response(w, http.StatusBadRequest, err)
		return
	}

	refreshTokenRow, err := h.DB.RetrieveRefreshToken(r.Context(), userID)
	if err != nil || refreshTokenRow.RefreshToken != req.RefreshToken {
		logger.Error("Invalid refresh token", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	secretkey := r.Context().Value(middleware.ContextValuesK).(middleware.Values).Get(string(middleware.SecretkeyContextK))
	newToken, err := auth.GenerateJWTToken(secretkey, userID, tokenTime)
	if err != nil {
		logger.Error("Couldn't re create new jwt token", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	err = h.DB.RotationRefreshToken(r.Context(), refreshTokenRow.ID)
	if err != nil {
		logger.Error("Rotaion refresh token fail", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(refreshkey, userID, refreshTime)
	if err != nil {
		logger.Error("Rotaion refresh token fail", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}
	_, err = h.DB.SaveRefreshToken(r.Context(), database.SaveRefreshTokenParams{
		UserID:       userID,
		RefreshToken: newRefreshToken,
	})
	if err != nil {
		logger.Error("Rotaion refresh token fail", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	h.responser.Response(w, http.StatusOK, response{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	})
}

func (h *Handlers) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Message string `json:"message"`
	}

	userID, ok := r.Context().Value(middleware.UserIDContextK).(uuid.UUID)
	if !ok {
		logger.Error("User id not in context", errors.New("Unauthorized"))
		h.responser.Response(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	err := h.DB.DeleteRefreshTokenByUserID(r.Context(), userID)
	if err != nil {
		logger.Error("Fail to delete refresh token", err)
		h.responser.Response(w, http.StatusInternalServerError, err)
		return
	}

	h.responser.Response(w, http.StatusNoContent, response{
		Message: "Logout successfully!",
	})
}
