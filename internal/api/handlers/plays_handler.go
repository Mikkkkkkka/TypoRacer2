package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/data"
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
	var payloadData model.Play

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := data.AddPlay(&payloadData, db); err != nil {
		log.Fatal(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
