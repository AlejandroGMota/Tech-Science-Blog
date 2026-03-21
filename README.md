# Blog — alejandrogmota.com

Un blog personal de alguien que construye cosas y aprende en el proceso. Tecnología, negocios, ideas y todo lo que vale la pena documentar.

**Live:** [blog.alejandrogmota.com](https://blog.alejandrogmota.com)

---

## Filosofía

> Publica lo que te hubiera servido a ti hace 6 meses.

El blog también alimenta tu portafolio y posicionamiento profesional. Un artículo técnico bueno vale más que 10 líneas en un CV.

---

## Categorías

| Categoría | Contenido |
|-----------|-----------|
| **Code** | Arquitectura, microservicios, Go, Node.js, DevOps, lecciones técnicas reales del trabajo |
| **Business** | Celinki, emprendimiento en México, importación, lecciones de escalar un negocio, SAT, operaciones |
| **Ideas** | Geopolítica, ciencia, análisis de cosas que te hacen pensar |
| **Stack de vida** | Productividad, herramientas, productos tech útiles, organización del tiempo |
| **Aprendiendo en público** | Cosas nuevas, errores, descubrimientos — el cajón comodín con honestidad |

## Formatos de artículo (ArticleType)

| Tipo | Extensión | Descripción |
|------|-----------|-------------|
| **Deep Dive** | ~1000-1500 palabras | Cuando dominas el tema o investigaste a fondo |
| **Nota Rápida** | ~300-500 palabras | Una lección, un error, una opinión. Publicar rápido sin perfeccionismo |
| **TIL (Today I Learned)** | Párrafo corto | Lo que aprendiste hoy. Bajo esfuerzo, alta consistencia |

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
| GET | `/api/articles` | List articles (`?category=X&search=Y&article_type=Z`) |
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
