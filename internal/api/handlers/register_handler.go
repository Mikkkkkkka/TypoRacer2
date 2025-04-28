package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/service"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

type RegisterHandler struct {
	service *service.UserService
}

func NewRegisterHandler(service *service.UserService) *RegisterHandler {
	return &RegisterHandler{service: service}
}

func (handler RegisterHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /api/v1/register", handler)
}

func (handler RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payloadData model.LoginInfo

	if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := handler.service.RegisterUser(payloadData.Username, payloadData.Password); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
