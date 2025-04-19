package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/pkg/model/requests"
)

func RegisterHandlerWithDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Default().Println("Method not allowed for api/v1/register")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payloadData requests.LoginInfo

		if err := json.NewDecoder(r.Body).Decode(&payloadData); err != nil {
			log.Default().Println(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := data.AddUser(payloadData.Username, payloadData.Password, db); err != nil {
			log.Default().Println(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
}
