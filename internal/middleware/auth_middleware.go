package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/evbgsl/bank-service-evbgsl/pkg/response"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "authorization header required")
				return
			}

			tokenString, err := extractBearerToken(authHeader)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			claims := &jwt.RegisteredClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}

				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				response.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			userID, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil || userID <= 0 {
				response.Error(w, http.StatusUnauthorized, "invalid token subject")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	return userID, ok
}

func extractBearerToken(authHeader string) (string, error) {
	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("invalid bearer token")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}
