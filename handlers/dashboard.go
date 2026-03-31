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

// GetDashboard returns role-specific dashboard data
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
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

	profile, _ := h.db.GetProfile(userID)

	if user.Role == "mentor" {
		// Mentor dashboard: incoming requests
		requests, err := h.db.GetRequestsForMentor(userID)
		if err != nil {
			jsonError(w, "failed to fetch requests", http.StatusInternalServerError)
			return
		}

		pending := 0
		accepted := 0
		for _, req := range requests {
			if req.Status == "pending" {
				pending++
			} else if req.Status == "accepted" || req.Status == "scheduled" {
				accepted++
			}
		}

		jsonResponse(w, http.StatusOK, models.APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"role":             "mentor",
				"user":             user,
				"profile":          profile,
				"requests":         requests,
				"pending_count":    pending,
				"accepted_count":   accepted,
				"total_requests":   len(requests),
			},
		})
	} else {
		// Learner dashboard: matched mentors + request statuses
		mentors, err := h.db.GetMatchedMentors(userID)
		if err != nil {
			jsonError(w, "failed to fetch mentors", http.StatusInternalServerError)
			return
		}

		requests, err := h.db.GetRequestsByLearner(userID)
		if err != nil {
			requests = []models.MentorshipRequestWithUser{}
		}

		jsonResponse(w, http.StatusOK, models.APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"role":     "learner",
				"user":     user,
				"profile":  profile,
				"mentors":  mentors,
				"requests": requests,
			},
		})
	}
}
