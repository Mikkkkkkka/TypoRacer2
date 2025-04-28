package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/service"
)

type QuotesHandler struct {
	service *service.QuoteService
}

func NewQuotesHandler(service *service.QuoteService) *QuotesHandler {
	return &QuotesHandler{service: service}
}

func (handler QuotesHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/quotes", handler)
}

func (handler QuotesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "Invalid quote id format", http.StatusBadRequest)
		log.Println(err)
		return
	}
	quote, err := handler.service.GetQuote(uint(id))
	if err != nil {
		http.Error(w, "Quote with given id does not exist", http.StatusBadRequest)
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
	quote, err := handler.service.GetRandomQuote()
	if err != nil {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
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
	quotes, err := handler.service.GetAllQuotes()
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
