package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/petraclara/quality-education-EduMentor/models"
	_ "modernc.org/sqlite"
)

// DB wraps the sql.DB connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection and runs migrations
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'mentee',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER UNIQUE NOT NULL,
			bio TEXT DEFAULT '',
			skills TEXT DEFAULT '[]',
			interests TEXT DEFAULT '[]',
			availability TEXT DEFAULT '[]',
			avatar_url TEXT DEFAULT '',
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS matches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			mentor_id INTEGER NOT NULL,
			mentee_id INTEGER NOT NULL,
			score REAL DEFAULT 0,
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (mentor_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (mentee_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(mentor_id, mentee_id)
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, q := range queries {
		if _, err := db.conn.Exec(q); err != nil {
			return fmt.Errorf("migration failed: %w\nQuery: %s", err, q)
		}
	}
	return nil
}

// ==================== User Operations ====================

// CreateUser inserts a new user and creates an empty profile
func (db *DB) CreateUser(name, email, passwordHash, role string) (*models.User, error) {
	result, err := db.conn.Exec(
		"INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)",
		name, email, passwordHash, role,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return nil, fmt.Errorf("email already registered")
		}
		return nil, err
	}

	id, _ := result.LastInsertId()
	user := &models.User{
		ID:    int(id),
		Name:  name,
		Email: email,
		Role:  role,
	}

	// Create empty profile
	_, err = db.conn.Exec(
		"INSERT INTO profiles (user_id) VALUES (?)", id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := db.conn.QueryRow(
		"SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (db *DB) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := db.conn.QueryRow(
		"SELECT id, name, email, password_hash, role, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ==================== Profile Operations ====================

// GetProfile retrieves a user's profile
func (db *DB) GetProfile(userID int) (*models.Profile, error) {
	profile := &models.Profile{}
	var skillsJSON, interestsJSON, availabilityJSON string

	err := db.conn.QueryRow(
		"SELECT id, user_id, bio, skills, interests, availability, avatar_url FROM profiles WHERE user_id = ?",
		userID,
	).Scan(&profile.ID, &profile.UserID, &profile.Bio, &skillsJSON, &interestsJSON, &availabilityJSON, &profile.AvatarURL)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(skillsJSON), &profile.Skills)
	json.Unmarshal([]byte(interestsJSON), &profile.Interests)
	json.Unmarshal([]byte(availabilityJSON), &profile.Availability)

	if profile.Skills == nil {
		profile.Skills = []string{}
	}
	if profile.Interests == nil {
		profile.Interests = []string{}
	}
	if profile.Availability == nil {
		profile.Availability = []string{}
	}

	return profile, nil
}

// UpdateProfile updates a user's profile
func (db *DB) UpdateProfile(userID int, req models.ProfileUpdateRequest) error {
	skillsJSON, _ := json.Marshal(req.Skills)
	interestsJSON, _ := json.Marshal(req.Interests)
	availabilityJSON, _ := json.Marshal(req.Availability)

	_, err := db.conn.Exec(
		"UPDATE profiles SET bio = ?, skills = ?, interests = ?, availability = ? WHERE user_id = ?",
		req.Bio, string(skillsJSON), string(interestsJSON), string(availabilityJSON), userID,
	)
	return err
}

// GetAllProfiles retrieves all profiles with user info (for matching)
func (db *DB) GetAllProfiles() ([]struct {
	User    models.User
	Profile models.Profile
}, error) {
	rows, err := db.conn.Query(`
		SELECT u.id, u.name, u.email, u.role, u.created_at,
		       p.id, p.user_id, p.bio, p.skills, p.interests, p.availability, p.avatar_url
		FROM users u
		JOIN profiles p ON u.id = p.user_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		User    models.User
		Profile models.Profile
	}

	for rows.Next() {
		var u models.User
		var p models.Profile
		var skillsJSON, interestsJSON, availabilityJSON string

		err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt,
			&p.ID, &p.UserID, &p.Bio, &skillsJSON, &interestsJSON, &availabilityJSON, &p.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(skillsJSON), &p.Skills)
		json.Unmarshal([]byte(interestsJSON), &p.Interests)
		json.Unmarshal([]byte(availabilityJSON), &p.Availability)

		if p.Skills == nil {
			p.Skills = []string{}
		}
		if p.Interests == nil {
			p.Interests = []string{}
		}
		if p.Availability == nil {
			p.Availability = []string{}
		}

		results = append(results, struct {
			User    models.User
			Profile models.Profile
		}{u, p})
	}
	return results, nil
}

// ==================== Match Operations ====================

// CreateMatch inserts a new match
func (db *DB) CreateMatch(mentorID, menteeID int, score float64) (*models.Match, error) {
	result, err := db.conn.Exec(
		"INSERT OR REPLACE INTO matches (mentor_id, mentee_id, score, status) VALUES (?, ?, ?, 'pending')",
		mentorID, menteeID, score,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.Match{
		ID:       int(id),
		MentorID: mentorID,
		MenteeID: menteeID,
		Score:    score,
		Status:   "pending",
	}, nil
}

// UpdateMatchStatus updates the status of a match
func (db *DB) UpdateMatchStatus(matchID int, status string) error {
	_, err := db.conn.Exec(
		"UPDATE matches SET status = ? WHERE id = ?",
		status, matchID,
	)
	return err
}

// GetMatchesByUser retrieves all matches for a user (as mentor or mentee) with user details
func (db *DB) GetMatchesByUser(userID int) ([]models.MatchWithUser, error) {
	rows, err := db.conn.Query(`
		SELECT m.id, m.mentor_id, m.mentee_id, m.score, m.status, m.created_at,
		       u.name, u.email, u.role,
		       p.bio, p.skills, p.interests
		FROM matches m
		JOIN users u ON (CASE WHEN m.mentor_id = ? THEN m.mentee_id ELSE m.mentor_id END) = u.id
		JOIN profiles p ON u.id = p.user_id
		WHERE m.mentor_id = ? OR m.mentee_id = ?
		ORDER BY m.score DESC
	`, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.MatchWithUser
	for rows.Next() {
		var m models.MatchWithUser
		var skillsJSON, interestsJSON string

		err := rows.Scan(
			&m.ID, &m.MentorID, &m.MenteeID, &m.Score, &m.Status, &m.CreatedAt,
			&m.UserName, &m.UserEmail, &m.UserRole,
			&m.Bio, &skillsJSON, &interestsJSON,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(skillsJSON), &m.Skills)
		json.Unmarshal([]byte(interestsJSON), &m.Interests)

		if m.Skills == nil {
			m.Skills = []string{}
		}
		if m.Interests == nil {
			m.Interests = []string{}
		}

		matches = append(matches, m)
	}
	return matches, nil
}

// GetMatchByID retrieves a match by ID
func (db *DB) GetMatchByID(id int) (*models.Match, error) {
	m := &models.Match{}
	err := db.conn.QueryRow(
		"SELECT id, mentor_id, mentee_id, score, status, created_at FROM matches WHERE id = ?",
		id,
	).Scan(&m.ID, &m.MentorID, &m.MenteeID, &m.Score, &m.Status, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetDashboardStats retrieves match statistics for a user
func (db *DB) GetDashboardStats(userID int) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	err := db.conn.QueryRow(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = 'accepted' THEN 1 ELSE 0 END) as accepted,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected
		FROM matches
		WHERE mentor_id = ? OR mentee_id = ?
	`, userID, userID).Scan(&stats.TotalMatches, &stats.PendingMatches, &stats.AcceptedMatches, &stats.RejectedMatches)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// ==================== Session Operations ====================

// CreateSession creates a new session
func (db *DB) CreateSession(token string, userID int, expiresAt time.Time) error {
	_, err := db.conn.Exec(
		"INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt,
	)
	return err
}

// GetSession retrieves a session by token
func (db *DB) GetSession(token string) (*models.Session, error) {
	s := &models.Session{}
	err := db.conn.QueryRow(
		"SELECT token, user_id, expires_at FROM sessions WHERE token = ?",
		token,
	).Scan(&s.Token, &s.UserID, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSession removes a session
func (db *DB) DeleteSession(token string) error {
	_, err := db.conn.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

// CleanExpiredSessions removes expired sessions
func (db *DB) CleanExpiredSessions() error {
	_, err := db.conn.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}
