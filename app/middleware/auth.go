package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"

	"goat/app/renderer"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			render.Status(r, http.StatusUnauthorized)
			renderer.PrettyJSON(w, r, "Missing Authorization Header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			render.Status(r, http.StatusUnauthorized)
			renderer.PrettyJSON(w, r, "Invalid token format")
			return
		}

		jwtKey := []byte(os.Getenv("JWT_SECRET"))
		if len(jwtKey) == 0 {
			jwtKey = []byte("default-secret-key")
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			render.Status(r, http.StatusUnauthorized)
			renderer.PrettyJSON(w, r, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Audience[0])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(UserRoleKey).(string)
			if !ok {
				render.Status(r, http.StatusForbidden)
				renderer.PrettyJSON(w, r, "Insufficient permissions")
				return
			}

			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			render.Status(r, http.StatusForbidden)
			renderer.PrettyJSON(w, r, "Insufficient permissions")
		})
	}
}
