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

- [ ] Image upload endpoint (OCI Object Storage bucket exists: tech-blog-images)
- [ ] Google Analytics + AdSense integration
- [ ] Responsive testing + Lighthouse audit
- [ ] sitemap.xml + robots.txt
- [ ] Rich text editor in admin (currently raw HTML)
- [ ] Article pagination (currently returns all)
- [ ] Draft/published article status
- [ ] Rebuild comerciantes-chavarria frontend (stale IP in bundle)

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
