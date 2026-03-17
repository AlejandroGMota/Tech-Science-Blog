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
