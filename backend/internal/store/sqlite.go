package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
)

// SQLiteStore keeps everything in one file, so moving the blog to another
// host is copying that file. The driver is pure Go: no cgo, no Instant Client.
type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	dsn := "file:" + path + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}
	// SQLite allows one writer at a time; a single connection avoids
	// SQLITE_BUSY between our own goroutines.
	db.SetMaxOpenConns(1)

	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("sqlite migrate: %w", err)
	}
	return s, nil
}

func (s *SQLiteStore) Close() error { return s.db.Close() }

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			excerpt TEXT NOT NULL DEFAULT '',
			category TEXT NOT NULL DEFAULT '',
			article_type TEXT NOT NULL DEFAULT 'deep_dive',
			tags TEXT NOT NULL DEFAULT '',
			cover_image TEXT NOT NULL DEFAULT '',
			author TEXT NOT NULL DEFAULT '',
			published_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_articles_published ON articles (published_at DESC);
		CREATE TABLE IF NOT EXISTS ratings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			article_slug TEXT NOT NULL,
			score INTEGER NOT NULL CHECK (score BETWEEN 1 AND 5),
			ip_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			UNIQUE (article_slug, ip_hash)
		);
		CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			subject TEXT NOT NULL DEFAULT '',
			message TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);`)
	return err
}

func isUnique(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}

// Articles

const articleCols = `id, slug, title, content, excerpt, category, article_type, tags, cover_image, author, published_at, created_at, updated_at`

func scanArticle(sc interface{ Scan(...any) error }) (models.Article, error) {
	var a models.Article
	var id int64
	err := sc.Scan(&id, &a.Slug, &a.Title, &a.Content, &a.Excerpt, &a.Category, &a.ArticleType, &a.Tags, &a.CoverImage, &a.Author, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	a.ID = fmt.Sprintf("%d", id)
	return a, err
}

func (s *SQLiteStore) GetArticles(category, search, articleType string) ([]models.Article, error) {
	query := `SELECT ` + articleCols + ` FROM articles WHERE 1=1`
	var args []any
	if category != "" {
		query += ` AND category = ? COLLATE NOCASE`
		args = append(args, category)
	}
	if search != "" {
		// LIKE is case-insensitive for ASCII in SQLite.
		query += ` AND (title LIKE ? OR tags LIKE ?)`
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if articleType != "" {
		query += ` AND article_type = ? COLLATE NOCASE`
		args = append(args, articleType)
	}
	query += ` ORDER BY published_at DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, rows.Err()
}

func (s *SQLiteStore) GetArticleBySlug(slug string) (*models.Article, error) {
	a, err := scanArticle(s.db.QueryRow(`SELECT `+articleCols+` FROM articles WHERE slug = ?`, slug))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found: %s", slug)
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *SQLiteStore) CreateArticle(a *models.Article) error {
	now := time.Now().UTC()
	if a.PublishedAt.IsZero() {
		a.PublishedAt = now
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = now
	}
	if a.ArticleType == "" {
		a.ArticleType = "deep_dive"
	}

	res, err := s.db.Exec(
		`INSERT INTO articles (slug, title, content, excerpt, category, article_type, tags, cover_image, author, published_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Slug, a.Title, a.Content, a.Excerpt, a.Category, a.ArticleType, a.Tags, a.CoverImage, a.Author, a.PublishedAt.UTC(), a.CreatedAt.UTC(), a.UpdatedAt.UTC(),
	)
	if isUnique(err) {
		return fmt.Errorf("article already exists: %s", a.Slug)
	}
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	a.ID = fmt.Sprintf("%d", id)
	return nil
}

func (s *SQLiteStore) UpdateArticle(slug string, a *models.Article) error {
	a.UpdatedAt = time.Now().UTC()
	res, err := s.db.Exec(
		`UPDATE articles SET title = ?, content = ?, excerpt = ?, category = ?, article_type = ?, tags = ?, cover_image = ?, author = ?, updated_at = ? WHERE slug = ?`,
		a.Title, a.Content, a.Excerpt, a.Category, a.ArticleType, a.Tags, a.CoverImage, a.Author, a.UpdatedAt, slug,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("article not found: %s", slug)
	}
	return nil
}

func (s *SQLiteStore) DeleteArticle(slug string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`DELETE FROM articles WHERE slug = ?`, slug)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("article not found: %s", slug)
	}
	if _, err := tx.Exec(`DELETE FROM ratings WHERE article_slug = ?`, slug); err != nil {
		return err
	}
	return tx.Commit()
}

// Ratings

func (s *SQLiteStore) RateArticle(slug string, r *models.Rating) error {
	r.ArticleID = slug
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.Exec(
		`INSERT INTO ratings (article_slug, score, ip_hash, created_at) VALUES (?, ?, ?, ?)`,
		slug, r.Score, r.IPHash, r.CreatedAt.UTC(),
	)
	if isUnique(err) {
		return fmt.Errorf("already rated")
	}
	return err
}

// GetArticleRating blends the all-time average with a recency-weighted one
// over the last 60 days, the same formula as the Oracle store. The weighting
// is done in Go because SQLite has no date arithmetic on TIMESTAMP values.
func (s *SQLiteStore) GetArticleRating(slug string) (*models.RatingSummary, error) {
	rows, err := s.db.Query(`SELECT score, created_at FROM ratings WHERE article_slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	now := time.Now()
	var sum, wSum, w float64
	var count, recentCount int
	for rows.Next() {
		var score int
		var at time.Time
		if err := rows.Scan(&score, &at); err != nil {
			return nil, err
		}
		count++
		sum += float64(score)
		if age := now.Sub(at).Hours() / 24; age <= 60 {
			weight := 1 - age/60
			wSum += float64(score) * weight
			w += weight
			recentCount++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	summary := &models.RatingSummary{Count: count}
	if count == 0 {
		return summary, nil
	}
	hist := sum / float64(count)
	summary.Average = hist
	if recentCount > 0 && w > 0 {
		alpha := float64(recentCount) / float64(count)
		if alpha > 0.7 {
			alpha = 0.7
		}
		summary.Average = alpha*(wSum/w) + (1-alpha)*hist
	}
	return summary, nil
}

// Contact

func (s *SQLiteStore) CreateContact(c *models.ContactMessage) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	res, err := s.db.Exec(
		`INSERT INTO contacts (name, email, subject, message, created_at) VALUES (?, ?, ?, ?, ?)`,
		c.Name, c.Email, c.Subject, c.Message, c.CreatedAt.UTC(),
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	c.ID = fmt.Sprintf("%d", id)
	return nil
}

func (s *SQLiteStore) GetContacts() ([]models.ContactMessage, error) {
	rows, err := s.db.Query(`SELECT id, name, email, subject, message, created_at FROM contacts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.ContactMessage
	for rows.Next() {
		var c models.ContactMessage
		var id int64
		if err := rows.Scan(&id, &c.Name, &c.Email, &c.Subject, &c.Message, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.ID = fmt.Sprintf("%d", id)
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func (s *SQLiteStore) DeleteContact(id string) error {
	res, err := s.db.Exec(`DELETE FROM contacts WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("contact not found: %s", id)
	}
	return nil
}
