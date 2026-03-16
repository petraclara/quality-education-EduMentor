package handlers

import (
	"net/http"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/middleware"
	"github.com/petraclara/quality-education-EduMentor/models"
)

type DashboardHandler struct {
	db *database.DB
}

func NewDashboardHandler(db *database.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetDashboard returns dashboard stats and recent matches
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)

	stats, err := h.db.GetDashboardStats(userID)
	if err != nil {
		jsonError(w, "failed to fetch stats", http.StatusInternalServerError)
		return
	}

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	profile, err := h.db.GetProfile(userID)
	if err != nil {
		jsonError(w, "profile not found", http.StatusNotFound)
		return
	}

	matches, err := h.db.GetMatchesByUser(userID)
	if err != nil {
		matches = []models.MatchWithUser{}
	}

	// Limit to 5 most recent
	recentMatches := matches
	if len(recentMatches) > 5 {
		recentMatches = recentMatches[:5]
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"user":           user,
			"profile":        profile,
			"stats":          stats,
			"recent_matches": recentMatches,
		},
	})
}
