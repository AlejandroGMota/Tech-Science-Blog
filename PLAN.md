# Tech-Science-Blog Restructure Plan

## Current State
- Static HTML/CSS/JS blog (no framework, no build tools)
- 10 articles as individual HTML files in `/entradas/`
- Article metadata in `entradas.json` (title, date, category, stars)
- Hosted on GitHub Pages → `blog.alejandrogmota.com`
- Google Analytics, AdSense, PWA manifest
- SEO optimized (Schema.org, Open Graph, sitemap)

## Goals
1. **Migrate to a framework** (React frontend + Go backend, monorepo)
2. **Host on Oracle Cloud VM** (free tier, same pattern as comerciantesChavarria)
3. **Admin page to publish new articles** without editing raw HTML
4. **Individual article rating/scoring system** (separate function per article)
5. **Free tier storage** for images and database

---

## Phase 1: Infrastructure — Oracle Cloud VM (Free Tier)

### VM Provisioning (desde consola de Oracle)

- [ ] Crear cuenta en Oracle Cloud (si no existe) o reutilizar la existente
- [ ] Crear instancia VM.Standard.A1.Flex (ARM) en Oracle Cloud Console:
  - **Shape:** VM.Standard.A1.Flex — 1 OCPU, 6 GB RAM (Always Free)
  - **Boot Volume:** 50 GB (Always Free)
  - **Region:** mx-queretaro-1 (o la más cercana)
  - **OS:** Ubuntu 22.04 LTS (Canonical)
  - **SSH Key:** Generar par de llaves o usar existente
- [ ] Configurar **VCN (Virtual Cloud Network)** desde la consola:
  - Crear VCN con CIDR 10.0.0.0/16
  - Crear subnet pública
  - Configurar Internet Gateway
  - Configurar Route Table (0.0.0.0/0 → Internet Gateway)
- [ ] Configurar **Security List** (firewall de Oracle):
  - Ingress: TCP 22 (SSH), TCP 80 (HTTP), TCP 443 (HTTPS), TCP 8080 (API dev)
  - Egress: All traffic allowed
- [ ] Asignar **IP pública reservada** (Always Free: 1 Reserved Public IP)
- [ ] Conectar por SSH y configurar el servidor:
  ```bash
  ssh -i <key> ubuntu@<IP>
  ```

### Configuración del Servidor

- [ ] Actualizar sistema: `sudo apt update && sudo apt upgrade -y`
- [ ] Instalar dependencias:
  ```bash
  sudo apt install -y git curl ufw nginx certbot python3-certbot-nginx
  ```
- [ ] Instalar Go (1.22+):
  ```bash
  wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
  sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
  ```
- [ ] Instalar Node.js 20 (para builds del frontend):
  ```bash
  curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
  sudo apt install -y nodejs
  ```
- [ ] Configurar UFW (firewall interno):
  ```bash
  sudo ufw allow 22/tcp
  sudo ufw allow 80/tcp
  sudo ufw allow 443/tcp
  sudo ufw enable
  ```
- [ ] Configurar Nginx como reverse proxy:
  ```nginx
  server {
      listen 80;
      server_name blog.alejandrogmota.com;

      location / {
          proxy_pass http://localhost:8080;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
          proxy_set_header X-Forwarded-Proto $scheme;
      }
  }
  ```
- [ ] Configurar SSL con Let's Encrypt:
  ```bash
  sudo certbot --nginx -d blog.alejandrogmota.com
  ```
- [ ] Actualizar DNS: cambiar CNAME de GitHub Pages a **A record** apuntando a la IP de Oracle

### Servicio Systemd

