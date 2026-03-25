package api

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/config"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/handlers"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/middleware"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/store"
)

func NewRouter(cfg *config.Config) http.Handler {
	// Storage
	var s store.Store
	switch cfg.DBType {
	case "oracle":
		oracleStore, err := store.NewOracleStore(cfg.OracleDBDSN)
		if err != nil {
			log.Fatalf("Failed to connect to Oracle DB: %v", err)
		}
		s = oracleStore
		log.Println("Using Oracle Autonomous Database")
	default:
		s = store.NewMemoryStore()
		store.SeedArticles(s)
		log.Println("Using in-memory storage (seeded with articles)")
	}

	// Session store
	sessions := middleware.NewSessionStore()
	auth := middleware.RequireAuth(sessions)

	// Handlers
	articleH := handlers.NewArticleHandler(s)
	ratingH := handlers.NewRatingHandler(s)
	contactH := handlers.NewContactHandler(s)
	seoH := handlers.NewSEOHandler(s)

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Auth
	mux.HandleFunc("POST /api/admin/login", loginHandler(cfg, sessions))
	mux.Handle("POST /api/admin/logout", auth(http.HandlerFunc(logoutHandler(sessions))))

	// Articles (public)
	mux.HandleFunc("GET /api/articles", articleH.List)
	mux.HandleFunc("GET /api/articles/{slug}", articleH.Get)

	// Articles (protected)
	mux.Handle("POST /api/articles", auth(http.HandlerFunc(articleH.Create)))
	mux.Handle("PUT /api/articles/{slug}", auth(http.HandlerFunc(articleH.Update)))
	mux.Handle("DELETE /api/articles/{slug}", auth(http.HandlerFunc(articleH.Delete)))

	// Ratings (public)
	mux.HandleFunc("GET /api/articles/{slug}/rating", ratingH.Get)
	mux.HandleFunc("POST /api/articles/{slug}/rating", ratingH.Rate)

	// Contact
	mux.HandleFunc("POST /api/contacto", contactH.Create)
	mux.Handle("GET /api/contacto", auth(http.HandlerFunc(contactH.List)))
	mux.Handle("DELETE /api/contacto/{id}", auth(http.HandlerFunc(contactH.Delete)))

	// SEO
	mux.HandleFunc("GET /robots.txt", seoH.RobotsTXT)
	mux.HandleFunc("GET /sitemap.xml", seoH.SitemapXML)

	// Serve frontend SPA
	mux.HandleFunc("GET /admin/", serveSPA("admin-dist"))
	publicSPA := handlers.NewSPAHandler(s, "public-dist")
	mux.Handle("GET /", publicSPA)

	// Apply CORS
	return middleware.CORS(cfg.AllowedOrigins)(mux)
}

func serveSPA(distDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := filepath.Join(distDir, r.URL.Path)
		if distDir == "admin-dist" {
			path = filepath.Join(distDir, r.URL.Path[len("/admin"):])
		}

		// Clean path to prevent traversal
		path = filepath.Clean(path)

		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		// Check if distDir exists
		if _, err := fs.Stat(os.DirFS("."), distDir); err != nil {
			http.NotFound(w, r)
			return
		}

		// Fallback to index.html (SPA routing)
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	}
}

func loginHandler(cfg *config.Config, sessions *middleware.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			User string `json:"user"`
			Pass string `json:"pass"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if creds.User != cfg.AdminUser || creds.Pass != cfg.AdminPass {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		token, err := sessions.Create()
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func logoutHandler(sessions *middleware.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")[len("Bearer "):]
		sessions.Delete(token)
		w.WriteHeader(http.StatusNoContent)
	}
}
