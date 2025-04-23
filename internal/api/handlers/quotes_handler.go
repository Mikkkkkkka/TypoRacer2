package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/data"
)

type QuotesHandler struct {
	db *sql.DB
}

func NewQuotesHandler(db *sql.DB) *QuotesHandler {
	return &QuotesHandler{db: db}
}

func (handler QuotesHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/quotes", handler)
}

func (handler QuotesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed for api/v1/quotes")
		return
	}

	string_id := r.URL.Query().Get("id")
	random := r.URL.Query().Get("random")

	switch {
	case string_id != "":
		handler.quoteIdHandler(string_id, w)
	case random == "true":
		handler.quoteRandomHandler(w)
	default:
		handler.quoteAllHandler(w)
	}
}

func (handler QuotesHandler) quoteIdHandler(string_id string, w http.ResponseWriter) {
	id, err := strconv.ParseUint(string_id, 10, 32)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println(err)
		return
	}
	quote, err := data.GetQuote(uint(id), handler.db)
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

func (handler QuotesHandler) quoteRandomHandler(w http.ResponseWriter) {
	quote, err := data.GetRandomQuote(handler.db)
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

func (handler QuotesHandler) quoteAllHandler(w http.ResponseWriter) {
	quotes, err := data.GetAllQuotes(handler.db)
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