- [ ] Crear servicio `tech-blog.service`:
  ```ini
  [Unit]
  Description=Tech Science Blog API
  After=network.target

  [Service]
  Type=simple
  User=ubuntu
  WorkingDirectory=/home/ubuntu/Tech-Science-Blog/backend
  ExecStart=/home/ubuntu/Tech-Science-Blog/backend/tech-blog-api
  Restart=always
  RestartSec=5
  Environment=PORT=8080
  Environment=ALLOWED_ORIGINS=https://blog.alejandrogmota.com
  Environment=ADMIN_USER=${ADMIN_USER}
  Environment=ADMIN_PASS=${ADMIN_PASS}
  Environment=DB_TYPE=supabase
  Environment=SUPABASE_URL=${SUPABASE_URL}
  Environment=SUPABASE_KEY=${SUPABASE_KEY}

  [Install]
  WantedBy=multi-user.target
  ```

---

## Phase 2: Framework Migration (React + Go Monorepo)

### Arquitectura (mismo patrón que comerciantesChavarria)

**Same-Origin Architecture:**
```
Browser Request
    ├─ GET /                    → frontend-dist/index.html (React SPA)
    ├─ GET /entradas            → React Router (SPA fallback)
    ├─ GET /admin               → admin-dist/index.html (Admin SPA)
    ├─ GET /api/articles        → Go handler
    ├─ POST /api/articles       → Go handler (auth required)
    ├─ POST /api/ratings        → Go handler
    └─ GET /api/health          → Go handler
```

### Backend (Go)

- [ ] Inicializar módulo Go: `go mod init github.com/AlejandroGMota/Tech-Science-Blog`
- [ ] Crear estructura del backend:
  ```
  backend/
  ├── cmd/api/main.go              # Entry point
  ├── api/routes.go                # HTTP router + SPA handler
  ├── internal/
  │   ├── config/config.go         # Env vars
  │   ├── models/models.go         # Article, Rating, User models
  │   ├── store/
  │   │   ├── store.go             # Storage interface
  │   │   ├── memory.go            # In-memory (dev)
  │   │   └── supabase.go          # Supabase/Postgres (prod)
  │   ├── handlers/
  │   │   ├── articles.go          # CRUD artículos
  │   │   ├── ratings.go           # Sistema de puntuación
  │   │   ├── contacto.go          # Formulario de contacto
  │   │   └── helpers.go           # JSON writer, validación
  │   └── middleware/
  │       ├── auth.go              # Session/Bearer token auth
  │       └── cors.go              # CORS middleware
  └── public-dist/                 # Frontend compilado (gitignored)
  ```

#### API Endpoints

**Públicos (sin auth):**
- `GET /api/health` — Health check
- `GET /api/articles` — Listar artículos (con filtros: category, date, search)
- `GET /api/articles/{slug}` — Obtener artículo individual
- `GET /api/articles/{slug}/rating` — Obtener rating de un artículo
- `POST /api/articles/{slug}/rating` — Votar (1-5 estrellas)
- `POST /api/contacto` — Enviar mensaje de contacto

**Protegidos (Bearer token):**
- `POST /api/admin/login` — Login
- `POST /api/admin/logout` — Logout
- `POST /api/articles` — Crear artículo
- `PUT /api/articles/{slug}` — Editar artículo
- `DELETE /api/articles/{slug}` — Eliminar artículo
- `POST /api/articles/upload-image` — Subir imagen (→ Cloudflare R2)
- `GET /api/contacto` — Ver mensajes (admin)
- `DELETE /api/contacto/{id}` — Eliminar mensaje

#### Modelos de Datos

