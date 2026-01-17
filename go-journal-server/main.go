package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go-journal-server/database"
	"go-journal-server/handlers"
)

func main() {
	// Initialize database
	db, err := database.NewJournalDB("./data")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create handler
	handler := handlers.NewHandler(db)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Setup routes
	handlers.SetupRoutes(r, handler)

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Journal API Server"))
	})

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
