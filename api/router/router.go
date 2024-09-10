package router

import (
	"database/sql"
	"net/http"

	"github.com/thaironsilva/messenger/api/resource/user"
)

func New(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	userStorage := user.NewRepository(db)
	router.HandleFunc("GET /users/{id}", user.GetUser(userStorage))
	router.HandleFunc("GET /users", user.GetUsers(userStorage))
	router.HandleFunc("POST /users", user.CreateUser(userStorage))
	router.HandleFunc("PUT /users/{id}", user.UpdateUser(userStorage))
	router.HandleFunc("DELETE /users/{id}", user.DeleteUser(userStorage))

	return router
}
