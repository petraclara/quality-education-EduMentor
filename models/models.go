package models

import "time"

// User represents a registered user
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // "mentor", "mentee", or "both"
	CreatedAt    time.Time `json:"created_at"`
}

// Profile contains extended user information for matching
type Profile struct {
	ID           int      `json:"id"`
	UserID       int      `json:"user_id"`
	Bio          string   `json:"bio"`
	Skills       []string `json:"skills"`
	Interests    []string `json:"interests"`
	Availability []string `json:"availability"` // e.g. ["monday-morning", "tuesday-afternoon"]
	AvatarURL    string   `json:"avatar_url"`
}

// Match represents a mentor-mentee pairing
type Match struct {
	ID        int       `json:"id"`
	MentorID  int       `json:"mentor_id"`
	MenteeID  int       `json:"mentee_id"`
	Score     float64   `json:"score"`
	Status    string    `json:"status"` // "pending", "accepted", "rejected"
	CreatedAt time.Time `json:"created_at"`
}

// MatchWithUser adds user details to a match for API responses
type MatchWithUser struct {
	Match
	UserName  string   `json:"user_name"`
	UserEmail string   `json:"user_email"`
	UserRole  string   `json:"user_role"`
	Bio       string   `json:"bio"`
	Skills    []string `json:"skills"`
	Interests []string `json:"interests"`
}

// Session represents an active user session
type Session struct {
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// DashboardStats contains aggregated match statistics
type DashboardStats struct {
	TotalMatches    int `json:"total_matches"`
	PendingMatches  int `json:"pending_matches"`
	AcceptedMatches int `json:"accepted_matches"`
	RejectedMatches int `json:"rejected_matches"`
}

// API request/response types

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
	Availability []string `json:"availability"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
