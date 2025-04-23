package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/pkg/model/requests"
)

type RegisterHandler struct {
	db *sql.DB
}

func NewRegisterHandler(db *sql.DB) *RegisterHandler {
	return &RegisterHandler{db: db}
}

func (handler RegisterHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/register", handler)
}

func (handler RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed for api/v1/register")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payloadData requests.LoginInfo

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := data.AddUser(payloadData.Username, payloadData.Password, handler.db); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
