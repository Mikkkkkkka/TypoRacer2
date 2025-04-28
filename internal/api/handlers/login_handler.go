package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/service"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

type LoginHandler struct {
	service *service.UserService
}

func NewLoginHandler(service *service.UserService) *LoginHandler {
	return &LoginHandler{service: service}
}

func (handler LoginHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /api/v1/login", handler)
}

func (handler LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payloadData model.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer r.Body.Close()

	user, err := handler.service.LoginUser(payloadData.Username, payloadData.Password)
	if err != nil && err.Error() != "LoginUser: failed to generate token" {
		http.Error(w, "Incorrect password or login", http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.Header().Add("Authorization", "Bearer: "+user.Token)
}
