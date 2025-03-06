package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/abhilash111/ecom/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	sql  *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr,
		db,
	}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler := user.NewUserHandler()
	userHandler.RegisterUserRoutes(subrouter)

	log.Println("Starting server on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
