package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/service"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

type LoginHandler struct {
	db *sql.DB
}

func NewLoginHandler(db *sql.DB) *LoginHandler {
	return &LoginHandler{db: db}
}

func (handler LoginHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/login", handler)
}

func (handler LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed for /api/v1/login")
		return
	}

	var payloadData model.LoginInfo

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer r.Body.Close()

	user, err := service.LoginUser(payloadData.Username, payloadData.Password, handler.db)
	if err != nil && err.Error() != "LoginUser: failed to generate token" {
		http.Error(w, "Incorrect password or login", http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.Header().Add("Authorization", "Bearer: "+user.Token)
}
