package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/petraclara/quality-education-EduMentor/database"
	"github.com/petraclara/quality-education-EduMentor/handlers"
	"github.com/petraclara/quality-education-EduMentor/middleware"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "edumentor.db"
	}

	db, err := database.New(dbUrl)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	if err := db.SeedDemoData(); err != nil {
		log.Println("Warning: failed to seed demo data:", err)
	}

	authHandler := handlers.NewAuthHandler(db)
	profileHandler := handlers.NewProfileHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)
	mentorHandler := handlers.NewMentorHandler(db)
	requestHandler := handlers.NewRequestHandler(db)
	bookingHandler := handlers.NewBookingHandler(db)

	authMW := middleware.Auth(db)
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/logout", authHandler.Logout)
	mux.HandleFunc("/api/mentors", mentorHandler.ListMentors)
	mux.HandleFunc("/api/mentors/", mentorHandler.GetMentor)

	// Protected routes
	mux.Handle("/api/me", authMW(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/profile", authMW(http.HandlerFunc(profileHandler.HandleProfile)))
	mux.Handle("/api/dashboard", authMW(http.HandlerFunc(dashboardHandler.GetDashboard)))
	mux.Handle("/api/requests", authMW(http.HandlerFunc(requestHandler.HandleRequests)))
	mux.Handle("/api/bookings", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookingHandler.GetBookings(w, r)
		case http.MethodPost:
			bookingHandler.CreateBooking(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Request action routes: /api/requests/{id}/{action} where action is accept/decline/confirm
	mux.Handle("/api/requests/", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 4 && (parts[3] == "accept" || parts[3] == "decline" || parts[3] == "confirm") {
			requestHandler.HandleRequestAction(w, r)
			return
		}
		http.NotFound(w, r)
	})))

	// Static Frontend Server
	// Serve static files from the frontend/dist directory
	distDir := "./frontend/dist"
	fs := http.FileServer(http.Dir(distDir))
	
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If it's explicitly an API route that wasn't matched above, let it 404 cleanly in API
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		
		// If requesting the root, serve index
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
			return
		}
		
		// Serve static request if the file actually exists
		path := filepath.Join(distDir, r.URL.Path)
		if _, err := os.Stat(path); err == nil {
			fs.ServeHTTP(w, r)
			return
		}
		
		// For React Router single-page apps, serve index.html for all other paths
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	}))

	handler := middleware.CORS(mux)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("🎓 MentorConnect API server starting on http://localhost:%s\n", port)
	fmt.Println("   Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}