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
	router.HandleFunc("GET /api/v0/user", user.GetUser(userHandler))
	router.HandleFunc("GET /api/v0/users", user.GetUsers(userHandler))
	router.HandleFunc("POST /api/v0/users", user.CreateUser(userHandler))
	router.HandleFunc("POST /api/v0/users/confirmation", user.ConfirmAccount(userHandler))
	router.HandleFunc("POST /api/v0/users/login", user.SignIn(userHandler))
	router.HandleFunc("PUT /api/v0/users/password", user.UpdatePassword(userHandler))
	router.HandleFunc("DELETE /api/v0/users", user.DeleteUser(userHandler))

	return router
}
