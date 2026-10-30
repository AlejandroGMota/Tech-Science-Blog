//go:build oracle

// Command migrate copies every article, rating and contact message from the
// Oracle Autonomous Database into a SQLite file, and writes a JSON dump of
// the same rows first so there is a backup that depends on neither database.
//
//	ORACLE_DB_DSN=... TNS_ADMIN=... migrate -sqlite data/blog.db -dump oracle-dump.json
//
// It refuses to write into a SQLite file that already holds articles, so
// running it twice cannot duplicate anything.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/godror/godror"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/store"
)

type dump struct {
	ExportedAt time.Time               `json:"exported_at"`
	Articles   []models.Article        `json:"articles"`
	Ratings    []models.Rating         `json:"ratings"`
	Contacts   []models.ContactMessage `json:"contacts"`
}

func main() {
	sqlitePath := flag.String("sqlite", "data/blog.db", "SQLite file to create")
	dumpPath := flag.String("dump", "oracle-dump.json", "JSON backup to write")
	flag.Parse()

	dsn := os.Getenv("ORACLE_DB_DSN")
	if dsn == "" {
		log.Fatal("ORACLE_DB_DSN is required")
	}
	ora, err := sql.Open("godror", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer ora.Close()
	if err := ora.Ping(); err != nil {
		log.Fatalf("oracle ping: %v", err)
	}

	d, err := read(ora)
	if err != nil {
		log.Fatalf("read oracle: %v", err)
	}
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*dumpPath, b, 0o600); err != nil {
		log.Fatal(err)
	}
	log.Printf("dump: %d articles, %d ratings, %d contacts → %s", len(d.Articles), len(d.Ratings), len(d.Contacts), *dumpPath)

	if err := write(*sqlitePath, d); err != nil {
		log.Fatalf("write sqlite: %v", err)
	}
	log.Printf("sqlite: copied into %s", *sqlitePath)
}

func read(db *sql.DB) (*dump, error) {
	d := &dump{ExportedAt: time.Now().UTC()}

	rows, err := db.Query(`SELECT slug, title, content, excerpt, category, article_type, tags, cover_image, author, published_at, created_at, updated_at FROM articles ORDER BY published_at`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a models.Article
		var content, excerpt, category, typ, tags, cover, author sql.NullString
		if err := rows.Scan(&a.Slug, &a.Title, &content, &excerpt, &category, &typ, &tags, &cover, &author, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.Content, a.Excerpt, a.Category, a.ArticleType, a.Tags, a.CoverImage, a.Author =
			content.String, excerpt.String, category.String, typ.String, tags.String, cover.String, author.String
		d.Articles = append(d.Articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	rows, err = db.Query(`SELECT article_slug, score, ip_hash, created_at FROM ratings ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.Rating
		if err := rows.Scan(&r.ArticleID, &r.Score, &r.IPHash, &r.CreatedAt); err != nil {
			return nil, err
		}
		d.Ratings = append(d.Ratings, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	rows, err = db.Query(`SELECT name, email, subject, message, created_at FROM contacts ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.ContactMessage
		var subject sql.NullString
		if err := rows.Scan(&c.Name, &c.Email, &subject, &c.Message, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.Subject = subject.String
		d.Contacts = append(d.Contacts, c)
	}
	return d, rows.Err()
}

func write(path string, d *dump) error {
	s, err := store.NewSQLiteStore(path)
	if err != nil {
		return err
	}
	defer s.Close()

	existing, err := s.GetArticles("", "", "")
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return fmt.Errorf("%s already has %d articles; refusing to copy twice", path, len(existing))
	}

	for i := range d.Articles {
		if err := s.CreateArticle(&d.Articles[i]); err != nil {
			return fmt.Errorf("article %s: %w", d.Articles[i].Slug, err)
		}
	}
	for i := range d.Ratings {
		r := d.Ratings[i]
		if err := s.RateArticle(r.ArticleID, &r); err != nil {
			return fmt.Errorf("rating %s: %w", r.ArticleID, err)
		}
	}
	for i := range d.Contacts {
		if err := s.CreateContact(&d.Contacts[i]); err != nil {
			return fmt.Errorf("contact %d: %w", i, err)
		}
	}
	return nil
}
