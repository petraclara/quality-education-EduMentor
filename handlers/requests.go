package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/middleware"
	"github.com/petraclara/quality-education-EduMentor/models"
)

type RequestHandler struct {
	db *database.DB
}

func NewRequestHandler(db *database.DB) *RequestHandler {
	return &RequestHandler{db: db}
}

// HandleRequests routes GET and POST
func (h *RequestHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getRequests(w, r)
	case http.MethodPost:
		h.createRequest(w, r)
	default:
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RequestHandler) createRequest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	if user.Role != "learner" {
		jsonError(w, "only learners can send mentorship requests", http.StatusForbidden)
		return
	}

	var req models.MentorshipRequestCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.MentorID == 0 || req.HelpWith == "" {
		jsonError(w, "mentor_id and help_with are required", http.StatusBadRequest)
		return
	}

	// Verify mentor exists
	mentor, err := h.db.GetMentorByID(req.MentorID)
	if err != nil || mentor.Role != "mentor" {
		jsonError(w, "mentor not found", http.StatusNotFound)
		return
	}

	mentorshipReq, err := h.db.CreateMentorshipRequest(userID, req.MentorID, req.HelpWith, req.Goal, req.Message)
	if err != nil {
		jsonError(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "mentorship request sent",
		Data:    mentorshipReq,
	})
}

func (h *RequestHandler) getRequests(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	user, err := h.db.GetUserByID(userID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	if user.Role == "mentor" {
		requests, err := h.db.GetRequestsForMentor(userID)
		if err != nil {
			jsonError(w, "failed to fetch requests", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, models.APIResponse{Success: true, Data: requests})
	} else {
		requests, err := h.db.GetRequestsByLearner(userID)
		if err != nil {
			jsonError(w, "failed to fetch requests", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, models.APIResponse{Success: true, Data: requests})
	}
}

// HandleRequestAction handles accept/decline at /api/requests/{id}/{action}
func (h *RequestHandler) HandleRequestAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		jsonError(w, "invalid URL", http.StatusBadRequest)
		return
	}

	reqID, err := strconv.Atoi(parts[2])
	if err != nil {
		jsonError(w, "invalid request ID", http.StatusBadRequest)
		return
	}
	action := parts[3]

	// Verify request
	mentorshipReq, err := h.db.GetRequestByID(reqID)
	if err != nil {
		jsonError(w, "request not found", http.StatusNotFound)
		return
	}

	switch action {
	case "accept":
		if mentorshipReq.MentorID != userID {
			jsonError(w, "unauthorized", http.StatusForbidden)
			return
		}
		var body models.RequestActionBody
		json.NewDecoder(r.Body).Decode(&body)
		if err := h.db.AcceptRequest(reqID, body.MeetingType, body.MeetingLink, body.ProposedSlots); err != nil {
			jsonError(w, "failed to accept", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, models.APIResponse{Success: true, Message: "request accepted"})
	case "decline":
		if mentorshipReq.MentorID != userID {
			jsonError(w, "unauthorized", http.StatusForbidden)
			return
		}
		var body models.RequestActionBody
		json.NewDecoder(r.Body).Decode(&body)
		reason := body.DeclineReason
		if reason == "" {
			reason = "No reason given"
		}
		if err := h.db.DeclineRequest(reqID, reason); err != nil {
			jsonError(w, "failed to decline", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, models.APIResponse{Success: true, Message: "request declined"})
	case "confirm":
		if mentorshipReq.LearnerID != userID {
			jsonError(w, "unauthorized", http.StatusForbidden)
			return
		}
		var body models.ConfirmSlotRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Date == "" || body.Time == "" {
			jsonError(w, "date and time are required", http.StatusBadRequest)
			return
		}
		if err := h.db.ConfirmRequestSlot(reqID, body.Date, body.Time); err != nil {
			jsonError(w, "failed to confirm slot", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, models.APIResponse{Success: true, Message: "slot confirmed"})
	default:
		jsonError(w, "action must be 'accept', 'decline' or 'confirm'", http.StatusBadRequest)
	}
}
