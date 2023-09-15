package server

import (
	"context"
	"encoding/json"
	"gomongojwt/internal/middleware"
	"gomongojwt/internal/repository"
	"gomongojwt/internal/util"
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
		AddSource: false,
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
	if err != nil || !exists {
		s.respond(w, r, http.StatusBadRequest, "no such user", err)
		return
	}
	access, refresh, err := util.GetTokenPair(guid)
	if err != nil {
		s.respond(w, r, http.StatusUnauthorized, nil, err)
		return
	}
	err = s.store.User().UpdateRefresh(guid, refresh)
	if err != nil {
		s.respond(w, r, http.StatusUnauthorized, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, map[string]string{
		"access":  access,
		"refresh": refresh,
	}, nil)
}
func (s *server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	guid, err := util.ValidateJWT(body.Access)
	if err != nil {
		s.respond(w, r, http.StatusUnauthorized, nil, err)
		return
	}
	same, err := s.store.User().CompareRefreshAndHash(body.Refresh, guid.User)
	if err != nil || !same {
		s.respond(w, r, http.StatusForbidden, nil, err)
		return
	}
	body.Access, body.Refresh, err = util.GetTokenPair(guid.User)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	if err := s.store.User().UpdateRefresh(guid.User, body.Refresh); err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, map[string]string{
		"access":  body.Access,
		"refresh": body.Refresh,
	}, nil)
}
