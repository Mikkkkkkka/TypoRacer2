package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/data"
)

func QuotesHandlerWithDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Println("Method not allowed for api/v1/quotes")
			return
		}

		string_id := r.URL.Query().Get("id")
		random := r.URL.Query().Get("random")

		switch {
		case string_id != "":
			quoteIdHandler(string_id, db, w)
		case random == "true":
			quoteRandomHandler(db, w)
		default:
			quoteAllHandler(db, w)
		}
	}
}

func quoteIdHandler(string_id string, db *sql.DB, w http.ResponseWriter) {
	id, err := strconv.ParseUint(string_id, 10, 32)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	quote, err := data.GetQuote(uint(id), db)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func quoteRandomHandler(db *sql.DB, w http.ResponseWriter) {
	quote, err := data.GetRandomQuote(db)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func quoteAllHandler(db *sql.DB, w http.ResponseWriter) {
	quotes, err := data.GetAllQuotes(db)
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		log.Println(err)
	}
}
