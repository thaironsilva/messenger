package router

import (
	"database/sql"
	"net/http"

	"github.com/thaironsilva/messenger/api/resource/user"
	"github.com/thaironsilva/messenger/cognitoClient"
)

func New(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	repository := user.NewRepository(db)
	cognito := cognitoClient.NewCognitoClient()
	userHandler := user.NewHandler(repository, cognito)
	router.HandleFunc("GET /users/{id}", user.GetUser(userHandler))
	router.HandleFunc("GET /users", user.GetUsers(userHandler))
	router.HandleFunc("POST /users", user.CreateUser(userHandler))
	router.HandleFunc("POST /users/confirmation", user.ConfirmAccount(userHandler))
	router.HandleFunc("PUT /users/{id}", user.UpdateUser(userHandler))
	router.HandleFunc("DELETE /users/{id}", user.DeleteUser(userHandler))

	return router
}
