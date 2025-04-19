package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/service"
	"github.com/Mikkkkkkka/typoracer/pkg/model/requests"
)

func LoginHandlerWithDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payloadData requests.LoginInfo

		if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
			log.Fatal(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		user, err := service.LoginUser(payloadData.Username, payloadData.Password, db)
		if err != nil && err.Error() != "LoginUser: failed to generate token" {
			http.Error(w, "Incorrect password or login", http.StatusBadRequest)
			return
		}

		w.Header().Add("Authorization", "Bearer: "+user.Token)
	}
}
