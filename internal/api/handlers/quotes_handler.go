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
			w.WriteHeader(http.StatusMethodNotAllowed)
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
	id, err := strconv.Atoi(string_id)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	quote, err := data.GetQuote(id, db)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
	w.WriteHeader(http.StatusOK)
}

func quoteRandomHandler(db *sql.DB, w http.ResponseWriter) {
	quote, err := data.GetRandomQuote(db)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
	w.WriteHeader(http.StatusOK)
}

func quoteAllHandler(db *sql.DB, w http.ResponseWriter) {
	quotes, err := data.GetAllQuotes(db)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
	w.WriteHeader(http.StatusOK)
}
