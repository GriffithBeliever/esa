# Stage 1: Build the React frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go binary with embedded frontend
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Place the built frontend where Go's embed directive expects it
COPY --from=frontend /app/cmd/server/frontend/dist ./cmd/server/frontend/dist
RUN go build -o ims ./cmd/server

# Stage 3: Minimal runtime image
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/ims ./ims
RUN mkdir -p data
EXPOSE 8080
CMD ["./ims"]
