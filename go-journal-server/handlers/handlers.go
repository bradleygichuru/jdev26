package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go-journal-server/database"
)

type Handler struct {
	db *database.JournalDB
}

func NewHandler(db *database.JournalDB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Content == "" {
		h.sendError(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	entry, err := h.db.CreateEntry(req.Title, req.Content, req.Tags)
	if err != nil {
		h.sendError(w, "Failed to create entry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.convertToAPIEntry(entry)
	h.sendResponse(w, response, http.StatusCreated)
}

func (h *Handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.db.GetEntry(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendError(w, "Entry not found", http.StatusNotFound)
		} else {
			h.sendError(w, "Failed to get entry: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := h.convertToAPIEntry(entry)
	h.sendResponse(w, response, http.StatusOK)
}

func (h *Handler) GetAllEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := h.db.GetAllEntries()
	if err != nil {
		h.sendError(w, "Failed to get entries: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response []JournalEntry
	for _, entry := range entries {
		response = append(response, *h.convertToAPIEntry(entry))
	}

	h.sendResponse(w, response, http.StatusOK)
}

func (h *Handler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	var req UpdateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.db.UpdateEntry(id, req.Title, req.Content, req.Tags)
	if err != nil {
		h.sendError(w, "Failed to update entry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated entry
	entry, err := h.db.GetEntry(id)
	if err != nil {
		h.sendError(w, "Failed to get updated entry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.convertToAPIEntry(entry)
	h.sendResponse(w, response, http.StatusOK)
}

func (h *Handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	err = h.db.DeleteEntry(id)
	if err != nil {
		h.sendError(w, "Failed to delete entry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, map[string]string{"message": "Entry deleted successfully"}, http.StatusOK)
}

func (h *Handler) SearchEntries(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.sendError(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Get all entries and filter in application layer since LIKE is not supported
	allEntries, err := h.db.GetAllEntries()
	if err != nil {
		h.sendError(w, "Failed to search entries: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var matchedEntries []*database.JournalEntryDB
	queryLower := strings.ToLower(query)

	for _, entry := range allEntries {
		if strings.Contains(strings.ToLower(entry.Title), queryLower) ||
			strings.Contains(strings.ToLower(entry.Content), queryLower) ||
			strings.Contains(strings.ToLower(entry.Tags), queryLower) {
			matchedEntries = append(matchedEntries, entry)
		}
	}

	var response []JournalEntry
	for _, entry := range matchedEntries {
		response = append(response, *h.convertToAPIEntry(entry))
	}

	h.sendResponse(w, response, http.StatusOK)
}

func (h *Handler) convertToAPIEntry(dbEntry *database.JournalEntryDB) *JournalEntry {
	tags := []string{}
	if dbEntry.Tags != "" {
		tags = strings.Split(dbEntry.Tags, ",")
	}
	// Filter out empty tags
	var filtered []string
	for _, t := range tags {
		if t != "" {
			filtered = append(filtered, t)
		}
	}

	return &JournalEntry{
		ID:        dbEntry.ID,
		Title:     dbEntry.Title,
		Content:   dbEntry.Content,
		CreatedAt: dbEntry.CreatedAt,
		UpdatedAt: dbEntry.UpdatedAt,
		Tags:      filtered,
	}
}

func (h *Handler) sendResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Success: true,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIResponse{
		Success: false,
		Error:   message,
	}

	json.NewEncoder(w).Encode(response)
}
