package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/internal/service"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

const BEARER_PREFIX_LENGTH = 8

type PlaysHandler struct {
	service *service.PlayService
}

func NewPlaysHandler(service *service.PlayService) *PlaysHandler {
	return &PlaysHandler{service: service}
}

func (handler PlaysHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/plays", handler)
}

func (handler PlaysHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.playsGetHandler(w, r)
	case http.MethodPost:
		handler.playsPostHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed for api/v1/plays")
		return
	}
}

func (handler PlaysHandler) playsGetHandler(w http.ResponseWriter, r *http.Request) {
	strUserId := r.URL.Query().Get("user_id")
	if strUserId != "" {
		handler.playsGetByUserIdHandler(strUserId, w)
		return
	}
	plays, err := handler.service.GetAllPlays()
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
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	plays, err := handler.service.GetPlaysByUserId(uint(userId))
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
	user, err := data.GetUserFromToken(token[BEARER_PREFIX_LENGTH:], handler.service.Db)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}

	play, err := handler.service.CalculatePlayResults(user.Id, &payloadData)
	if err != nil {
		http.Error(w, "Invalid quote id", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := handler.service.AddPlay(play); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(play); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
