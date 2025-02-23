package middleware

import (
	"context"
	"net/http"
	"strings"
	"test-ai-api/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "No authorization header")
			return
		}

		// Bearer token format
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		claims, err := utils.ValidateJWT(bearerToken[1])
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
