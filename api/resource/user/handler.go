package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/thaironsilva/messenger/cognitoClient"
)

var badRequestResponse = []byte(`{"message":"bad request"}`)
var methodNotAllowedResponse = []byte(`{"message":"method not allowed"}`)
var notFoundResponse = []byte(`{"message":"user not found"}`)
var unauthorizedResponse = []byte(`{"message":"unauthorized token"}`)

type Storage interface {
	GetById(id string) (User, error)
	GetAll() ([]User, error)
	Create(user User) error
	Update(user User) error
}

type UserHandler struct {
	storage Storage
	cognito cognitoClient.CognitoInterface
}

func NewHandler(storage Storage, cognito cognitoClient.CognitoInterface) UserHandler {
	return UserHandler{
		storage: storage,
		cognito: cognito,
	}
}

func GetUser(h UserHandler) http.HandlerFunc {
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

		user := &cognitoClient.UserResponse{}

		for _, attribute := range cognitoUser.UserAttributes {
			switch *attribute.Name {
			case "sub":
				user.ID = *attribute.Value
			case "nickname":
				user.Username = *attribute.Value
			case "email":
				user.Email = *attribute.Value
			case "email_verified":
				emailVerified, err := strconv.ParseBool(*attribute.Value)
				if err == nil {
					user.EmailVerified = emailVerified
				}
			}
		}

		err = json.NewEncoder(w).Encode(user)

		if err != nil {
			log.Println("Error encoding user:", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetUsers(h UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		users, err := h.storage.GetAll()

		if err != nil {
			log.Println("Error listing users:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		err = json.NewEncoder(w).Encode(users)

		if err != nil {
			log.Println("Error encoding users:", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func CreateUser(h UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		if r.Body == nil {
			log.Println("create requires a request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		var cognitoUser cognitoClient.CognitoUser

		if err := json.NewDecoder(r.Body).Decode(&cognitoUser); err != nil {
			log.Println("Error decoding user:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		err := h.cognito.SignUp(&cognitoUser)

		if err != nil {
			log.Println("Error occurred while trying to sign up:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		newUser := User{
			Username: cognitoUser.NickName,
			Email:    cognitoUser.Email,
		}

		if err := h.storage.Create(newUser); err != nil {
			log.Println("Error occurred while trying to create user:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}

func UpdatePassword(h UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if r.Body == nil {
			log.Println("update password requires a request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		var user cognitoClient.UserLogin

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		for _, attribute := range cognitoUser.UserAttributes {
			if *attribute.Name == "email" && *attribute.Value != user.Email {
				log.Println("Error updating password.")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(badRequestResponse)
				return
			}
		}

		if err := h.cognito.UpdatePassword(&user); err != nil {
			log.Println("Error occurred while trying to update user password:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func ConfirmAccount(h UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		if r.Body == nil {
			log.Println("confirm requires a request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		var cognitoUser cognitoClient.UserConfirmation

		if err := json.NewDecoder(r.Body).Decode(&cognitoUser); err != nil {
			log.Println("Error decoding user:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		err := h.cognito.ConfirmAccount(&cognitoUser)

		if err != nil {
			log.Println("Error occurred while trying to confirm account:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cognitoUser)
	}
}

func SignIn(h UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		var user cognitoClient.UserLogin

		if r.Body == nil {
			log.Println("login requires a request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		token, err := h.cognito.SignIn(&user)

		if err != nil {
			log.Println("Error occurred while trying to signin:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(token)
	}
}
