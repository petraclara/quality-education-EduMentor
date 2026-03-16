package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/middleware"
	"github.com/petraclara/quality-education-EduMentor/models"
)

type ProfileHandler struct {
	db *database.DB
}

func NewProfileHandler(db *database.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

// HandleProfile routes GET/PUT requests
func (h *ProfileHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getProfile(w, r)
	case http.MethodPut:
		h.updateProfile(w, r)
	default:
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProfileHandler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	profile, err := h.db.GetProfile(userID)
	if err != nil {
		jsonError(w, "profile not found", http.StatusNotFound)
		return
	}

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"user":    user,
			"profile": profile,
		},
	})
}

func (h *ProfileHandler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateProfile(userID, req); err != nil {
		jsonError(w, "failed to update profile", http.StatusInternalServerError)
		return
	}

	// Return updated profile
	profile, _ := h.db.GetProfile(userID)
	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "profile updated",
		Data:    profile,
	})
}
