# Tech & Science Blog

A full-stack blog platform for technology and science articles, built with React and Go.

**Live:** [blog.alejandrogmota.com](https://blog.alejandrogmota.com)

---

## Tech Stack

- **Frontend:** React 19 + Vite (public SPA + admin SPA)
- **Backend:** Go (net/http stdlib)
- **Database:** Oracle Autonomous JSON DB (godror driver)
- **Editor:** TipTap v2 (rich text, admin panel)
- **Infrastructure:** Oracle Cloud Free Tier (ARM VM, 1 OCPU / 6GB)
- **CI/CD:** GitHub Actions -> SSH deploy to VM
- **Reverse Proxy:** Nginx + Let's Encrypt SSL

---

## Architecture

```
blog.alejandrogmota.com
        |
      nginx (SSL, port 443)
        |
   Go API (:8081)
   |-- /api/*        -> REST endpoints
   |-- /             -> public-dist/ (React SPA)
   +-- /admin/       -> admin-dist/ (React SPA)
        |
  Oracle Autonomous DB (techblogdb)
```

---

## API Endpoints

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/articles` | List articles |
| GET | `/api/articles/{slug}` | Get article by slug |
| GET | `/api/articles/{slug}/rating` | Get article rating |
| POST | `/api/articles/{slug}/rating` | Submit rating (1-5) |
| POST | `/api/contacto` | Send contact message |

### Admin (Bearer token auth)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/login` | Login |
| POST | `/api/admin/logout` | Logout |
| POST | `/api/articles` | Create article |
| PUT | `/api/articles/{slug}` | Update article |
| DELETE | `/api/articles/{slug}` | Delete article |
| GET | `/api/contacto` | List contact messages |
| DELETE | `/api/contacto/{id}` | Delete contact message |

---

## Article Loading

### Current flow

```
EntradasPage → GET /api/articles?category=X&search=Y
               → Store.GetArticles() (no limit, ordered by published_at DESC)
               → returns all articles in one response

ArticlePage  → GET /api/articles/{slug}
             → GET /api/articles/{slug}/rating  (separate request)
```

### Bottlenecks

| Area | Issue |
|------|-------|
| No pagination | Full table scan returned every time |
| No HTTP cache headers | Every navigation re-fetches from DB |
| Dual rating request | Article detail fires 2 sequential API calls |
| No client cache | React state reset on route change |

### Optimizations (planned)

1. **Pagination** — `?page=1&limit=10` on `GET /api/articles` (already in Roadmap)
2. **HTTP caching** — `Cache-Control: public, max-age=60` on list endpoint; `stale-while-revalidate` for article detail
3. **Embed rating in article** — merge `RatingSummary` into `Article` response to eliminate the extra request on detail page
4. **Frontend SWR** — use a simple cache map in `useApi.js` keyed by URL, invalidated on mutations
5. **DB index** — ensure `published_at DESC` and `category` columns are indexed (Oracle auto-indexes PKs/unique; `category` needs manual index)

---

## Local Development

### Backend
```bash
cd backend
cp .env.example .env
go run ./cmd/api/
```

### Frontend (public)
```bash
cd frontend
npm install
npm run dev          # public SPA on :5173
```

### Frontend (admin)
```bash
cd frontend
npx vite --config vite.admin.config.js   # admin SPA on :5174
```

---

## Deployment

Push to `main` branch triggers automatic deployment via GitHub Actions:

1. SSH into Oracle Cloud VM
2. `git pull origin main`
3. Generate `.env` from GitHub Secrets
4. Build frontend (public + admin SPAs)
5. Build Go binary
6. Restart systemd service

### GitHub Secrets Required

| Secret | Description |
|--------|-------------|
| `VM_HOST` | VM public IP |
| `VM_USER` | SSH username |
| `VM_SSH_KEY` | SSH private key |
| `ADMIN_USER` | Admin panel username |
| `ADMIN_PASS` | Admin panel password |
| `ORACLE_DB_DSN` | Oracle DB connection string |

---

## Roadmap

- [ ] **Image Upload** — OCI Object Storage integration, drag-and-drop in editor
- [ ] **Pagination** — `?page=1&limit=10` for articles list
- [ ] **Draft/Published** — save articles as draft, toggle status from admin
- [ ] **SEO** — dynamic sitemap.xml, robots.txt, Open Graph meta tags

---

## License

MIT
