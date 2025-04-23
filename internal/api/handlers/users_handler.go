package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/internal/service"
)

type UsersHandler struct {
	db *sql.DB
}

func NewUsersHandler(db *sql.DB) *UsersHandler {
	return &UsersHandler{db: db}
}

func (handler UsersHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/users/{id}", handler)
}

func (handler UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id format", http.StatusBadRequest)
		log.Println(err)
		return
	}

	user, err := data.GetUserWithoutTokenById(userId, handler.db)
	if err != nil {
		http.Error(w, "User with id does not exist", http.StatusBadRequest)
		log.Println(err)
		return
	}

	stats, err := service.CalculateStats(user, handler.db)
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