```go
type Article struct {
    ID          string    `json:"id"`
    Slug        string    `json:"slug"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`      // HTML or Markdown
    Excerpt     string    `json:"excerpt"`
    Category    string    `json:"category"`
    CoverImage  string    `json:"cover_image"`
    Author      string    `json:"author"`
    PublishedAt time.Time `json:"published_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Rating struct {
    ID        string    `json:"id"`
    ArticleID string    `json:"article_id"`
    Score     int       `json:"score"`          // 1-5
    IPHash    string    `json:"ip_hash"`        // Anti-duplicate
    CreatedAt time.Time `json:"created_at"`
}
```

### Frontend (React + Vite)

- [ ] Inicializar proyecto: `npm create vite@latest frontend -- --template react`
- [ ] Crear estructura del frontend:
  ```
  frontend/
  ├── src/
  │   ├── main.jsx                    # Entry: public SPA
  │   ├── admin-main.jsx              # Entry: admin panel
  │   ├── App.jsx                     # Router público
  │   ├── components/
  │   │   ├── Navbar/
  │   │   ├── Hero/
  │   │   ├── ArticleCard/
  │   │   ├── ArticleRating/          # Rating component (estrellas)
  │   │   ├── SearchFilter/
  │   │   └── Footer/
  │   ├── pages/
  │   │   ├── HomePage/
  │   │   ├── EntradasPage/           # Listado de artículos
  │   │   ├── ArticlePage/            # Artículo individual + rating
  │   │   └── AdminPage/              # Panel admin (CRUD artículos)
  │   ├── hooks/
  │   │   └── useApi.js               # Fetch wrapper
  │   └── styles/
  │       └── global.css              # Migrado de styles.css actual
  ├── index.html
  ├── admin.html
  ├── vite.config.js                  # Build → ../backend/public-dist/
  ├── vite.admin.config.js            # Build → ../backend/admin-dist/
  └── package.json
  ```
- [ ] Migrar estilos actuales (`css/styles.css`) → `frontend/src/styles/global.css`
- [ ] Migrar contenido del `index.html` actual → componentes React
- [ ] Migrar `entradas.html` → `EntradasPage` con filtros y búsqueda
- [ ] Configurar Vite proxy para desarrollo: `/api` → `http://localhost:8080`
- [ ] Preservar SEO: meta tags dinámicos con react-helmet-async
- [ ] Preservar Google Analytics y AdSense

---

## Phase 3: Admin Page — Publicación de Artículos

- [ ] Crear `AdminPage` con login (Bearer token, mismo patrón que comerciantesChavarria)
- [ ] Editor de artículos:
  - Campo: título, slug (auto-generado), categoría, excerpt
  - Editor rich text (usar react-quill o tiptap) para el contenido
  - Upload de imagen de portada (→ Cloudflare R2)
  - Preview antes de publicar
- [ ] Listado de artículos existentes con opciones: editar, eliminar
- [ ] Gestión de mensajes de contacto (ver, eliminar)

---

## Phase 4: Sistema de Puntuación Individual

### Arquitectura
- Cada artículo tiene su propio rating independiente
- Función standalone `rateArticle()` en el backend
- Anti-spam: hash de IP para evitar votos duplicados

### Implementación

- [ ] Backend: handler `POST /api/articles/{slug}/rating`
  ```go
  func rateArticle(slug string, score int, ipHash string) error
  func getArticleRating(slug string) (average float64, count int, err error)
  ```
- [ ] Frontend: componente `ArticleRating`
  - UI de estrellas clickeables (1-5)
  - Muestra promedio + total de votos
  - Feedback visual al votar
  - Almacena en localStorage que ya votó (UX, no seguridad)
- [ ] Mostrar rating en las cards del listado de artículos
- [ ] Endpoint `GET /api/articles/{slug}/rating` para consultar sin auth

---

## Phase 5: Almacenamiento (Free Tier)

### Imágenes/Assets → Cloudflare R2 (Recomendado)

| Aspecto     | Detalle                                               |
|-------------|-------------------------------------------------------|
| Storage     | 10 GB gratis                                          |
| Reads       | 10M Class B ops/mes                                   |
| Writes      | 1M Class A ops/mes                                    |
| Egress      | **$0 — sin costo de descarga** (killer feature)       |
| Duración    | Free tier permanente                                  |
| API         | Compatible con S3                                     |

- [ ] Crear cuenta Cloudflare y bucket R2
- [ ] Configurar credenciales S3-compatible en el backend Go
- [ ] Handler `POST /api/articles/upload-image` → sube a R2, retorna URL pública
- [ ] Configurar custom domain en R2: `media.alejandrogmota.com` (opcional)

