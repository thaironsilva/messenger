package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type UserHandler struct {
	logger     *log.Logger
	repository *Repository
}

func New(logger *log.Logger, db *sql.DB) *UserHandler {
	return &UserHandler{
		logger:     logger,
		repository: NewRepository(db),
	}
}

func (h *UserHandler) List(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	users := h.repository.list()
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newUser User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.logger.Println("Error decoding user: ", err)
		return
	}

	h.repository.create(newUser)
	w.WriteHeader(http.StatusCreated)
}
