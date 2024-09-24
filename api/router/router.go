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
	router.HandleFunc("GET /user", user.GetUser(userHandler))
	router.HandleFunc("GET /users", user.GetUsers(userHandler))
	router.HandleFunc("POST /users", user.CreateUser(userHandler))
	router.HandleFunc("POST /users/confirmation", user.ConfirmAccount(userHandler))
	router.HandleFunc("POST /users/login", user.SignIn(userHandler))
	router.HandleFunc("PUT /users/password", user.UpdatePassword(userHandler))

	return router
}
