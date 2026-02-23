# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Full-stack application with a Go backend API and React frontend.

## Commands

### Backend (Go)
```bash
go run ./cmd/server          # Start the dev server
go build ./...               # Build all packages
go test ./...                # Run all tests
go test ./path/to/pkg/...    # Run tests in a specific package
go test -run TestFunctionName ./...  # Run a single test
go vet ./...                 # Static analysis
```

### Frontend (React)
```bash
npm install                  # Install dependencies
npm run dev                  # Start dev server (Vite) or: npm start (CRA)
npm run build                # Production build
npm test                     # Run tests
npm test -- --testNamePattern="test name"  # Run a single test
npm run lint                 # Run ESLint
```

## Architecture

### Backend (Go)
- Entry point: `cmd/server/main.go`
- Handlers/controllers live in `internal/handler/` or `internal/api/`
- Business logic in `internal/service/`
- Data access in `internal/repository/` or `internal/store/`
- Domain types/models in `internal/domain/` or `internal/model/`
- Configuration loaded from environment variables (use a `.env` file locally)

Follow standard Go project layout: https://github.com/golang-standards/project-layout

### Frontend (React)
- Source in `frontend/src/` or `web/src/`
- Feature-based folder structure: group files by feature/domain, not by type
- API calls centralized in a `api/` directory using `fetch` or `axios`
- State management: prefer React context + hooks for simple state; use Zustand or Redux Toolkit for complex global state

## Conventions

### Go
- Return errors explicitly; never ignore them
- Use `context.Context` as the first argument in functions that do I/O
- Interfaces defined where they are used (consumer side), not where implemented
- Table-driven tests with `t.Run` subtests
- Gracefully shutdown enabled using context

### React
- Functional components only; no class components
- Co-locate component styles, tests, and types with the component file
- Custom hooks (`use*.ts`) for reusable stateful logic
- Avoid `any` in TypeScript — type all API responses

## API Contract
- Backend serves gRPC
- Prefix all API routes with `/api/v1/`
- Use standard HTTP status codes; errors return `{"error": "message"}`
- CORS is configured on the Go server to allow the React dev server origin
