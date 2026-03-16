package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/store"
)

type ArticleHandler struct {
	store store.Store
}

func NewArticleHandler(s store.Store) *ArticleHandler {
	return &ArticleHandler{store: s}
}

// GET /api/articles
func (h *ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	articles, err := h.store.GetArticles(category, search)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if articles == nil {
		articles = []models.Article{}
	}
	writeJSON(w, http.StatusOK, articles)
}

// GET /api/articles/{slug}
func (h *ArticleHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	article, err := h.store.GetArticleBySlug(slug)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, article)
}

// POST /api/articles
func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var a models.Article
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if a.Title == "" || a.Slug == "" {
		writeError(w, http.StatusBadRequest, "title and slug are required")
		return
	}

	if err := h.store.CreateArticle(&a); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

// PUT /api/articles/{slug}
func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	var a models.Article
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	a.Slug = slug
	if err := h.store.UpdateArticle(slug, &a); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// DELETE /api/articles/{slug}
func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	if err := h.store.DeleteArticle(slug); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
