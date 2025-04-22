package middleware

import (
	"context"
	"gemsvietnambe/pkg/auth"
	"gemsvietnambe/pkg/httputils"
	"net/http"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils := httputils.New()
		secretkey := r.Context().Value(ContextValuesK).(Values).Get(string(SecretkeyContextK))

		authID, err := auth.GetBearerToken(r.Header)
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, err)
			return
		}

		userID, err := auth.ValidateJWT(authID, secretkey)
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextK, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
