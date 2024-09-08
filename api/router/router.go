package router

import (
	"database/sql"
	"net/http"

	"github.com/thaironsilva/messenger/api/resource/user"
)

func New(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	userStorage := user.NewRepository(db)
	router.HandleFunc("GET /users", user.GetUsers(userStorage))
	router.HandleFunc("POST /users", user.CreateUsers(userStorage))

	return router
}
