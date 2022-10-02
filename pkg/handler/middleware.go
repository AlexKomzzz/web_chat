package handler

import (
	"context"
	"net/http"
	"strings"
)

type myCtx string

// парсинг хедера, определение JWT, определение id
func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization") // выделяем из заголовка поле "Authorization"
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if userId < 1 {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}

		username, err := h.service.GetUsername(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var keyName myCtx = "username"

		h.webClient.ctx = context.WithValue(h.webClient.ctx, keyName, username)
		//c.Set("userId", userId)

		next.ServeHTTP(w, r)
	})
}
