package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/models"
)

type MentorHandler struct {
	db *database.DB
}

func NewMentorHandler(db *database.DB) *MentorHandler {
	return &MentorHandler{db: db}
}

// ListMentors returns all mentors
func (h *MentorHandler) ListMentors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mentors, err := h.db.GetMentors()
	if err != nil {
		jsonError(w, "failed to fetch mentors", http.StatusInternalServerError)
		return
	}

	if mentors == nil {
		mentors = []models.MentorCard{}
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    mentors,
	})
}

// GetMentor returns a single mentor by ID
func (h *MentorHandler) GetMentor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		jsonError(w, "invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		jsonError(w, "invalid mentor ID", http.StatusBadRequest)
		return
	}

	mentor, err := h.db.GetMentorByID(id)
	if err != nil {
		jsonError(w, "mentor not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    mentor,
	})
}
