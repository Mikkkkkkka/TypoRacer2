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

func PlaysHandlerWithDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			playsGetHandler(w, r, db)
		case http.MethodPost:
			playsPostHandler(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Println("Method not allowed for api/v1/plays")
			return
		}
	}
}

func playsGetHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	strUserId := r.URL.Query().Get("user_id")
	if strUserId != "" {
		playsGetByUserIdHandler(strUserId, w, db)
		return
	}
	plays, err := data.GetAllPlays(db)
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

func playsGetByUserIdHandler(strUserId string, w http.ResponseWriter, db *sql.DB) {
	userId, err := strconv.ParseUint(strUserId, 10, 32)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	plays, err := data.GetPlaysByUserId(uint(userId), db)
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

func playsPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var payloadData model.PlayRecord

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer r.Body.Close()
	token := r.Header.Get("Authorization")
	user, err := data.GetUserFromToken(token[BEARER_PREFIX_LENGTH:], db)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}

	play, err := service.CalculatePlayResults(user.Id, &payloadData, db)
	if err != nil {
		http.Error(w, "Invalid quote id", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err := data.AddPlay(play, db); err != nil {
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
