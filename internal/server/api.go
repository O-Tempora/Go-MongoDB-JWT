package server

import (
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
		client:   nil, //temp
		database: nil, //temp
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
func (s *server) initRouter() {

}
