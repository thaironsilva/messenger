package message

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/thaironsilva/messenger/api/cognitoClient"
	"github.com/thaironsilva/messenger/api/resource/user"
)

var badRequestResponse = []byte(`{"message":"bad request"}`)
var methodNotAllowedResponse = []byte(`{"message":"method not allowed"}`)
var notFoundResponse = []byte(`{"message":"user not found"}`)
var unauthorizedResponse = []byte(`{"message":"unauthorized token"}`)

type Storage interface {
	GetAll(sender_id string, receiver string) ([]Message, error)
	Create(message Message) error
}

type MessageHandler struct {
	storage     Storage
	userStorage user.Storage
	cognito     cognitoClient.CognitoInterface
}

func NewHandler(storage Storage, userStorage user.Storage, cognito cognitoClient.CognitoInterface) MessageHandler {
	return MessageHandler{
		storage:     storage,
		userStorage: userStorage,
		cognito:     cognito,
	}
}

func GetMessages(h MessageHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		cognitoUser, err := h.cognito.GetUserByToken(token)

		if err != nil {
			if err.Error() == "NotAuthorizedException: Could not verify signature for Access Token" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(unauthorizedResponse)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		var email string

		for _, attribute := range cognitoUser.UserAttributes {
			if *attribute.Name == "email" {
				email = *attribute.Value
			}
		}

		sender, err := h.userStorage.GetByEmail(email)

		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				w.Write(notFoundResponse)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		username := strings.TrimPrefix(r.URL.Path, "/api/v0/messages/")

		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		receiver, err := h.userStorage.GetByUsername(username)

		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				w.Write(notFoundResponse)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		messages, err := h.storage.GetAll(sender.Id, receiver.Id)

		if err != nil {
			log.Println("Error listing messages:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		err = json.NewEncoder(w).Encode(messages)

		if err != nil {
			log.Println("Error encoding messages:", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
