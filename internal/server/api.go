package server

import (
	"context"
	"encoding/json"
	"fmt"
	"gomongojwt/internal/middleware"
	"gomongojwt/internal/service"
	"io"
	"net/http"
	"os"

	_ "gomongojwt/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

type server struct {
	logger  *slog.Logger
	client  *mongo.Client
	router  *mux.Router
	service service.Service
	config  *Config
}

func initServer(config *Config) *server {
	s := &server{
		client:  nil,
		service: nil,
		logger:  initLogger(os.Stdout),
		router:  mux.NewRouter(),
		config:  config,
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
	s.router.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost%s/swagger/doc.json", s.config.Port)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	s.router.Use(middleware.LogRequest(s.logger))
	s.router.HandleFunc("/auth", s.handleAuth).Methods("POST")
	s.router.HandleFunc("/refresh", s.handleRefresh).Methods("POST")
}

// AuthorizeUser godoc
// @Summary      Performs user authorization via tokens
// @Description  Get Access and Refresh tokens by GUID
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param		 guid	query	string true "User's GUID"
// @Router       /auth [post]
// @Success 200 {object} TokenPair
// @Failure 401 {string}	error
func (s *server) handleAuth(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	access, refresh, err := s.service.AuthorizeUser(guid)
	if err != nil {
		s.respond(w, r, http.StatusUnauthorized, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil)
}

// RefreshTokens godoc
// @Summary      Refreshes Access and Refresh tokens
// @Description  Refresh tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param		 tokenPair	body	TokenPair	true	"Access and Refresh tokens"
// @Router       /refresh [post]
// @Success 200 {object} TokenPair
// @Failure 400 {string}	error
// @Failure 401 {string}	error
func (s *server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	body := &TokenPair{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	newAccess, newRefresh, err := s.service.RefreshTokens(body.Access, body.Refresh)
	if err != nil {
		s.respond(w, r, http.StatusUnauthorized, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, TokenPair{
		Access:  newAccess,
		Refresh: newRefresh,
	}, nil)
}
