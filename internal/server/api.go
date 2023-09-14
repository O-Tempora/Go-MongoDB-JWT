package server

import (
	"context"
	"encoding/json"
	"fmt"
	"gomongojwt/internal/middleware"
	"gomongojwt/internal/repository"
	"gomongojwt/internal/util"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

type server struct {
	logger   *slog.Logger
	client   *mongo.Client
	database *mongo.Database
	router   *mux.Router
	store    *repository.Store
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
		AddSource: true,
	}))
	return logger
}
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}, err error) {
	w.WriteHeader(code)
	if err != nil {
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		s.logger.LogAttrs(context.Background(), slog.LevelError, "Response:",
			slog.String("URL", r.URL.Path),
			slog.String("Method", r.Method),
			slog.Int("HTTP Code", code),
			slog.String("HTTP Status", http.StatusText(code)),
			slog.String("Error", err.Error()),
		)
		return
	}

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Response:",
		slog.String("URL", r.URL.Path),
		slog.String("Method", r.Method),
		slog.Int("HTTP Code", code),
		slog.String("HTTP Status", http.StatusText(code)),
	)
}

func (s *server) initRouter() {
	s.router.Use(middleware.LogRequest(s.logger))
	s.router.HandleFunc("/auth", s.handleAuth).Methods("POST")
	s.router.HandleFunc("/refresh", s.handleRefresh).Methods("POST")
}

func (s *server) handleAuth(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	exists, err := s.store.User().IsPresent(guid)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		s.respond(w, r, http.StatusBadRequest, exists, nil)
		return
	}
	access, refresh, err := util.GetTokenPair(guid)
	if err != nil {
		s.respond(w, r, http.StatusUnauthorized, nil, err)
		return
	}
	id, _ := primitive.ObjectIDFromHex(guid)
	err = s.store.User().UpdateRefresh(id, refresh)
	if err != nil {
		s.respond(w, r, http.StatusConflict, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, map[string]string{
		"Access":  access,
		"Refresh": refresh,
	}, nil)
}
func (s *server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Refresh, %s", r.Host)
}
