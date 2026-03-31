package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/petraclara/quality-education-EduMentor/database"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Auth middleware checks for valid session cookie
func Auth(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, `{"success":false,"message":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			session, err := db.GetSession(cookie.Value)
			if err != nil || session.ExpiresAt.Before(time.Now()) {
				http.Error(w, `{"success":false,"message":"session expired"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CORS middleware for React dev server & Production
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts user ID from context
func GetUserID(r *http.Request) int {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		return 0
	}
	return userID
}
