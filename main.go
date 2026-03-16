package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/handlers"
	"github.com/petraclara/quality-education-EduMentor/middleware"
)

func main() {
	// Initialize database
	db, err := database.New("edumentor.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	profileHandler := handlers.NewProfileHandler(db)
	matchHandler := handlers.NewMatchHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)

	// Auth middleware
	authMW := middleware.Auth(db)

	// Create router
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/logout", authHandler.Logout)

	// Protected routes
	mux.Handle("/api/me", authMW(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/profile", authMW(http.HandlerFunc(profileHandler.HandleProfile)))
	mux.Handle("/api/matches/find", authMW(http.HandlerFunc(matchHandler.FindMatches)))
	mux.Handle("/api/matches", authMW(http.HandlerFunc(matchHandler.GetMyMatches)))
	mux.Handle("/api/dashboard", authMW(http.HandlerFunc(dashboardHandler.GetDashboard)))

	// Match action routes (accept/reject) - uses path prefix matching
	mux.Handle("/api/matches/", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 4 && (parts[3] == "accept" || parts[3] == "reject") {
			matchHandler.UpdateMatch(w, r)
			return
		}
		http.NotFound(w, r)
	})))

	// Apply CORS middleware
	handler := middleware.CORS(mux)

	fmt.Println("🎓 EduMentor API server starting on http://localhost:8080")
	fmt.Println("   Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(":8080", handler))
}