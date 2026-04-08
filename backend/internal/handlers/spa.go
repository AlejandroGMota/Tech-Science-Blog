package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/store"
)

type SPAHandler struct {
	store   store.Store
	distDir string
}

func NewSPAHandler(s store.Store, distDir string) *SPAHandler {
	return &SPAHandler{store: s, distDir: distDir}
}

func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Try to serve static file directly
	path := filepath.Clean(filepath.Join(h.distDir, r.URL.Path))
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	// Check if distDir exists
	if _, err := os.Stat(h.distDir); err != nil {
		http.NotFound(w, r)
		return
	}

	// Read index.html
	indexPath := filepath.Join(h.distDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	htmlStr := string(content)

	// Check if this is an article page
	if slug := extractSlug(r.URL.Path); slug != "" {
		if article, err := h.store.GetArticleBySlug(slug); err == nil {
			htmlStr = injectArticleMeta(htmlStr, article.Title, article.Excerpt, article.Slug, article.CoverImage, article.Author, article.PublishedAt.Format("2006-01-02"), article.UpdatedAt.Format("2006-01-02"), article.Category, article.Content)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlStr))
}

func extractSlug(path string) string {
	path = strings.TrimSuffix(path, "/")
	if strings.HasPrefix(path, "/entradas/") {
		parts := strings.SplitN(path, "/entradas/", 2)
		if len(parts) == 2 && parts[1] != "" && !strings.Contains(parts[1], "/") {
			return parts[1]
		}
	}
	return ""
}

func injectArticleMeta(htmlStr, title, excerpt, slug, coverImage, author, publishedAt, updatedAt, category, content string) string {
	safeTitle := html.EscapeString(title)
	safeExcerpt := html.EscapeString(excerpt)
	articleURL := domain + "/entradas/" + slug

	// Replace title
	htmlStr = strings.Replace(htmlStr,
		"<title>Alejandro G. Mota — Blog</title>",
		fmt.Sprintf("<title>%s — Alejandro G. Mota</title>", safeTitle), 1)

	// Replace meta description
	htmlStr = strings.Replace(htmlStr,
		`<meta name="description" content="Blog personal de Alejandro G. Mota. Code, Business, Ideas y todo lo que vale la pena documentar." />`,
		fmt.Sprintf(`<meta name="description" content="%s" />`, safeExcerpt), 1)

	// Replace canonical
	htmlStr = strings.Replace(htmlStr,
		`<link rel="canonical" href="https://blog.alejandrogmota.com/" />`,
		fmt.Sprintf(`<link rel="canonical" href="%s" />`, articleURL), 1)

	// Replace OG tags
	htmlStr = strings.Replace(htmlStr,
		`<meta property="og:type" content="website" />`,
		`<meta property="og:type" content="article" />`, 1)
	htmlStr = strings.Replace(htmlStr,
		`<meta property="og:title" content="Alejandro G. Mota — Blog" />`,
		fmt.Sprintf(`<meta property="og:title" content="%s" />`, safeTitle), 1)
	htmlStr = strings.Replace(htmlStr,
		`<meta property="og:description" content="Code, Business, Ideas y todo lo que vale la pena documentar." />`,
		fmt.Sprintf(`<meta property="og:description" content="%s" />`, safeExcerpt), 1)
	htmlStr = strings.Replace(htmlStr,
		`<meta property="og:url" content="https://blog.alejandrogmota.com/" />`,
		fmt.Sprintf(`<meta property="og:url" content="%s" />`, articleURL), 1)

	// Replace Twitter tags
	htmlStr = strings.Replace(htmlStr,
		`<meta name="twitter:title" content="Alejandro G. Mota — Blog" />`,
		fmt.Sprintf(`<meta name="twitter:title" content="%s" />`, safeTitle), 1)
	htmlStr = strings.Replace(htmlStr,
		`<meta name="twitter:description" content="Code, Business, Ideas y todo lo que vale la pena documentar." />`,
		fmt.Sprintf(`<meta name="twitter:description" content="%s" />`, safeExcerpt), 1)

	// Add og:image and twitter:image if cover image exists
	var imageTag string
	if coverImage != "" {
		imageTag = fmt.Sprintf("\n    <meta property=\"og:image\" content=\"%s\" />\n    <meta name=\"twitter:image\" content=\"%s\" />", html.EscapeString(coverImage), html.EscapeString(coverImage))
	}

	// Build JSON-LD
	jsonLD := map[string]interface{}{
		"@context":      "https://schema.org",
		"@type":         "BlogPosting",
		"headline":      title,
		"description":   excerpt,
		"url":           articleURL,
		"author": map[string]string{
			"@type": "Person",
			"name":  author,
			"url":   "https://alejandrogmota.com",
		},
		"datePublished": publishedAt,
		"dateModified":  updatedAt,
		"publisher": map[string]interface{}{
			"@type": "Organization",
			"name":  "Alejandro G. Mota",
			"logo": map[string]interface{}{
				"@type": "ImageObject",
				"url":   "https://blog.alejandrogmota.com/favicon-64.png",
			},
		},
	}
	if coverImage != "" {
		jsonLD["image"] = coverImage
	}
	if category != "" {
		jsonLD["articleSection"] = category
	}

	jsonLDBytes, _ := json.Marshal(jsonLD)
	scriptTag := fmt.Sprintf(`<script type="application/ld+json">%s</script>`, string(jsonLDBytes))

	// Inject image tags and JSON-LD before </head>
	injection := imageTag + "\n    " + scriptTag
	htmlStr = strings.Replace(htmlStr, "</head>", injection+"\n  </head>", 1)

	// Inject article content into root div for crawler visibility
	htmlStr = strings.Replace(htmlStr,
		`<div id="root"></div>`,
		fmt.Sprintf(`<div id="root"><article class="crawlable-content">%s</article></div>`, content), 1)

	return htmlStr
}
