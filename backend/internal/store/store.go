package store

import "github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"

type Store interface {
	// Articles
	GetArticles(category, search string) ([]models.Article, error)
	GetArticleBySlug(slug string) (*models.Article, error)
	CreateArticle(a *models.Article) error
	UpdateArticle(slug string, a *models.Article) error
	DeleteArticle(slug string) error

	// Ratings
	RateArticle(slug string, r *models.Rating) error
	GetArticleRating(slug string) (*models.RatingSummary, error)

	// Contact
	CreateContact(c *models.ContactMessage) error
	GetContacts() ([]models.ContactMessage, error)
	DeleteContact(id string) error
}
