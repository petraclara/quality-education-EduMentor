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
	db, err := database.New("edumentor.db")
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

	handler := middleware.CORS(mux)
	fmt.Println("🎓 EduMentor API server starting on http://localhost:8080")
	fmt.Println("   Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(":8080", handler))
}