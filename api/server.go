package api

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	db "github.com/arcedo/financial-ai-backend/database"
)

type Server struct {
	listenAddress string
	store         db.MongoStorage
	router        *http.ServeMux
}

func NewServer(listenAddress string, store db.MongoStorage) *Server {
	return &Server{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *Server) Start() error {
	s.setupRoutes()

	server := &http.Server{
		Addr:    s.listenAddress,
		Handler: Cors(s.router),
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Printf("Shutting down server...")
		server.Close()
	}()

	log.Printf("Server running on port %s", s.listenAddress)
	return server.ListenAndServe()
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: change in production to our domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
