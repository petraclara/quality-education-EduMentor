package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/middleware"
	"github.com/petraclara/quality-education-EduMentor/models"
)

type BookingHandler struct {
	db *database.DB
}

func NewBookingHandler(db *database.DB) *BookingHandler {
	return &BookingHandler{db: db}
}

// CreateBooking creates a new booking
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)

	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.MentorID == 0 || req.Date == "" || req.TimeSlot == "" {
		jsonError(w, "mentor_id, date, and time_slot are required", http.StatusBadRequest)
		return
	}

	// Don't allow booking yourself
	if req.MentorID == userID {
		jsonError(w, "you cannot book yourself", http.StatusBadRequest)
		return
	}

	// Verify mentor exists
	mentor, err := h.db.GetMentorByID(req.MentorID)
	if err != nil {
		jsonError(w, "mentor not found", http.StatusNotFound)
		return
	}

	// Check if the time slot is in the mentor's availability
	slotAvailable := false
	for _, avail := range mentor.Availability {
		if avail == req.TimeSlot {
			slotAvailable = true
			break
		}
	}
	if !slotAvailable {
		jsonError(w, "mentor is not available at this time", http.StatusBadRequest)
		return
	}

	// Check for conflicting bookings
	existing, _ := h.db.GetBookingsForMentorOnDate(req.MentorID, req.Date)
	for _, b := range existing {
		if b.TimeSlot == req.TimeSlot {
			jsonError(w, "this time slot is already booked", http.StatusConflict)
			return
		}
	}

	booking, err := h.db.CreateBooking(req.MentorID, userID, req.Date, req.TimeSlot, req.Note)
	if err != nil {
		jsonError(w, "failed to create booking", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "session booked successfully",
		Data:    booking,
	})
}

// GetBookings returns all bookings for current user
func (h *BookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	bookings, err := h.db.GetBookingsByUser(userID)
	if err != nil {
		jsonError(w, "failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    bookings,
	})
}
