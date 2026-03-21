package models

import "time"

// User represents a registered user
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // "learner" or "mentor"
	CreatedAt    time.Time `json:"created_at"`
}

// Profile contains extended user information for matching
type Profile struct {
	ID           int      `json:"id"`
	UserID       int      `json:"user_id"`
	Bio          string   `json:"bio"`
	Skills       []string `json:"skills"`
	Interests    []string `json:"interests"`
	Level        string   `json:"level"` // "beginner", "intermediate", "advanced"
	Goal         string   `json:"goal"`
	Availability []string `json:"availability"`
	AvatarURL    string   `json:"avatar_url"`
	MaxMentees   int      `json:"max_mentees"`
	Rating       float64  `json:"rating"`
	RatingCount  int      `json:"rating_count"`
}

// MentorCard is a flattened view for the mentor listing
type MentorCard struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	Bio          string   `json:"bio"`
	Skills       []string `json:"skills"`
	Interests    []string `json:"interests"`
	Level        string   `json:"level"`
	Availability []string `json:"availability"`
	AvatarURL    string   `json:"avatar_url"`
	Rating       float64  `json:"rating"`
	RatingCount  int      `json:"rating_count"`
	MatchScore   float64  `json:"match_score,omitempty"`
}

// MentorshipRequest represents a learner's request to a mentor
type MentorshipRequest struct {
	ID            int       `json:"id"`
	LearnerID     int       `json:"learner_id"`
	MentorID      int       `json:"mentor_id"`
	HelpWith      string    `json:"help_with"`
	Goal          string    `json:"goal"`
	Message       string    `json:"message"`
	Status        string    `json:"status"` // "pending", "accepted", "declined"
	DeclineReason string    `json:"decline_reason,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// MentorshipRequestWithUser includes user names
type MentorshipRequestWithUser struct {
	MentorshipRequest
	LearnerName  string `json:"learner_name"`
	LearnerEmail string `json:"learner_email"`
	LearnerLevel string `json:"learner_level"`
	MentorName   string `json:"mentor_name"`
}

// Booking represents a scheduled mentoring session
type Booking struct {
	ID        int       `json:"id"`
	MentorID  int       `json:"mentor_id"`
	MenteeID  int       `json:"mentee_id"`
	Date      string    `json:"date"`
	TimeSlot  string    `json:"time_slot"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// BookingWithUser includes mentor/mentee name in booking
type BookingWithUser struct {
	Booking
	MentorName string `json:"mentor_name"`
	MenteeName string `json:"mentee_name"`
}

// Match represents a mentor-mentee pairing (kept for algorithmic matching)
type Match struct {
	ID        int       `json:"id"`
	MentorID  int       `json:"mentor_id"`
	MenteeID  int       `json:"mentee_id"`
	Score     float64   `json:"score"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type MatchWithUser struct {
	Match
	UserName  string   `json:"user_name"`
	UserEmail string   `json:"user_email"`
	UserRole  string   `json:"user_role"`
	Bio       string   `json:"bio"`
	Skills    []string `json:"skills"`
	Interests []string `json:"interests"`
}

type Session struct {
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DashboardStats struct {
	TotalMatches    int `json:"total_matches"`
	PendingMatches  int `json:"pending_matches"`
	AcceptedMatches int `json:"accepted_matches"`
	RejectedMatches int `json:"rejected_matches"`
}

// ==================== API Request/Response Types ====================

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProfileUpdateRequest struct {
	Bio          string   `json:"bio"`
	Skills       []string `json:"skills"`
	Interests    []string `json:"interests"`
	Level        string   `json:"level"`
	Goal         string   `json:"goal"`
	Availability []string `json:"availability"`
	MaxMentees   int      `json:"max_mentees"`
}

type MentorshipRequestCreate struct {
	MentorID int    `json:"mentor_id"`
	HelpWith string `json:"help_with"`
	Goal     string `json:"goal"`
	Message  string `json:"message"`
}

type RequestActionBody struct {
	DeclineReason string `json:"decline_reason,omitempty"`
}

type BookingRequest struct {
	MentorID int    `json:"mentor_id"`
	Date     string `json:"date"`
	TimeSlot string `json:"time_slot"`
	Note     string `json:"note"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
