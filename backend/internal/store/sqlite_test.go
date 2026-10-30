package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/models"
)

func newTestSQLite(t *testing.T) *SQLiteStore {
	t.Helper()
	s, err := NewSQLiteStore(filepath.Join(t.TempDir(), "blog.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestSQLiteArticles(t *testing.T) {
	s := newTestSQLite(t)

	old := &models.Article{Slug: "old", Title: "Docker basics", Category: "DevOps", Tags: "docker", PublishedAt: time.Now().Add(-48 * time.Hour)}
	newer := &models.Article{Slug: "new", Title: "Go channels", Category: "Go", ArticleType: "short"}
	for _, a := range []*models.Article{old, newer} {
		if err := s.CreateArticle(a); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.CreateArticle(&models.Article{Slug: "old", Title: "dup"}); err == nil {
		t.Fatal("duplicate slug accepted")
	}

	all, err := s.GetArticles("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].Slug != "new" {
		t.Fatalf("want newest first, got %+v", all)
	}
	if all[1].ArticleType != "deep_dive" {
		t.Errorf("default article_type = %q", all[1].ArticleType)
	}

	for _, tc := range []struct{ cat, q, typ, want string }{
		{cat: "devops", want: "old"},
		{q: "DOCKER", want: "old"},
		{typ: "SHORT", want: "new"},
	} {
		got, err := s.GetArticles(tc.cat, tc.q, tc.typ)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].Slug != tc.want {
			t.Errorf("filter %+v: got %+v", tc, got)
		}
	}

	newer.Title = "Go channels, revisited"
	if err := s.UpdateArticle("new", newer); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetArticleBySlug("new")
	if err != nil || got.Title != "Go channels, revisited" {
		t.Fatalf("update not persisted: %+v, %v", got, err)
	}
	if err := s.UpdateArticle("missing", newer); err == nil {
		t.Error("update of missing article succeeded")
	}

	if err := s.RateArticle("old", &models.Rating{Score: 4, IPHash: "a"}); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteArticle("old"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetArticleBySlug("old"); err == nil {
		t.Error("deleted article still readable")
	}
	if r, _ := s.GetArticleRating("old"); r.Count != 0 {
		t.Errorf("ratings survived article delete: %+v", r)
	}
}

func TestSQLiteRatings(t *testing.T) {
	s := newTestSQLite(t)

	if r, err := s.GetArticleRating("x"); err != nil || r.Count != 0 || r.Average != 0 {
		t.Fatalf("empty rating = %+v, %v", r, err)
	}
	old := time.Now().Add(-100 * 24 * time.Hour)
	if err := s.RateArticle("x", &models.Rating{Score: 1, IPHash: "a", CreatedAt: old}); err != nil {
		t.Fatal(err)
	}
	if err := s.RateArticle("x", &models.Rating{Score: 5, IPHash: "b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.RateArticle("x", &models.Rating{Score: 3, IPHash: "b"}); err == nil {
		t.Fatal("second rating from the same ip accepted")
	}

	r, err := s.GetArticleRating("x")
	if err != nil {
		t.Fatal(err)
	}
	// hist = 3, recent = 5, alpha = 1/2 → 4.
	if r.Count != 2 || r.Average < 3.99 || r.Average > 4.01 {
		t.Errorf("rating = %+v, want count 2 average 4", r)
	}
}

func TestSQLiteContacts(t *testing.T) {
	s := newTestSQLite(t)

	c := &models.ContactMessage{Name: "Ana", Email: "ana@example.com", Message: "hola"}
	if err := s.CreateContact(c); err != nil {
		t.Fatal(err)
	}
	list, err := s.GetContacts()
	if err != nil || len(list) != 1 || list[0].ID != c.ID {
		t.Fatalf("contacts = %+v, %v", list, err)
	}
	if err := s.DeleteContact(c.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteContact(c.ID); err == nil {
		t.Error("second delete succeeded")
	}
}
