package server

import (
	"fmt"
	"gomongojwt/internal/middleware"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

type server struct {
	logger   *slog.Logger
	client   *mongo.Client
	database *mongo.Database
	router   *mux.Router
}

func initServer() *server {
	s := &server{
		client:   nil,
		database: nil,
		logger:   initLogger(os.Stdout),
		router:   mux.NewRouter(),
	}
	s.initRouter()
	return s
}
func initLogger(wr io.Writer) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(wr, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}))
	return logger
}
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) initRouter() {
	s.router.Use(middleware.LogRequest(s.logger))
	s.router.HandleFunc("/auth", s.handleAuth).Methods("POST")
	s.router.HandleFunc("/refresh", s.handleRefresh).Methods("POST")
}

func (s *server) handleAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth, %s", r.Host)
}
func (s *server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Refresh, %s", r.Host)
}
