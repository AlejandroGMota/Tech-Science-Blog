package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/godror/godror"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
)

type OracleStore struct {
	db *sql.DB
}

func NewOracleStore(dsn string) (*OracleStore, error) {
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, fmt.Errorf("oracle connect: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("oracle ping: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	s := &OracleStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("oracle migrate: %w", err)
	}
	return s, nil
}

func (s *OracleStore) migrate() error {
	queries := []string{
		`CREATE TABLE articles (
			id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			slug VARCHAR2(500) UNIQUE NOT NULL,
			title VARCHAR2(500) NOT NULL,
			content CLOB,
			excerpt VARCHAR2(2000),
			category VARCHAR2(100),
			cover_image VARCHAR2(1000),
			author VARCHAR2(200),
			published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE ratings (
			id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			article_slug VARCHAR2(500) NOT NULL,
			score NUMBER(1) CHECK (score BETWEEN 1 AND 5),
			ip_hash VARCHAR2(64) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT uq_rating_ip UNIQUE (article_slug, ip_hash)
		)`,
		`CREATE TABLE contacts (
			id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name VARCHAR2(200) NOT NULL,
			email VARCHAR2(300) NOT NULL,
			subject VARCHAR2(500),
			message CLOB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		_, err := s.db.Exec(q)
		if err != nil {
			// ORA-00955: name is already used — table exists, skip
			if strings.Contains(err.Error(), "ORA-00955") {
				continue
			}
			return err
		}
	}
	return nil
}

// Articles

func (s *OracleStore) GetArticles(category, search string) ([]models.Article, error) {
	query := `SELECT id, slug, title, content, excerpt, category, cover_image, author, published_at, created_at, updated_at FROM articles WHERE 1=1`
	var args []any
	argIdx := 1

	if category != "" {
		query += fmt.Sprintf(" AND UPPER(category) = UPPER(:%d)", argIdx)
		args = append(args, category)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND UPPER(title) LIKE UPPER(:%d)", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	query += " ORDER BY published_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		var coverImg sql.NullString
		err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.Excerpt, &a.Category, &coverImg, &a.Author, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if coverImg.Valid {
			a.CoverImage = coverImg.String
		}
		articles = append(articles, a)
	}
	return articles, nil
}

func (s *OracleStore) GetArticleBySlug(slug string) (*models.Article, error) {
	var a models.Article
	var coverImg sql.NullString
	err := s.db.QueryRow(
		`SELECT id, slug, title, content, excerpt, category, cover_image, author, published_at, created_at, updated_at FROM articles WHERE slug = :1`,
		slug,
	).Scan(&a.ID, &a.Slug, &a.Title, &a.Content, &a.Excerpt, &a.Category, &coverImg, &a.Author, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found: %s", slug)
	}
	if err != nil {
		return nil, err
	}
	if coverImg.Valid {
		a.CoverImage = coverImg.String
	}
	return &a, nil
}

func (s *OracleStore) CreateArticle(a *models.Article) error {
	now := time.Now()
	if a.PublishedAt.IsZero() {
		a.PublishedAt = now
	}
	a.CreatedAt = now
	a.UpdatedAt = now

	var id int64
	_, err := s.db.Exec(
		`INSERT INTO articles (slug, title, content, excerpt, category, cover_image, author, published_at, created_at, updated_at)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)
		RETURNING id INTO :11`,
		a.Slug, a.Title, a.Content, a.Excerpt, a.Category, a.CoverImage, a.Author, a.PublishedAt, a.CreatedAt, a.UpdatedAt,
		sql.Out{Dest: &id},
	)
	if err != nil {
		if strings.Contains(err.Error(), "ORA-00001") {
			return fmt.Errorf("article already exists: %s", a.Slug)
		}
		return err
	}
	a.ID = fmt.Sprintf("%d", id)
	return nil
}

func (s *OracleStore) UpdateArticle(slug string, a *models.Article) error {
	a.UpdatedAt = time.Now()
	result, err := s.db.Exec(
		`UPDATE articles SET title = :1, content = :2, excerpt = :3, category = :4, cover_image = :5, author = :6, updated_at = :7 WHERE slug = :8`,
		a.Title, a.Content, a.Excerpt, a.Category, a.CoverImage, a.Author, a.UpdatedAt, slug,
	)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("article not found: %s", slug)
	}
	return nil
}

func (s *OracleStore) DeleteArticle(slug string) error {
	result, err := s.db.Exec(`DELETE FROM articles WHERE slug = :1`, slug)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("article not found: %s", slug)
	}
	// Also delete ratings
	s.db.Exec(`DELETE FROM ratings WHERE article_slug = :1`, slug)
	return nil
}

// Ratings

func (s *OracleStore) RateArticle(slug string, r *models.Rating) error {
	r.ArticleID = slug
	r.CreatedAt = time.Now()

	_, err := s.db.Exec(
		`INSERT INTO ratings (article_slug, score, ip_hash, created_at) VALUES (:1, :2, :3, :4)`,
		slug, r.Score, r.IPHash, r.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "ORA-00001") {
			return fmt.Errorf("already rated")
		}
		return err
	}
	return nil
}

func (s *OracleStore) GetArticleRating(slug string) (*models.RatingSummary, error) {
	var avg sql.NullFloat64
	var count int
	err := s.db.QueryRow(
		`SELECT AVG(score), COUNT(*) FROM ratings WHERE article_slug = :1`,
		slug,
	).Scan(&avg, &count)
	if err != nil {
		return nil, err
	}
	summary := &models.RatingSummary{Count: count}
	if avg.Valid {
		summary.Average = avg.Float64
	}
	return summary, nil
}

// Contact

func (s *OracleStore) CreateContact(c *models.ContactMessage) error {
	c.CreatedAt = time.Now()
	var id int64
	_, err := s.db.Exec(
		`INSERT INTO contacts (name, email, subject, message, created_at) VALUES (:1, :2, :3, :4, :5) RETURNING id INTO :6`,
		c.Name, c.Email, c.Subject, c.Message, c.CreatedAt,
		sql.Out{Dest: &id},
	)
	if err != nil {
		return err
	}
	c.ID = fmt.Sprintf("%d", id)
	return nil
}

func (s *OracleStore) GetContacts() ([]models.ContactMessage, error) {
	rows, err := s.db.Query(`SELECT id, name, email, subject, message, created_at FROM contacts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.ContactMessage
	for rows.Next() {
		var c models.ContactMessage
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Subject, &c.Message, &c.CreatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (s *OracleStore) DeleteContact(id string) error {
	result, err := s.db.Exec(`DELETE FROM contacts WHERE id = :1`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("contact not found: %s", id)
	}
	return nil
}
