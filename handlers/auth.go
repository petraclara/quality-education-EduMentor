package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/middleware"
	"github.com/petraclara/quality-education-EduMentor/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register creates a new user account
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate fields
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))

	if req.Name == "" || req.Email == "" || req.Password == "" {
		jsonError(w, "name, email, and password are required", http.StatusBadRequest)
		return
	}

	if req.Role != "mentor" && req.Role != "learner" {
		jsonError(w, "role must be 'mentor' or 'learner'", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		jsonError(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := h.db.CreateUser(req.Name, req.Email, string(hash), req.Role)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			jsonError(w, "email already registered", http.StatusConflict)
			return
		}
		jsonError(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	// Auto-login after registration
	token, err := h.createSession(w, user.ID)
	if err != nil {
		jsonError(w, "account created but login failed", http.StatusInternalServerError)
		return
	}
	_ = token

	jsonResponse(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "registration successful",
		Data:    user,
	})
}

// Login authenticates a user
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := h.db.GetUserByEmail(req.Email)
	if err != nil {
		jsonError(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		jsonError(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	_, err = h.createSession(w, user.ID)
	if err != nil {
		jsonError(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "login successful",
		Data:    user,
	})
}

// Logout destroys the current session
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err == nil {
		h.db.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "logged out",
	})
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	prefsSet, _ := h.db.HasPreferencesSet(userID)

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"user":            user,
			"preferences_set": prefsSet,
		},
	})
}

func (h *AuthHandler) createSession(w http.ResponseWriter, userID int) (string, error) {
	token := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour * 7) // 7 days

	if err := h.db.CreateSession(token, userID, expiresAt); err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	return token, nil
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Helper functions for JSON responses
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	jsonResponse(w, status, models.APIResponse{
		Success: false,
		Message: message,
	})
}
