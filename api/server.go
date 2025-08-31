package api

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/bolusarz/task-manager/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store    db.Store
	validate *validator.Validate
	router   *chi.Mux
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store:    store,
		validate: validator.New(),
	}

	server.validate.RegisterValidation("strong", IsPasswordStrong)

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Post("/api/v1/users", s.Register)
	router.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am good")
	})

	s.router = router
}

func (s *Server) StartServer() {
	http.ListenAndServe(":3000", s.router)
}
