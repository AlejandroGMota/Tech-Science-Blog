package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/store"
)

const domain = "https://blog.alejandrogmota.com"

type SEOHandler struct {
	store store.Store
}

func NewSEOHandler(s store.Store) *SEOHandler {
	return &SEOHandler{store: s}
}

// GET /robots.txt
func (h *SEOHandler) RobotsTXT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, `User-agent: *
Allow: /
Disallow: /admin/

Sitemap: %s/sitemap.xml
`, domain)
}

// GET /sitemap.xml
func (h *SEOHandler) SitemapXML(w http.ResponseWriter, r *http.Request) {
	type URL struct {
		Loc        string `xml:"loc"`
		LastMod    string `xml:"lastmod,omitempty"`
		ChangeFreq string `xml:"changefreq,omitempty"`
		Priority   string `xml:"priority,omitempty"`
	}
	type URLSet struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []URL    `xml:"url"`
	}

	urls := []URL{
		{Loc: domain + "/", ChangeFreq: "weekly", Priority: "1.0"},
		{Loc: domain + "/entradas", ChangeFreq: "daily", Priority: "0.9"},
	}

	articles, err := h.store.GetArticles("", "", "")
	if err == nil {
		for _, a := range articles {
			urls = append(urls, URL{
				Loc:        domain + "/entradas/" + a.Slug,
				LastMod:    a.UpdatedAt.Format(time.DateOnly),
				ChangeFreq: "monthly",
				Priority:   "0.8",
			})
		}
	}

	sitemap := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(sitemap)
}
