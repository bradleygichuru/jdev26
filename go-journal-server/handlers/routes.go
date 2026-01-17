package handlers

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux, handler *Handler) {
	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/entries", handler.CreateEntry)
		r.Get("/entries", handler.GetAllEntries)
		r.Get("/entries/{id}", handler.GetEntry)
		r.Put("/entries/{id}", handler.UpdateEntry)
		r.Delete("/entries/{id}", handler.DeleteEntry)
		r.Get("/entries/search", handler.SearchEntries)
	})
}
