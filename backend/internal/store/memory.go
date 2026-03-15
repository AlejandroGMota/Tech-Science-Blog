package store

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
)

type MemoryStore struct {
	mu       sync.RWMutex
	articles map[string]*models.Article
	ratings  map[string][]models.Rating // key: article slug
	contacts map[string]*models.ContactMessage
	nextID   int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		articles: make(map[string]*models.Article),
		ratings:  make(map[string][]models.Rating),
		contacts: make(map[string]*models.ContactMessage),
		nextID:   1,
	}
}

func (s *MemoryStore) genID() string {
	id := fmt.Sprintf("%d", s.nextID)
	s.nextID++
	return id
}

// Articles

func (s *MemoryStore) GetArticles(category, search string) ([]models.Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Article
	for _, a := range s.articles {
		if category != "" && !strings.EqualFold(a.Category, category) {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(a.Title), strings.ToLower(search)) {
			continue
		}
		result = append(result, *a)
	}
	return result, nil
}

func (s *MemoryStore) GetArticleBySlug(slug string) (*models.Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	a, ok := s.articles[slug]
	if !ok {
		return nil, fmt.Errorf("article not found: %s", slug)
	}
	return a, nil
}

func (s *MemoryStore) CreateArticle(a *models.Article) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.articles[a.Slug]; exists {
		return fmt.Errorf("article already exists: %s", a.Slug)
	}

	a.ID = s.genID()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	if a.PublishedAt.IsZero() {
		a.PublishedAt = now
	}
	s.articles[a.Slug] = a
	return nil
}

func (s *MemoryStore) UpdateArticle(slug string, a *models.Article) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.articles[slug]
	if !ok {
		return fmt.Errorf("article not found: %s", slug)
	}

	a.ID = existing.ID
	a.CreatedAt = existing.CreatedAt
	a.UpdatedAt = time.Now()
	s.articles[slug] = a
	return nil
}

func (s *MemoryStore) DeleteArticle(slug string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.articles[slug]; !ok {
		return fmt.Errorf("article not found: %s", slug)
	}
	delete(s.articles, slug)
	delete(s.ratings, slug)
	return nil
}

// Ratings

func (s *MemoryStore) RateArticle(slug string, r *models.Rating) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.articles[slug]; !ok {
		return fmt.Errorf("article not found: %s", slug)
	}

	// Check duplicate by IP hash
	for _, existing := range s.ratings[slug] {
		if existing.IPHash == r.IPHash {
			return fmt.Errorf("already rated")
		}
	}

	r.ID = s.genID()
	r.ArticleID = slug
	r.CreatedAt = time.Now()
	s.ratings[slug] = append(s.ratings[slug], *r)
	return nil
}

func (s *MemoryStore) GetArticleRating(slug string) (*models.RatingSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ratings := s.ratings[slug]
	if len(ratings) == 0 {
		return &models.RatingSummary{Average: 0, Count: 0}, nil
	}

	var total int
	for _, r := range ratings {
		total += r.Score
	}
	return &models.RatingSummary{
		Average: float64(total) / float64(len(ratings)),
		Count:   len(ratings),
	}, nil
}

// Contact

func (s *MemoryStore) CreateContact(c *models.ContactMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c.ID = s.genID()
	c.CreatedAt = time.Now()
	s.contacts[c.ID] = c
	return nil
}

func (s *MemoryStore) GetContacts() ([]models.ContactMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.ContactMessage
	for _, c := range s.contacts {
		result = append(result, *c)
	}
	return result, nil
}

func (s *MemoryStore) DeleteContact(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.contacts[id]; !ok {
		return fmt.Errorf("contact not found: %s", id)
	}
	delete(s.contacts, id)
	return nil
}
