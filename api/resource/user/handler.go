package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var badRequestResponse = []byte(`{"message":"bad request"}`)
var methodNotAllowedResponse = []byte(`{"message":"method not allowed"}`)

type Storage interface {
	GetAll() ([]User, error)
	Create(user User) error
}

func GetUsers(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(methodNotAllowedResponse)
			return
		}
		users, err := storage.GetAll()
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

func CreateUsers(storage Storage) http.HandlerFunc {
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

		var newUser User

		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			log.Println("Error decoding user:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		if err := newUser.IsValid(); err != nil {
			log.Println("Error validating user:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		now := time.Now()
		newUser.CreatedAt = now
		newUser.UpdatedAt = now

		if err := storage.Create(newUser); err != nil {
			log.Println("Error occurred while trying to create user:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}
