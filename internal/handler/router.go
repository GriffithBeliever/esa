package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/solomon/ims/internal/repository"
	"github.com/solomon/ims/internal/service"
)

func NewRouter(
	authSvc service.AuthService,
	eventSvc service.EventService,
	invSvc service.InvitationService,
	aiSvc service.AIService,
	userRepo repository.UserRepository,
	frontendOrigin string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORSMiddleware(frontendOrigin))

	authH := NewAuthHandler(authSvc, userRepo)
	eventH := NewEventHandler(eventSvc)
	invH := NewInvitationHandler(invSvc, userRepo)
	aiH := NewAIHandler(aiSvc, eventSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/refresh", authH.Refresh)

			r.Group(func(r chi.Router) {
				r.Use(AuthMiddleware(authSvc))
				r.Post("/logout", authH.Logout)
				r.Get("/me", authH.Me)
			})
		})

		// Event routes (all authenticated)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authSvc))

			r.Route("/events", func(r chi.Router) {
				r.Get("/", eventH.List)
				r.Post("/", eventH.Create)
				r.Get("/conflicts", eventH.CheckConflicts)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", eventH.Get)
					r.Put("/", eventH.Update)
					r.Delete("/", eventH.Delete)
					r.Patch("/status", eventH.UpdateStatus)

					// Invitation sub-routes
					r.Post("/invitations", invH.Send)
					r.Get("/invitations", invH.ListByEvent)
					r.Delete("/invitations/{invID}", invH.Delete)
				})
			})

			// Incoming invitations
			r.Get("/invitations/incoming", invH.ListIncoming)

			// AI routes
			r.Route("/ai", func(r chi.Router) {
				r.Post("/generate-description", aiH.GenerateDescription)
				r.Post("/parse-event", aiH.ParseEvent)
				r.Post("/suggest-times", aiH.SuggestTimes)
			})
		})

		// Invitation respond: token-based, auth optional (picks up userID if logged in)
		r.With(optionalAuth(authSvc)).Post("/invitations/{token}/respond", invH.Respond)
	})

	return r
}

// optionalAuth sets userID in context if a valid Bearer token is present, but doesn't reject requests without one.
func optionalAuth(authSvc service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if len(header) > 7 && header[:7] == "Bearer " {
				if userID, err := authSvc.ValidateAccessToken(header[7:]); err == nil {
					ctx := context.WithValue(r.Context(), userIDKey, userID)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
