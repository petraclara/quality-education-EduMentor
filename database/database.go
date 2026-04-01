package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/petraclara/quality-education-EduMentor/models"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	conn   *sql.DB
	driver string
}

func New(dbPath string) (*DB, error) {
	driver := "sqlite"
	if strings.HasPrefix(dbPath, "postgres://") || strings.HasPrefix(dbPath, "postgresql://") {
		driver = "postgres"
	}

	conn, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if driver == "sqlite" {
		if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
			return nil, fmt.Errorf("failed to set WAL mode: %w", err)
		}
	}
	db := &DB{conn: conn, driver: driver}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return db, nil
}

func (db *DB) Close() error { return db.conn.Close() }

// formatQuery replaces ? with $1, $2, etc for postgres
func (db *DB) formatQuery(q string) string {
	if db.driver != "postgres" {
		return q
	}
	count := 1
	for {
		idx := strings.Index(q, "?")
		if idx == -1 {
			break
		}
		q = q[:idx] + fmt.Sprintf("$%d", count) + q[idx+1:]
		count++
	}
	return q
}

// executeInsert handles differences in sqlite vs postgres insertion and ID retrieval
func (db *DB) executeInsert(query string, args ...interface{}) (int64, error) {
	query = db.formatQuery(query)
	if db.driver == "postgres" {
		query += " RETURNING id"
		var id int64
		err := db.conn.QueryRow(query, args...).Scan(&id)
		return id, err
	}
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'learner',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER UNIQUE NOT NULL,
			bio TEXT DEFAULT '',
			skills TEXT DEFAULT '[]',
			interests TEXT DEFAULT '[]',
			level TEXT DEFAULT '',
			goal TEXT DEFAULT '',
			availability TEXT DEFAULT '[]',
			avatar_url TEXT DEFAULT '',
			max_mentees INTEGER DEFAULT 5,
			rating REAL DEFAULT 0,
			rating_count INTEGER DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS mentorship_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			learner_id INTEGER NOT NULL,
			mentor_id INTEGER NOT NULL,
			help_with TEXT NOT NULL,
			goal TEXT DEFAULT '',
			message TEXT DEFAULT '',
			status TEXT DEFAULT 'pending',
			decline_reason TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (learner_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (mentor_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			mentor_id INTEGER NOT NULL,
			mentee_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			time_slot TEXT NOT NULL,
			status TEXT DEFAULT 'upcoming',
			note TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (mentor_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (mentee_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}
	for _, q := range queries {
		if db.driver == "postgres" {
			q = strings.ReplaceAll(q, "INTEGER PRIMARY KEY AUTOINCREMENT", "SERIAL PRIMARY KEY")
			q = strings.ReplaceAll(q, "DATETIME", "TIMESTAMP")
		}
		if _, err := db.conn.Exec(q); err != nil {
			return fmt.Errorf("migration failed: %w\nQuery: %s", err, q)
		}
	}

	// Add new columns to mentorship_requests if they don't exist
	alterQueries := []string{
		`ALTER TABLE mentorship_requests ADD COLUMN meeting_type TEXT DEFAULT ''`,
		`ALTER TABLE mentorship_requests ADD COLUMN meeting_link TEXT DEFAULT ''`,
		`ALTER TABLE mentorship_requests ADD COLUMN proposed_slots TEXT DEFAULT '[]'`,
	}
	for _, q := range alterQueries {
		db.conn.Exec(q) // ignore errors since column might already exist
	}

	return nil
}

// ==================== Seed Data ====================

func (db *DB) SeedDemoData() error {
	var count int
	db.conn.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'mentor'").Scan(&count)
	if count > 0 {
		return nil
	}

	mentors := []struct {
		Name         string
		Email        string
		Bio          string
		Skills       []string
		Interests    []string
		Level        string
		Availability []string
		Rating       float64
		RatingCount  int
	}{
		{
			Name: "Dr. Amara Okonkwo", Email: "amara@edumentor.io",
			Bio:   "Senior software engineer with 12 years of experience in distributed systems. Previously led engineering teams at major tech companies. Passionate about helping the next generation of engineers grow.",
			Skills: []string{"Go", "Python", "Kubernetes", "AWS", "System Design"},
			Interests: []string{"Cloud Computing", "Open Source", "Tech Leadership"},
			Level: "beginner", Availability: []string{"monday-morning", "wednesday-morning", "friday-morning"},
			Rating: 4.9, RatingCount: 47,
		},
		{
			Name: "Raj Patel", Email: "raj@edumentor.io",
			Bio:   "Full-stack developer and educator specializing in modern web technologies. I run workshops on React and Node.js. My goal is to make complex concepts simple and accessible.",
			Skills: []string{"React", "JavaScript", "TypeScript", "Node.js", "CSS"},
			Interests: []string{"Web Development", "UI Design", "Teaching"},
			Level: "beginner", Availability: []string{"tuesday-morning", "thursday-morning", "saturday-morning"},
			Rating: 4.7, RatingCount: 35,
		},
		{
			Name: "Lin Wei Chen", Email: "lin@edumentor.io",
			Bio:   "Data scientist at a leading AI research lab. PhD in Machine Learning from MIT. I love breaking down complex ML concepts and helping students build real-world data pipelines.",
			Skills: []string{"Python", "Machine Learning", "TensorFlow", "Data Science", "SQL"},
			Interests: []string{"Artificial Intelligence", "Deep Learning", "Data Visualization"},
			Level: "intermediate", Availability: []string{"monday-evening", "wednesday-evening", "friday-afternoon"},
			Rating: 4.8, RatingCount: 52,
		},
		{
			Name: "Marcus Thompson", Email: "marcus@edumentor.io",
			Bio:   "DevOps engineer and infrastructure architect. Extensive experience in CI/CD pipelines, container orchestration, and site reliability. Mentor at several coding bootcamps.",
			Skills: []string{"Docker", "Kubernetes", "Terraform", "Linux", "AWS", "Go"},
			Interests: []string{"Infrastructure", "Automation", "DevOps"},
			Level: "intermediate", Availability: []string{"tuesday-evening", "thursday-afternoon", "saturday-afternoon"},
			Rating: 4.6, RatingCount: 28,
		},
		{
			Name: "Sofia Rodriguez", Email: "sofia@edumentor.io",
			Bio:   "UX/UI designer turned product manager. 8 years creating user-centered digital experiences. I help aspiring designers build portfolios and develop design thinking skills.",
			Skills: []string{"UI Design", "UX Research", "Figma", "Prototyping", "CSS"},
			Interests: []string{"Product Design", "User Research", "Design Thinking"},
			Level: "beginner", Availability: []string{"monday-morning", "wednesday-afternoon", "friday-morning"},
			Rating: 4.8, RatingCount: 41,
		},
		{
			Name: "Dr. Eleanor Hughes", Email: "eleanor@edumentor.io",
			Bio:   "Professor of Computer Science with 20 years of academic and industry experience. Specializing in algorithms and cybersecurity. Published 50+ peer-reviewed papers.",
			Skills: []string{"Algorithms", "Cybersecurity", "Java", "C++", "Cryptography"},
			Interests: []string{"Computer Science Education", "Security Research"},
			Level: "advanced", Availability: []string{"tuesday-morning", "wednesday-morning", "friday-afternoon"},
			Rating: 4.5, RatingCount: 19,
		},
		{
			Name: "David Kimani", Email: "david@edumentor.io",
			Bio:   "Mobile app developer with experience shipping iOS and Android apps to millions of users. Former lead at a fintech startup. I love helping people build their first apps.",
			Skills: []string{"Swift", "Kotlin", "React Native", "Flutter", "Firebase"},
			Interests: []string{"Mobile Apps", "Startup Culture", "Fintech"},
			Level: "beginner", Availability: []string{"monday-evening", "thursday-evening", "saturday-morning"},
			Rating: 4.7, RatingCount: 33,
		},
		{
			Name: "Yusuf Al-Rashid", Email: "yusuf@edumentor.io",
			Bio:   "Backend engineer specializing in high-performance APIs and microservices. Expert in Go and Rust. I enjoy pair programming and code reviews to accelerate learning.",
			Skills: []string{"Go", "Rust", "PostgreSQL", "Redis", "Microservices", "API Design"},
			Interests: []string{"Systems Programming", "Performance", "Open Source"},
			Level: "advanced", Availability: []string{"monday-morning", "tuesday-evening", "wednesday-morning", "thursday-morning"},
			Rating: 4.9, RatingCount: 38,
		},
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("mentor123"), bcrypt.DefaultCost)
	for _, m := range mentors {
		id, err := db.executeInsert(
			"INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, 'mentor')",
			m.Name, m.Email, string(hash),
		)
		if err != nil {
			continue
		}
		skillsJSON, _ := json.Marshal(m.Skills)
		interestsJSON, _ := json.Marshal(m.Interests)
		availJSON, _ := json.Marshal(m.Availability)
		db.conn.Exec(
			db.formatQuery(`INSERT INTO profiles (user_id, bio, skills, interests, level, availability, rating, rating_count) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`),
			id, m.Bio, string(skillsJSON), string(interestsJSON), m.Level, string(availJSON), m.Rating, m.RatingCount,
		)
	}
	fmt.Printf("✅ Seeded %d demo mentors\n", len(mentors))
	return nil
}

// ==================== User Operations ====================

func (db *DB) CreateUser(name, email, passwordHash, role string) (*models.User, error) {
	id, err := db.executeInsert(
		"INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)",
		name, email, passwordHash, role,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, fmt.Errorf("email already registered")
		}
		return nil, err
	}
	user := &models.User{ID: int(id), Name: name, Email: email, Role: role}
	_, err = db.conn.Exec(db.formatQuery("INSERT INTO profiles (user_id) VALUES (?)"), id)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := db.conn.QueryRow(
		db.formatQuery("SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = ?"), email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil { return nil, err }
	return user, nil
}

func (db *DB) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := db.conn.QueryRow(
		db.formatQuery("SELECT id, name, email, password_hash, role, created_at FROM users WHERE id = ?"), id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil { return nil, err }
	return user, nil
}

// ==================== Profile Operations ====================

func (db *DB) GetProfile(userID int) (*models.Profile, error) {
	profile := &models.Profile{}
	var skillsJSON, interestsJSON, availabilityJSON string
	err := db.conn.QueryRow(
		db.formatQuery(`SELECT id, user_id, bio, skills, interests, level, goal, availability, avatar_url, max_mentees, rating, rating_count 
		 FROM profiles WHERE user_id = ?`), userID,
	).Scan(&profile.ID, &profile.UserID, &profile.Bio, &skillsJSON, &interestsJSON,
		&profile.Level, &profile.Goal, &availabilityJSON, &profile.AvatarURL,
		&profile.MaxMentees, &profile.Rating, &profile.RatingCount)
	if err != nil { return nil, err }

	json.Unmarshal([]byte(skillsJSON), &profile.Skills)
	json.Unmarshal([]byte(interestsJSON), &profile.Interests)
	json.Unmarshal([]byte(availabilityJSON), &profile.Availability)
	if profile.Skills == nil { profile.Skills = []string{} }
	if profile.Interests == nil { profile.Interests = []string{} }
	if profile.Availability == nil { profile.Availability = []string{} }
	return profile, nil
}

func (db *DB) UpdateProfile(userID int, req models.ProfileUpdateRequest) error {
	skillsJSON, _ := json.Marshal(req.Skills)
	interestsJSON, _ := json.Marshal(req.Interests)
	availabilityJSON, _ := json.Marshal(req.Availability)
	_, err := db.conn.Exec(
		db.formatQuery(`UPDATE profiles SET bio = ?, skills = ?, interests = ?, level = ?, goal = ?, 
		 availability = ?, max_mentees = ? WHERE user_id = ?`),
		req.Bio, string(skillsJSON), string(interestsJSON), req.Level, req.Goal,
		string(availabilityJSON), req.MaxMentees, userID,
	)
	return err
}

// HasPreferencesSet checks if a user has filled in their preferences (skills or interests)
func (db *DB) HasPreferencesSet(userID int) (bool, error) {
	var skills, interests string
	err := db.conn.QueryRow(
		db.formatQuery("SELECT skills, interests FROM profiles WHERE user_id = ?"), userID,
	).Scan(&skills, &interests)
	if err != nil {
		return false, err
	}
	// Default value is '[]', so check for non-empty arrays
	return (skills != "[]" && skills != "") || (interests != "[]" && interests != ""), nil
}

// ==================== Mentor Operations ====================

func (db *DB) GetMentors() ([]models.MentorCard, error) {
	rows, err := db.conn.Query(`
		SELECT u.id, u.name, u.role, p.bio, p.skills, p.interests, p.level, p.availability, p.avatar_url, p.rating, p.rating_count
		FROM users u JOIN profiles p ON u.id = p.user_id
		WHERE u.role = 'mentor' ORDER BY p.rating DESC
	`)
	if err != nil { return nil, err }
	defer rows.Close()

	var mentors []models.MentorCard
	for rows.Next() {
		var m models.MentorCard
		var skillsJSON, interestsJSON, availJSON string
		err := rows.Scan(&m.ID, &m.Name, &m.Role, &m.Bio, &skillsJSON, &interestsJSON, &m.Level, &availJSON, &m.AvatarURL, &m.Rating, &m.RatingCount)
		if err != nil { return nil, err }
		json.Unmarshal([]byte(skillsJSON), &m.Skills)
		json.Unmarshal([]byte(interestsJSON), &m.Interests)
		json.Unmarshal([]byte(availJSON), &m.Availability)
		if m.Skills == nil { m.Skills = []string{} }
		if m.Interests == nil { m.Interests = []string{} }
		if m.Availability == nil { m.Availability = []string{} }
		mentors = append(mentors, m)
	}
	return mentors, nil
}

func (db *DB) GetMentorByID(id int) (*models.MentorCard, error) {
	m := &models.MentorCard{}
	var skillsJSON, interestsJSON, availJSON string
	err := db.conn.QueryRow(db.formatQuery(`
		SELECT u.id, u.name, u.role, p.bio, p.skills, p.interests, p.level, p.availability, p.avatar_url, p.rating, p.rating_count
		FROM users u JOIN profiles p ON u.id = p.user_id WHERE u.id = ?
	`), id).Scan(&m.ID, &m.Name, &m.Role, &m.Bio, &skillsJSON, &interestsJSON, &m.Level, &availJSON, &m.AvatarURL, &m.Rating, &m.RatingCount)
	if err != nil { return nil, err }
	json.Unmarshal([]byte(skillsJSON), &m.Skills)
	json.Unmarshal([]byte(interestsJSON), &m.Interests)
	json.Unmarshal([]byte(availJSON), &m.Availability)
	if m.Skills == nil { m.Skills = []string{} }
	if m.Interests == nil { m.Interests = []string{} }
	if m.Availability == nil { m.Availability = []string{} }
	return m, nil
}

// GetMatchedMentors returns mentors ranked by compatibility with a learner's interests
func (db *DB) GetMatchedMentors(learnerID int) ([]models.MentorCard, error) {
	learnerProfile, err := db.GetProfile(learnerID)
	if err != nil { return nil, err }

	allMentors, err := db.GetMentors()
	if err != nil { return nil, err }

	var matches []models.MentorCard

	// Score each mentor based on interest/skill overlap
	for i := range allMentors {
		score := 0.0
		score += jaccardSimilarity(learnerProfile.Interests, allMentors[i].Skills) * 0.5
		score += jaccardSimilarity(learnerProfile.Interests, allMentors[i].Interests) * 0.3
		// level match bonus
		if learnerProfile.Level != "" && allMentors[i].Level == learnerProfile.Level {
			score += 0.2
		}
		
		// Only include mentors with some overlap in skills or interests
		// Alternatively, just requiring score > 0 after the interest/skill checks is good.
		// Wait, level match gives 0.2, so a mentor with no skill match but same level could get a score of 0.2.
		// The user explicitly requested: "what user wants to learn matches what mentor is teaching"
		skillOverlap := jaccardSimilarity(learnerProfile.Interests, allMentors[i].Skills)
		interestOverlap := jaccardSimilarity(learnerProfile.Interests, allMentors[i].Interests)
		
		if skillOverlap > 0 || interestOverlap > 0 {
			allMentors[i].MatchScore = score
			matches = append(matches, allMentors[i])
		}
	}

	// Sort by match score descending
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].MatchScore > matches[i].MatchScore {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}
	return matches, nil
}

func jaccardSimilarity(a, b []string) float64 {
	if len(a) == 0 && len(b) == 0 { return 0 }
	setA := make(map[string]bool)
	for _, v := range a { setA[strings.ToLower(strings.TrimSpace(v))] = true }
	setB := make(map[string]bool)
	for _, v := range b { setB[strings.ToLower(strings.TrimSpace(v))] = true }
	intersection := 0
	for k := range setA { if setB[k] { intersection++ } }
	union := len(setA)
	for k := range setB { if !setA[k] { union++ } }
	if union == 0 { return 0 }
	return float64(intersection) / float64(union)
}

// ==================== Mentorship Request Operations ====================

func (db *DB) CreateMentorshipRequest(learnerID, mentorID int, helpWith, goal, message string) (*models.MentorshipRequest, error) {
	id, err := db.executeInsert(
		"INSERT INTO mentorship_requests (learner_id, mentor_id, help_with, goal, message) VALUES (?, ?, ?, ?, ?)",
		learnerID, mentorID, helpWith, goal, message,
	)
	if err != nil { return nil, err }
	return &models.MentorshipRequest{
		ID: int(id), LearnerID: learnerID, MentorID: mentorID,
		HelpWith: helpWith, Goal: goal, Message: message, Status: "pending",
	}, nil
}

func (db *DB) GetRequestsForMentor(mentorID int) ([]models.MentorshipRequestWithUser, error) {
	rows, err := db.conn.Query(db.formatQuery(`
		SELECT r.id, r.learner_id, r.mentor_id, r.help_with, r.goal, r.message, r.status, r.decline_reason, r.meeting_type, r.meeting_link, r.proposed_slots, r.created_at,
		       learner.name, learner.email, COALESCE(p.level, '')
		FROM mentorship_requests r
		JOIN users learner ON r.learner_id = learner.id
		LEFT JOIN profiles p ON learner.id = p.user_id
		WHERE r.mentor_id = ?
		ORDER BY CASE r.status WHEN 'pending' THEN 0 WHEN 'accepted' THEN 1 ELSE 2 END, r.created_at DESC
	`), mentorID)
	if err != nil { return nil, err }
	defer rows.Close()

	var requests []models.MentorshipRequestWithUser
	for rows.Next() {
		var r models.MentorshipRequestWithUser
		var slotsJSON string
		err := rows.Scan(&r.ID, &r.LearnerID, &r.MentorID, &r.HelpWith, &r.Goal, &r.Message,
			&r.Status, &r.DeclineReason, &r.MeetingType, &r.MeetingLink, &slotsJSON, &r.CreatedAt, &r.LearnerName, &r.LearnerEmail, &r.LearnerLevel)
		if err != nil { return nil, err }
		json.Unmarshal([]byte(slotsJSON), &r.ProposedSlots)
		if r.ProposedSlots == nil { r.ProposedSlots = []models.ProposedSlot{} }
		requests = append(requests, r)
	}
	if requests == nil { requests = []models.MentorshipRequestWithUser{} }
	return requests, nil
}

func (db *DB) GetRequestsByLearner(learnerID int) ([]models.MentorshipRequestWithUser, error) {
	rows, err := db.conn.Query(db.formatQuery(`
		SELECT r.id, r.learner_id, r.mentor_id, r.help_with, r.goal, r.message, r.status, r.decline_reason, r.meeting_type, r.meeting_link, r.proposed_slots, r.created_at,
		       '', '', '', mentor.name
		FROM mentorship_requests r
		JOIN users mentor ON r.mentor_id = mentor.id
		WHERE r.learner_id = ?
		ORDER BY r.created_at DESC
	`), learnerID)
	if err != nil { return nil, err }
	defer rows.Close()

	var requests []models.MentorshipRequestWithUser
	for rows.Next() {
		var r models.MentorshipRequestWithUser
		var slotsJSON string
		err := rows.Scan(&r.ID, &r.LearnerID, &r.MentorID, &r.HelpWith, &r.Goal, &r.Message,
			&r.Status, &r.DeclineReason, &r.MeetingType, &r.MeetingLink, &slotsJSON, &r.CreatedAt, &r.LearnerName, &r.LearnerEmail, &r.LearnerLevel, &r.MentorName)
		if err != nil { return nil, err }
		json.Unmarshal([]byte(slotsJSON), &r.ProposedSlots)
		if r.ProposedSlots == nil { r.ProposedSlots = []models.ProposedSlot{} }
		requests = append(requests, r)
	}
	if requests == nil { requests = []models.MentorshipRequestWithUser{} }
	return requests, nil
}

func (db *DB) GetRequestByID(id int) (*models.MentorshipRequest, error) {
	r := &models.MentorshipRequest{}
	var slotsJSON string
	err := db.conn.QueryRow(
		db.formatQuery("SELECT id, learner_id, mentor_id, help_with, goal, message, status, decline_reason, meeting_type, meeting_link, proposed_slots, created_at FROM mentorship_requests WHERE id = ?"), id,
	).Scan(&r.ID, &r.LearnerID, &r.MentorID, &r.HelpWith, &r.Goal, &r.Message, &r.Status, &r.DeclineReason, &r.MeetingType, &r.MeetingLink, &slotsJSON, &r.CreatedAt)
	if err != nil { return nil, err }
	json.Unmarshal([]byte(slotsJSON), &r.ProposedSlots)
	if r.ProposedSlots == nil { r.ProposedSlots = []models.ProposedSlot{} }
	return r, nil
}

func (db *DB) AcceptRequest(id int, meetingType, meetingLink string, proposedSlots []models.ProposedSlot) error {
	slotsJSON, _ := json.Marshal(proposedSlots)
	_, err := db.conn.Exec(db.formatQuery("UPDATE mentorship_requests SET status = 'accepted', meeting_type = ?, meeting_link = ?, proposed_slots = ? WHERE id = ?"), meetingType, meetingLink, string(slotsJSON), id)
	return err
}

func (db *DB) ConfirmRequestSlot(reqID int, date, timeStr string) error {
	req, err := db.GetRequestByID(reqID)
	if err != nil { return err }
	
	// Overwrite proposed_slots with just the single confirmed slot so UI can easily show it
	confirmedSlot := []models.ProposedSlot{{Date: date, Time: timeStr}}
	slotsJSON, _ := json.Marshal(confirmedSlot)
	
	_, err = db.conn.Exec(db.formatQuery("UPDATE mentorship_requests SET status = 'scheduled', proposed_slots = ? WHERE id = ?"), string(slotsJSON), reqID)
	if err != nil { return err }
	
	_, err = db.CreateBooking(req.MentorID, req.LearnerID, date, timeStr, "Mentorship session via "+req.MeetingType+": "+req.MeetingLink)
	return err
}

func (db *DB) DeclineRequest(id int, reason string) error {
	_, err := db.conn.Exec(db.formatQuery("UPDATE mentorship_requests SET status = 'declined', decline_reason = ? WHERE id = ?"), reason, id)
	return err
}

// ==================== Booking Operations ====================

func (db *DB) CreateBooking(mentorID, menteeID int, date, timeSlot, note string) (*models.Booking, error) {
	id, err := db.executeInsert(
		"INSERT INTO bookings (mentor_id, mentee_id, date, time_slot, note) VALUES (?, ?, ?, ?, ?)",
		mentorID, menteeID, date, timeSlot, note,
	)
	if err != nil { return nil, err }
	return &models.Booking{ID: int(id), MentorID: mentorID, MenteeID: menteeID, Date: date, TimeSlot: timeSlot, Status: "upcoming", Note: note}, nil
}

func (db *DB) GetBookingsByUser(userID int) ([]models.BookingWithUser, error) {
	rows, err := db.conn.Query(db.formatQuery(`
		SELECT b.id, b.mentor_id, b.mentee_id, b.date, b.time_slot, b.status, b.note, b.created_at,
		       mentor.name, mentee.name
		FROM bookings b
		JOIN users mentor ON b.mentor_id = mentor.id
		JOIN users mentee ON b.mentee_id = mentee.id
		WHERE b.mentor_id = ? OR b.mentee_id = ?
		ORDER BY b.date ASC
	`), userID, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var bookings []models.BookingWithUser
	for rows.Next() {
		var b models.BookingWithUser
		rows.Scan(&b.ID, &b.MentorID, &b.MenteeID, &b.Date, &b.TimeSlot, &b.Status, &b.Note, &b.CreatedAt, &b.MentorName, &b.MenteeName)
		bookings = append(bookings, b)
	}
	if bookings == nil { bookings = []models.BookingWithUser{} }
	return bookings, nil
}

func (db *DB) GetBookingsForMentorOnDate(mentorID int, date string) ([]models.Booking, error) {
	rows, err := db.conn.Query(
		db.formatQuery("SELECT id, mentor_id, mentee_id, date, time_slot, status, note, created_at FROM bookings WHERE mentor_id = ? AND date = ? AND status = 'upcoming'"),
		mentorID, date,
	)
	if err != nil { return nil, err }
	defer rows.Close()
	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		rows.Scan(&b.ID, &b.MentorID, &b.MenteeID, &b.Date, &b.TimeSlot, &b.Status, &b.Note, &b.CreatedAt)
		bookings = append(bookings, b)
	}
	if bookings == nil { bookings = []models.Booking{} }
	return bookings, nil
}

// ==================== Session Operations ====================

func (db *DB) CreateSession(token string, userID int, expiresAt time.Time) error {
	_, err := db.conn.Exec(db.formatQuery("INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)"), token, userID, expiresAt)
	return err
}

func (db *DB) GetSession(token string) (*models.Session, error) {
	s := &models.Session{}
	err := db.conn.QueryRow(db.formatQuery("SELECT token, user_id, expires_at FROM sessions WHERE token = ?"), token).Scan(&s.Token, &s.UserID, &s.ExpiresAt)
	if err != nil { return nil, err }
	return s, nil
}

func (db *DB) DeleteSession(token string) error {
	_, err := db.conn.Exec(db.formatQuery("DELETE FROM sessions WHERE token = ?"), token)
	return err
}

func (db *DB) CleanExpiredSessions() error {
	_, err := db.conn.Exec(db.formatQuery("DELETE FROM sessions WHERE expires_at < ?"), time.Now())
	return err
}
