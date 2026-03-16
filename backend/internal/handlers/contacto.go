package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/store"
)

type ContactHandler struct {
	store store.Store
}

func NewContactHandler(s store.Store) *ContactHandler {
	return &ContactHandler{store: s}
}

// POST /api/contacto
func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c models.ContactMessage
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if c.Name == "" || c.Email == "" || c.Message == "" {
		writeError(w, http.StatusBadRequest, "name, email and message are required")
		return
	}

	if err := h.store.CreateContact(&c); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

// GET /api/contacto
func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.store.GetContacts()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if contacts == nil {
		contacts = []models.ContactMessage{}
	}
	writeJSON(w, http.StatusOK, contacts)
}

// DELETE /api/contacto/{id}
func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.store.DeleteContact(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