### Base de Datos → Opciones Free Tier

| Servicio | Storage | Límites | Duración | Notas |
|----------|---------|---------|----------|-------|
| **Supabase (Postgres)** | 500 MB | 50K MAU auth, 2 GB egress | Permanente | ⚠️ Pausa tras 1 semana inactivo |
| **Oracle Autonomous DB** | 20 GB × 2 instancias | 1 OCPU cada una | Always Free | Ya estás en Oracle, sin pausa |
| **Turso (SQLite edge)** | 5 GB | 500M reads, 10M writes/mes | Permanente | Sin cold starts, ligero |
| **Neon (Serverless Postgres)** | 0.5 GB | 100 compute-hrs/mes | Permanente | Cold start 5 min inactivo |
| **Firebase Firestore** | 1 GB | 50K reads, 20K writes/día | Permanente | NoSQL, quotas diarios |

### Recomendación

**Opción A (más simple):** Supabase Postgres — REST API incluida, auth built-in, 500 MB suficiente para blog. Configurar un cron ping para evitar pausa por inactividad.

**Opción B (sin dependencias externas):** Oracle Autonomous DB — 20 GB, ya estás en el ecosistema Oracle. Más complejo de configurar pero 0 riesgo de pausa.

**Opción C (más ligero):** Turso (libSQL) — 5 GB, sin pausa, edge performance. Driver Go disponible.

- [ ] Elegir opción de base de datos
- [ ] Implementar storage layer en el backend (interface pattern, como comerciantesChavarria)
- [ ] Configurar variables de entorno para conexión

---

## Phase 6: CI/CD — GitHub Actions

### Workflow: Deploy Backend + Frontend

**Archivo:** `.github/workflows/deploy.yml`

**Trigger:** Push a `main` con cambios en `backend/`, `frontend/`, `Makefile`, `.github/workflows/`

```yaml
name: Deploy to Oracle Cloud VM

on:
  push:
    branches: [main]
    paths:
      - 'backend/**'
      - 'frontend/**'
      - 'Makefile'
      - '.github/workflows/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Deploy via SSH
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.VM_HOST }}
          username: ${{ secrets.VM_USER }}
          key: ${{ secrets.VM_SSH_KEY }}
          script: |
            export PATH=$PATH:/usr/local/go/bin:/usr/local/node/bin

            cd ~/Tech-Science-Blog
            git pull origin main

            # Build frontend
            cd frontend
            npm ci
            npx vite build

            # Build admin
            npx vite build --config vite.admin.config.js
            mv ../backend/admin-dist/admin.html ../backend/admin-dist/index.html
            cp public/logo.webp ../backend/admin-dist/ 2>/dev/null || true

            # Build backend (ARM binary)
            cd ../backend
            CGO_ENABLED=0 go build -o tech-blog-api ./cmd/api/

            # Restart service
            sudo systemctl restart tech-blog
```

### GitHub Secrets Requeridos

| Secret       | Valor                          |
|-------------|--------------------------------|
| `VM_HOST`    | IP pública de la VM Oracle     |
| `VM_USER`    | Usuario SSH (ubuntu)           |
| `VM_SSH_KEY` | Llave SSH privada              |
| `ADMIN_USER` | Usuario admin del blog         |
| `ADMIN_PASS` | Password admin                 |
| `GH_PAT`     | GitHub Personal Access Token   |

### Makefile

