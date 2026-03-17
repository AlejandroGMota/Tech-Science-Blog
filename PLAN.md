# Tech-Science-Blog — Project Status

## Current State: LIVE

**URL:** https://blog.alejandrogmota.com
**VM:** Oracle Cloud ARM (shared with comerciantes-chavarria)
**IP:** 160.34.222.215 (ephemeral)

---

## Completed

- [x] Oracle Cloud VM (shared chavarria-api, 1 OCPU / 6GB ARM)
- [x] Go backend API (net/http stdlib)
- [x] React 19 + Vite frontend (public + admin SPAs)
- [x] Oracle Autonomous JSON DB connected (godror driver)
- [x] 10 articles seeded in Oracle DB
- [x] Article rating system (1-5 stars, one per IP)
- [x] Contact form (public create, admin list/delete)
- [x] Admin panel with Bearer token auth
- [x] Nginx reverse proxy + SSL (Let's Encrypt)
- [x] GitHub Actions CI/CD (triggers on push to main)
- [x] Systemd service (auto-restart)
- [x] DNS: blog.alejandrogmota.com -> A record -> VM IP
- [x] Oracle Instant Client 23.7 (aarch64)
- [x] Wallet configured for Autonomous DB
- [x] Favicon + logo branding
- [x] README rewritten (UTF-8)

## Pending

### 1. Rich Text Editor (admin panel) — DONE
- [x] Install TipTap v2 editor (lightweight, extensible, no jQuery)
- [x] Replace raw HTML textarea with TipTap in AdminPage
- [x] Toolbar: bold, italic, underline, strike, headings (h2, h3), links, lists, blockquote, code block
- [x] Image insert (paste URL)
- [x] HTML preview toggle (see/edit raw HTML output)
- [x] Content outputs clean HTML for storage

### 2. Image Upload (OCI Object Storage)
- [ ] Backend: POST /api/articles/upload-image handler
- [ ] OCI SDK integration (pre-authenticated requests or S3-compatible API)
- [ ] Accept multipart/form-data, validate file type/size (max 5MB, jpg/png/webp)
- [ ] Return public URL after upload
- [ ] Frontend: drag-and-drop or button upload in editor
- [ ] Bucket: tech-blog-images (namespace: axva0xxfvkwr, region: mx-queretaro-1)

### 3. Article Pagination
- [ ] Backend: add ?page=1&limit=10 query params to GET /api/articles
- [ ] Return total count in response header or body
- [ ] Frontend: pagination controls in EntradasPage
- [ ] Admin: paginate article list if > 20

### 4. Draft/Published Status
- [ ] Add `status` field to Article model (draft/published)
- [ ] Oracle DB: ALTER TABLE articles ADD status VARCHAR2(20) DEFAULT 'published'
- [ ] Backend: filter by status on public GET (only published)
- [ ] Admin: show all articles, toggle draft/published
- [ ] Save as draft button in editor

### 5. SEO
- [ ] Backend: GET /sitemap.xml — dynamic from published articles
- [ ] Backend: GET /robots.txt — allow all, point to sitemap
- [ ] Meta tags per article (Open Graph, Twitter Cards)

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

## VM Details

- Host: chavarria-api (shared)
- Port 8080: comerciantes-chavarria
- Port 8081: tech-blog
- Path: /home/ubuntu/tech-blog
- Service: tech-blog.service
- Env: /home/ubuntu/tech-blog/.env
- Wallet: /home/ubuntu/tech-blog/wallet/
- Instant Client: /opt/oracle/instantclient_23_7/

## CI/CD Flow

1. Push to main (backend/** or frontend/**)
2. GitHub Actions SSH into VM
3. git pull origin main
4. npm ci + vite build (public + admin)
5. go build
6. systemctl restart tech-blog

## GitHub Secrets

| Secret | Description |
|--------|-------------|
| VM_HOST | 160.34.222.215 |
| VM_USER | ubuntu |
| VM_SSH_KEY | chavarria-vm private key |
