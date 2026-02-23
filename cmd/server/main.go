package main

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/solomon/ims/internal/config"
	"github.com/solomon/ims/internal/handler"
	"github.com/solomon/ims/internal/repository"
	"github.com/solomon/ims/internal/service"
)

// Frontend is built to cmd/server/frontend/dist (see frontend/vite.config.ts outDir setting).
// This allows the embed to work since Go embed paths must be relative to the source file.
//
//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	loadDotEnv(".env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	db, err := repository.NewDB(ctx, cfg.DatabasePath)
	cancel()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	invRepo := repository.NewInvitationRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)
	eventSvc := service.NewEventService(eventRepo)
	invSvc := service.NewInvitationService(invRepo, eventRepo, userRepo)
	aiSvc := service.NewAIService(cfg.GeminiAPIKey, cfg.GeminiModel, cfg.GeminiTimeout)

	// Router
	apiRouter := handler.NewRouter(authSvc, eventSvc, invSvc, aiSvc, userRepo, cfg.FrontendOrigin)

	mux := http.NewServeMux()
	mux.Handle("/api/", apiRouter)

	dist, subErr := fs.Sub(frontendFS, "frontend/dist")
	if subErr != nil {
		// No frontend built yet — just serve the API
		log.Printf("warning: frontend not embedded (run: cd frontend && npm run build): %v", subErr)
		mux.Handle("/", apiRouter)
	} else {
		mux.Handle("/", spaHandler(dist))
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server starting on :%s (env=%s)", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}

// spaHandler serves static files and falls back to index.html for SPA routing.
func spaHandler(distFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(distFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		_, err := distFS.Open(path)
		if err != nil {
			// Serve index.html for client-side routing
			r2 := *r
			r2.URL.Path = "/"
			fileServer.ServeHTTP(w, &r2)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

// loadDotEnv reads key=value pairs from the given file and sets them as env vars
// if not already set. Silently ignores missing files.
func loadDotEnv(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove surrounding quotes
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') ||
			(val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