```makefile
# === Desarrollo Local ===
install:
	cd frontend && npm install

dev-front:
	cd frontend && npm run dev

dev-back:
	cd backend && go run ./cmd/api/

dev:
	make dev-back & make dev-front

build-front:
	cd frontend && npx vite build

build-admin:
	cd frontend && npx vite build --config vite.admin.config.js
	mv backend/admin-dist/admin.html backend/admin-dist/index.html

build-back:
	cd backend && CGO_ENABLED=0 go build -o tech-blog-api ./cmd/api/

build: build-front build-admin build-back

clean:
	rm -rf backend/public-dist backend/admin-dist backend/tech-blog-api

# === VM Setup (ejecutar dentro de la VM) ===
vm-deps:
	sudo apt update && sudo apt install -y git curl ufw nginx certbot python3-certbot-nginx

vm-go:
	wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
	sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz
	rm go1.22.0.linux-arm64.tar.gz

vm-node:
	curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
	sudo apt install -y nodejs

vm-clone:
	cd ~ && git clone https://github.com/AlejandroGMota/Tech-Science-Blog.git

vm-service:
	sudo cp deploy/tech-blog.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable tech-blog
	sudo systemctl restart tech-blog

vm-firewall:
	sudo ufw allow 22/tcp
	sudo ufw allow 80/tcp
	sudo ufw allow 443/tcp
	sudo ufw --force enable

vm-nginx:
	sudo cp deploy/nginx-blog.conf /etc/nginx/sites-available/blog
	sudo ln -sf /etc/nginx/sites-available/blog /etc/nginx/sites-enabled/
	sudo nginx -t && sudo systemctl reload nginx

vm-ssl:
	sudo certbot --nginx -d blog.alejandrogmota.com

vm-setup: vm-deps vm-go vm-node vm-clone vm-firewall vm-nginx vm-service

# === VM Deploy (ejecutar dentro de la VM) ===
vm-deploy:
	cd ~/Tech-Science-Blog && git pull origin main
	make build
	sudo systemctl restart tech-blog

vm-status:
	sudo systemctl status tech-blog

vm-logs:
	sudo journalctl -u tech-blog -f

vm-restart:
	sudo systemctl restart tech-blog
```

---

## Phase 7: Polish & Go Live

