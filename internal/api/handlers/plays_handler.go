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

func PlaysHandlerWithDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			playsGetHandler(w, r, db)
		case http.MethodPost:
			playsPostHandler(w, r, db)
		default:
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
}

func playsGetHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	strUserId := r.URL.Query().Get("user_id")
	if strUserId != "" {
		playsGetByUserIdHandler(strUserId, w, r, db)
		return
	}
	plays, err := data.GetAllPlays(db)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(plays)
}

func playsGetByUserIdHandler(strUserId string, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	plays, err := data.GetPlaysByUserId(userId, db)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(plays)
}

func playsPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var payloadData model.PlayRecord

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	play, err := service.CalculatePlayResults(&payloadData, db)
	if err != nil {
		http.Error(w, "Invalid quote id", http.StatusBadRequest)
		return
	}

	if err := data.AddPlay(play, db); err != nil {
		log.Fatal(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
