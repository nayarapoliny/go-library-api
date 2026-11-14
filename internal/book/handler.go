package book

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/books", h.handleSearch)
	mux.HandleFunc("POST /api/v1/books", h.handleCreate)
	mux.HandleFunc("GET /api/v1/books/{id}", h.handleGetByID)
	mux.HandleFunc("DELETE /api/v1/books/{id}", h.handleDelete)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	year, _ := strconv.Atoi(q.Get("year"))

	filter := SearchFilter{
		Title:         q.Get("title"),
		Author:        q.Get("author"),
		Genre:         q.Get("genre"),
		Publisher:     q.Get("publisher"),
		PublishedYear: year,
	}

	books, err := h.service.SearchBooks(r.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, books)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}

	if err := h.service.CreateBook(r.Context(), &b); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrInvalidInput) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (h *Handler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := h.service.GetBook(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteBook(r.Context(), id); err != nil {
		if errors.Is(err, ErrBookNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
