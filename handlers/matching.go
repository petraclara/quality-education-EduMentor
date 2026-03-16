package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/middleware"
	"github.com/petraclara/quality-education-EduMentor/models"
)

type MatchHandler struct {
	db *database.DB
}

func NewMatchHandler(db *database.DB) *MatchHandler {
	return &MatchHandler{db: db}
}

// FindMatches computes and returns ranked matches for current user
func (h *MatchHandler) FindMatches(w http.ResponseWriter, r *http.Request) {
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

	myProfile, err := h.db.GetProfile(userID)
	if err != nil {
		jsonError(w, "profile not found", http.StatusNotFound)
		return
	}

	allProfiles, err := h.db.GetAllProfiles()
	if err != nil {
		jsonError(w, "failed to fetch profiles", http.StatusInternalServerError)
		return
	}

	var matches []models.MatchWithUser
	for _, other := range allProfiles {
		if other.User.ID == userID {
			continue
		}

		// Check role compatibility
		if !isCompatible(user.Role, other.User.Role) {
			continue
		}

		// Compute match score
		score := computeScore(myProfile, &other.Profile)
		if score <= 0 {
			continue
		}

		// Determine mentor/mentee IDs
		mentorID, menteeID := determinePairing(user, &other.User)

		// Create or update match in database
		match, err := h.db.CreateMatch(mentorID, menteeID, score)
		if err != nil {
			continue
		}

		matches = append(matches, models.MatchWithUser{
			Match:     *match,
			UserName:  other.User.Name,
			UserEmail: other.User.Email,
			UserRole:  other.User.Role,
			Bio:       other.Profile.Bio,
			Skills:    other.Profile.Skills,
			Interests: other.Profile.Interests,
		})
	}

	// Sort by score descending (already computed with highest first)
	sortMatchesByScore(matches)

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    matches,
	})
}

// GetMyMatches returns existing matches for current user
func (h *MatchHandler) GetMyMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	matches, err := h.db.GetMatchesByUser(userID)
	if err != nil {
		jsonError(w, "failed to fetch matches", http.StatusInternalServerError)
		return
	}

	if matches == nil {
		matches = []models.MatchWithUser{}
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    matches,
	})
}

// UpdateMatch handles accept/reject
func (h *MatchHandler) UpdateMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract match ID and action from URL: /api/matches/{id}/{action}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		jsonError(w, "invalid URL", http.StatusBadRequest)
		return
	}

	matchID, err := strconv.Atoi(parts[2])
	if err != nil {
		jsonError(w, "invalid match ID", http.StatusBadRequest)
		return
	}

	action := parts[3] // "accept" or "reject"
	if action != "accept" && action != "reject" {
		jsonError(w, "action must be 'accept' or 'reject'", http.StatusBadRequest)
		return
	}

	// Verify the match belongs to the current user
	userID := middleware.GetUserID(r)
	match, err := h.db.GetMatchByID(matchID)
	if err != nil {
		jsonError(w, "match not found", http.StatusNotFound)
		return
	}

	if match.MentorID != userID && match.MenteeID != userID {
		jsonError(w, "unauthorized", http.StatusForbidden)
		return
	}

	status := "accepted"
	if action == "reject" {
		status = "rejected"
	}

	if err := h.db.UpdateMatchStatus(matchID, status); err != nil {
		jsonError(w, "failed to update match", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "match " + status,
	})
}

// ==================== Matching Algorithm ====================

// computeScore calculates compatibility between two profiles
// Weights: skills 40%, interests 30%, availability 30%
func computeScore(a, b *models.Profile) float64 {
	skillScore := jaccardSimilarity(a.Skills, b.Skills)
	interestScore := jaccardSimilarity(a.Interests, b.Interests)
	availabilityScore := jaccardSimilarity(a.Availability, b.Availability)

	return skillScore*0.4 + interestScore*0.3 + availabilityScore*0.3
}

// jaccardSimilarity computes |A ∩ B| / |A ∪ B|
func jaccardSimilarity(a, b []string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}

	setA := make(map[string]bool)
	for _, v := range a {
		setA[strings.ToLower(strings.TrimSpace(v))] = true
	}

	setB := make(map[string]bool)
	for _, v := range b {
		setB[strings.ToLower(strings.TrimSpace(v))] = true
	}

	intersection := 0
	for k := range setA {
		if setB[k] {
			intersection++
		}
	}

	union := len(setA)
	for k := range setB {
		if !setA[k] {
			union++
		}
	}

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// isCompatible checks if two users can be matched
func isCompatible(roleA, roleB string) bool {
	// "both" matches with anyone
	if roleA == "both" || roleB == "both" {
		return true
	}
	// mentor matches with mentee and vice versa
	return (roleA == "mentor" && roleB == "mentee") ||
		(roleA == "mentee" && roleB == "mentor")
}

// determinePairing decides who is mentor and who is mentee
func determinePairing(userA, userB *models.User) (mentorID, menteeID int) {
	if userA.Role == "mentor" || (userA.Role == "both" && userB.Role == "mentee") {
		return userA.ID, userB.ID
	}
	return userB.ID, userA.ID
}

// sortMatchesByScore sorts matches by score descending
func sortMatchesByScore(matches []models.MatchWithUser) {
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].Score > matches[i].Score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}
}
