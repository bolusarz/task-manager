package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/bolusarz/task-manager/db/sqlc"
	"github.com/bolusarz/task-manager/token"
	"github.com/bolusarz/task-manager/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	validate   *validator.Validate
	router     *chi.Mux
	config     util.Config
	tokenMaker token.TokenMaker
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		validate:   validator.New(),
		config:     config,
		tokenMaker: tokenMaker,
	}

	server.validate.RegisterValidation("strong", IsPasswordStrong)

	server.setupRoutes()

	return server, nil
}

func (s *Server) setupRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Post("/api/v1/users", s.Register)
	router.Post("/api/v1/login", s.Login)

	router.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am good")
	})

	s.router = router
}

func (s *Server) StartServer(addr string) {
	fmt.Println("Server listening at ", addr)
	err := http.ListenAndServe(addr, s.router)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to start server: %v", err))
	}
}
