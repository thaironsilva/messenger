package router

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/thaironsilva/messenger/api/resource/user"
)

func New(l *log.Logger, db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	userHandler := user.New(l, db)
	router.HandleFunc("GET /users", userHandler.List)
	router.HandleFunc("POST /users", userHandler.Create)

	return router
}
