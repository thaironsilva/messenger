package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/thaironsilva/messenger/cognitoClient"
)

var badRequestResponse = []byte(`{"message":"bad request"}`)
var methodNotAllowedResponse = []byte(`{"message":"method not allowed"}`)
var notFoundResponse = []byte(`{"message":"user not found"}`)

type Storage interface {
	GetById(id string) (User, error)
	GetAll() ([]User, error)
	Create(user User) error
	Update(user User) error
	Delete(is string) error
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

		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message":"user id should be provided as the value of an 'id' querystring parameter"}`))
			return
		}

		user, err := h.storage.GetById(id)

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

		err = json.NewEncoder(w).Encode(user)

		if err != nil {
			log.Println("Error encoding users:", err)
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

func UpdateUser(h UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message":"user id should be provided as the value of an 'id' querystring parameter"}`))
			return
		}

		user, err := h.storage.GetById(id)

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

		if r.Body == nil {
			log.Println("update requires a request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		var updatedUser User

		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			log.Println("Error decoding user:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		user.Username = updatedUser.Username
		user.Email = updatedUser.Email

		if err := h.storage.Update(user); err != nil {
			log.Println("Error occurred while trying to update user:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(h UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}

		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message":"user id should be provided as the value of an 'id' querystring parameter"}`))
			return
		}

		_, err := h.storage.GetById(id)

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

		if err := h.storage.Delete(id); err != nil {
			log.Println("Error occurred while trying to update user:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusOK)
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
			log.Println("create requires a request body")
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cognitoUser)
	}
}
