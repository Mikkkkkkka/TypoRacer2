package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/service"
)

type UsersHandler struct {
	service *service.UserService
}

func NewUsersHandler(service *service.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (handler UsersHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/users/{id}", handler)
}

func (handler UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user id format", http.StatusBadRequest)
		log.Println(err)
		return
	}
	stats, err := handler.service.CalculateStats(uint(userId))
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
