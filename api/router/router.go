package router

import (
	"database/sql"
	"net/http"

	"github.com/thaironsilva/messenger/api/cognitoClient"
	"github.com/thaironsilva/messenger/api/connectionManager"
	"github.com/thaironsilva/messenger/api/resource/message"
	"github.com/thaironsilva/messenger/api/resource/user"
)

func New(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	cognito := cognitoClient.NewCognitoClient()
	userRepository := user.NewRepository(db)
	messageRepository := message.NewRepository(db)

	connHandler := connectionManager.NewConnectionHandler(userRepository, cognito)
	router.HandleFunc("/messages/{username}", connHandler.HandleConnections)

	messageHandler := message.NewHandler(messageRepository, userRepository, cognito)
	router.HandleFunc("GET /api/v0/messages/{username}", message.GetMessages(messageHandler))

	userHandler := user.NewHandler(userRepository, cognito)
	router.HandleFunc("GET /api/v0/user", user.GetUser(userHandler))
	router.HandleFunc("GET /api/v0/users", user.GetUsers(userHandler))
	router.HandleFunc("POST /api/v0/users", user.CreateUser(userHandler))
	router.HandleFunc("POST /api/v0/users/confirmation", user.ConfirmAccount(userHandler))
	router.HandleFunc("POST /api/v0/users/login", user.SignIn(userHandler))
	router.HandleFunc("PUT /api/v0/users/password", user.UpdatePassword(userHandler))
	router.HandleFunc("DELETE /api/v0/users", user.DeleteUser(userHandler))

	return router
}
