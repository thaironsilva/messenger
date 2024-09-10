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
var notFoundResponse = []byte(`{"message":"user not found"}`)

type Storage interface {
	GetById(id string) (User, error)
	GetAll() ([]User, error)
	Create(user User) error
	Update(user User) error
	Delete(is string) error
}

func GetUser(storage Storage) http.HandlerFunc {
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

		user, err := storage.GetById(id)

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

func CreateUser(storage Storage) http.HandlerFunc {
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
			log.Println("Error validating parameters:", err)
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

func UpdateUser(storage Storage) http.HandlerFunc {
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

		user, err := storage.GetById(id)

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

		if err := updatedUser.IsValid(); err != nil {
			log.Println("Error validating parameters:", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(badRequestResponse)
			return
		}

		user.Username = updatedUser.Username
		user.Email = updatedUser.Email
		user.Password = updatedUser.Password
		user.UpdatedAt = time.Now()

		if err := storage.Update(user); err != nil {
			log.Println("Error occurred while trying to update user:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(storage Storage) http.HandlerFunc {
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

		_, err := storage.GetById(id)

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

		if err := storage.Delete(id); err != nil {
			log.Println("Error occurred while trying to update user:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": %s}`, err)))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
