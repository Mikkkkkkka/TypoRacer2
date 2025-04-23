package handlers

import (
	"database/sql"
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
	db *sql.DB
}

func NewPlaysHandler(db *sql.DB) *PlaysHandler {
	return &PlaysHandler{db: db}
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
	plays, err := data.GetAllPlays(handler.db)
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
	plays, err := data.GetPlaysByUserId(uint(userId), handler.db)
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
	user, err := data.GetUserFromToken(token[BEARER_PREFIX_LENGTH:], handler.db)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}

	play, err := service.CalculatePlayResults(user.Id, &payloadData, handler.db)
	if err != nil {
		http.Error(w, "Invalid quote id", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := data.AddPlay(play, handler.db); err != nil {
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