- [ ] Responsive testing en mobile/tablet/desktop
- [ ] Lighthouse audit (performance, SEO, accessibility)
- [ ] Migrar los 10 artículos existentes a la base de datos
- [ ] Configurar sitemap.xml dinámico desde el backend
- [ ] Configurar robots.txt
- [ ] Actualizar DNS: A record → IP Oracle VM
- [ ] Configurar SSL (Let's Encrypt + auto-renewal)
- [ ] Verificar Google Analytics y AdSense
- [ ] Eliminar archivos deprecados (HTML artículos, Python scripts)
- [ ] Actualizar README.md

---

## Proposed File Structure (Final)

```
Tech-Science-Blog/
├── backend/
│   ├── cmd/api/
│   │   └── main.go
│   ├── api/
│   │   └── routes.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── models/models.go
│   │   ├── store/
│   │   │   ├── store.go           # Interface
│   │   │   ├── memory.go          # Dev
│   │   │   └── supabase.go        # Prod (o turso.go, oracle.go)
│   │   ├── handlers/
│   │   │   ├── articles.go
│   │   │   ├── ratings.go         # Puntuación individual
│   │   │   ├── contacto.go
│   │   │   └── helpers.go
│   │   └── middleware/
│   │       ├── auth.go
│   │       └── cors.go
│   ├── public-dist/               # Frontend build (gitignored)
│   ├── admin-dist/                # Admin build (gitignored)
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── main.jsx
│   │   ├── admin-main.jsx
│   │   ├── App.jsx
│   │   ├── components/
│   │   │   ├── Navbar/
│   │   │   ├── Hero/
│   │   │   ├── ArticleCard/
│   │   │   ├── ArticleRating/
│   │   │   ├── SearchFilter/
│   │   │   └── Footer/
│   │   ├── pages/
│   │   │   ├── HomePage/
│   │   │   ├── EntradasPage/
│   │   │   ├── ArticlePage/
│   │   │   └── AdminPage/
│   │   ├── hooks/
│   │   └── styles/
│   │       └── global.css
│   ├── index.html
│   ├── admin.html
│   ├── vite.config.js
│   ├── vite.admin.config.js
│   └── package.json
├── deploy/
│   ├── tech-blog.service          # Systemd service file
│   └── nginx-blog.conf            # Nginx config
├── .github/
│   └── workflows/
│       └── deploy.yml             # CI/CD pipeline
├── Makefile
├── PLAN.md
├── README.md
└── .gitignore
```

---

## Tech Stack Summary

| Layer          | Technology                                    |
|---------------|-----------------------------------------------|
| Frontend      | React 19 + Vite                               |
| Backend       | Go (net/http stdlib)                          |
| Database      | Oracle Autonomous JSON DB (20 GB, Always Free)|
| Image Storage | Oracle Object Storage (10 GB, Always Free)    |
| Auth          | Sessions + Bearer tokens                      |
| Hosting       | Oracle Cloud VM ARM A1.Flex (3 OCPU, 18 GB)  |
| Reverse Proxy | Nginx + SSL (Let's Encrypt)                   |
| CI/CD         | GitHub Actions (SSH deploy)                   |
| Analytics     | Google Analytics (GTM)                        |
| Ads           | Google AdSense                                |

---

## Progress Tracker

| Tarea | Estado |
|-------|--------|
| Oracle VCN + Subnet + IGW + Routes + Security List | DONE |
| Oracle Object Storage bucket (tech-blog-images) | DONE |
| Oracle Autonomous JSON DB (tech-blog-db, 20 GB) | DONE |
| Oracle VM ARM A1.Flex (3 OCPU, 18 GB) | ESPERANDO CAPACIDAD (retry script) |
| SSH Keys generadas en Cloud Shell | DONE |
| Backend Go (API, handlers, auth, CORS) | DONE |
| Frontend React (public: Home, Entradas, Article) | DONE |
| Frontend React (admin: CRUD artículos) | DONE |
| ArticleRating component (estrellas 1-5) | DONE |
| Memory store (dev) | DONE |
| Oracle DB store (prod) con auto-migrate | DONE |
| Seed de 10 artículos existentes | DONE |
| Makefile (dev + vm-setup + vm-deploy) | DONE |
| Systemd service file | DONE |
| Nginx reverse proxy config | DONE |
| GitHub Actions CI/CD workflow | DONE |
| .env con credenciales Oracle | DONE |
| Limpieza proyecto (solo backend/ frontend/ deploy/) | DONE |
| Configurar VM (Go, Node, Nginx, systemd) | PENDIENTE (necesita VM) |
| Reservar IP pública | PENDIENTE (necesita VM) |
| DNS blog.alejandrogmota.com → VM IP | PENDIENTE (necesita VM) |
| SSL Let's Encrypt | PENDIENTE (necesita VM) |
| GitHub Secrets (VM_HOST, VM_SSH_KEY, etc.) | PENDIENTE (necesita VM) |
| Conectar backend a Oracle Autonomous DB | PENDIENTE (necesita wallet) |
| Subida de imágenes a Object Storage | PENDIENTE |
| Migrar artículos a Oracle DB (seed en prod) | PENDIENTE |
| Google Analytics + AdSense | PENDIENTE |
| Responsive testing + Lighthouse | PENDIENTE |

---

## Execution Order

1. ~~**Phase 1** — Infraestructura Oracle VM~~ (VCN, DB, bucket DONE — VM esperando capacidad)
2. ~~**Phase 2** — Migración a React + Go~~ (DONE)
3. ~~**Phase 3** — Admin page~~ (DONE)
4. ~~**Phase 4** — Sistema de puntuación~~ (DONE)
5. ~~**Phase 5** — Almacenamiento~~ (Oracle DB store DONE, bucket DONE)
6. ~~**Phase 6** — CI/CD GitHub Actions~~ (DONE)
7. **Phase 7** — VM setup + DNS + SSL + go live (BLOQUEADO: esperando VM)
