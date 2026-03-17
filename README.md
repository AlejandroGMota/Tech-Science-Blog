# Tech & Science Blog

A full-stack blog platform for technology and science articles, built with React and Go.

**Live:** [blog.alejandrogmota.com](https://blog.alejandrogmota.com)

---

## Tech Stack

- **Frontend:** React 19 + Vite
- **Backend:** Go (net/http stdlib)
- **Database:** Oracle Autonomous JSON DB
- **Infrastructure:** Oracle Cloud Free Tier (ARM VM)
- **CI/CD:** GitHub Actions -> SSH deploy

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

### Frontend
```bash
cd frontend
npm install
npm run dev
```

---

## Deployment

Push to `main` branch triggers automatic deployment via GitHub Actions.

## License

MIT
