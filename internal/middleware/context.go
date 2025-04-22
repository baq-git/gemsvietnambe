package middleware

import (
	"context"
	"net/http"
)

type Values struct {
	m map[string]string
}

func (v Values) Get(k string) string {
	return v.m[k]
}

func AddContext(values map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			valuesCtx := Values{map[string]string{}}
			for k, v := range values {
				valuesCtx.m[k] = v
			}
			ctx := context.WithValue(r.Context(), ContextValuesK, valuesCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
