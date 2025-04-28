package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/service"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

const BEARER_PREFIX_LENGTH = 8

type PlaysHandler struct {
	playService *service.PlayService
	userService *service.UserService
}

func NewPlaysHandler(playService *service.PlayService, userService *service.UserService) *PlaysHandler {
	return &PlaysHandler{
		playService: playService,
		userService: userService,
	}
}

func (handler PlaysHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/plays", handler.playsGetHandler)
	mux.HandleFunc("POST /api/v1/plays", handler.playsPostHandler)
}

func (handler PlaysHandler) playsGetHandler(w http.ResponseWriter, r *http.Request) {
	strUserId := r.URL.Query().Get("userId")
	if strUserId != "" {
		handler.playsGetByUserIdHandler(strUserId, w)
		return
	}
	plays, err := handler.playService.GetAllPlays()
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	if err := json.NewEncoder(w).Encode(plays); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (handler PlaysHandler) playsGetByUserIdHandler(strUserId string, w http.ResponseWriter) {
	userId, err := strconv.ParseUint(strUserId, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user id format", http.StatusBadRequest)
		log.Println(err)
		return
	}
	plays, err := handler.playService.GetPlaysByUserId(uint(userId))
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	if err := json.NewEncoder(w).Encode(plays); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (handler PlaysHandler) playsPostHandler(w http.ResponseWriter, r *http.Request) {
	var payloadData model.PlayRecord

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer r.Body.Close()
	token := r.Header.Get("Authorization")
	user, err := handler.userService.AuthorizeUser(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println(err)
		return
	}
	play, err := handler.playService.RegisterPlayRecord(user, &payloadData)
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err := json.NewEncoder(w).Encode(play); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
