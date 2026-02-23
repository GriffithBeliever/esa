# Event Scheduler

A full-stack event scheduling application with AI-powered features.

**Stack:** Go backend (chi, SQLite, JWT) · React/TypeScript frontend (Vite, TanStack Query, Tailwind, Zustand) · Gemini AI

## Quick Start

### Prerequisites
- Go 1.22+
- Node.js 20+

### 1. Configure environment

```bash
cp .env.example .env
# Edit .env and set:
#   JWT_SECRET   — random 64-character string
#   GEMINI_API_KEY — your Google AI Studio key (optional, disables AI features if missing)
```

### 2. Run the backend

```bash
go run ./cmd/server
# Server starts at http://localhost:8080
```

### 3. Run the frontend (dev mode)

```bash
cd frontend
npm install
npm run dev
# Frontend at http://localhost:5173 — proxies /api → :8080
```

### 4. Build for production (single binary)

```bash
cd frontend && npm run build && cd ..
go build -o ims ./cmd/server
./ims   # Serves frontend + API on :8080
```

## Features

- **Auth** — register, login, JWT access tokens + refresh tokens
- **Events** — create, edit, delete, search, filter by date/location/status
- **Calendar view** — monthly grid with color-coded events
- **Invitations** — invite by email, token-based accept/decline links
- **AI: Description Generator** — write event descriptions with Gemini
- **AI: Natural Language Input** — parse "lunch next Tuesday at noon" into a form
- **AI: Conflict Detection + Time Suggestions** — detect scheduling conflicts, suggest alternatives

## API

All endpoints are under `/api/v1/`. See [Implementation Plan](./req.md) for full API reference.

```bash
# Health
curl http://localhost:8080/api/v1/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","username":"you","password":"secret"}'

# Login → get access_token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"secret"}'
```

## Deploy to Railway.app

1. Push to GitHub
2. Create a Railway project, connect the repo
3. Set environment variables in the Railway dashboard:
   - `JWT_SECRET`
   - `GEMINI_API_KEY`
   - `DATABASE_PATH` (Railway volumes, or use the default `./data/events.db`)
   - `FRONTEND_ORIGIN` (your Railway URL, e.g. `https://ims.up.railway.app`)
4. Railway reads `railway.toml` and `nixpacks.toml` to build automatically

## Project Structure

```
ims/
├── cmd/server/         Go entry point + embedded frontend/dist
├── internal/
│   ├── config/         Environment configuration
│   ├── domain/         Domain types and errors
│   ├── repository/     SQLite data access
│   ├── service/        Business logic + AI calls
│   └── handler/        HTTP handlers + router
└── frontend/           Vite + React + TypeScript
    └── src/
        ├── api/        Axios API client
        ├── store/      Zustand auth store
        ├── hooks/      TanStack Query hooks
        ├── auth/       Login/Register pages
        ├── events/     Event CRUD + calendar
        ├── invitations/ Invitation flow
        └── ai/         AI feature components
```
