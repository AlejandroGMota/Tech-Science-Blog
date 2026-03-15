package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/store"
)

type RatingHandler struct {
	store store.Store
}

func NewRatingHandler(s store.Store) *RatingHandler {
	return &RatingHandler{store: s}
}

// GET /api/articles/{slug}/rating
func (h *RatingHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	summary, err := h.store.GetArticleRating(slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// POST /api/articles/{slug}/rating
func (h *RatingHandler) Rate(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	var input struct {
		Score int `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if input.Score < 1 || input.Score > 5 {
		writeError(w, http.StatusBadRequest, "score must be between 1 and 5")
		return
	}

	rating := &models.Rating{
		Score:  input.Score,
		IPHash: hashIP(r),
	}

	if err := h.store.RateArticle(slug, rating); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	summary, _ := h.store.GetArticleRating(slug)
	writeJSON(w, http.StatusCreated, summary)
}

func hashIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	hash := sha256.Sum256([]byte(ip + "tech-blog-salt"))
	return fmt.Sprintf("%x", hash[:16])
}
