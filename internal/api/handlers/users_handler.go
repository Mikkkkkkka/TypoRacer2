package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/internal/service"
)

func UsersHandlerWithDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Invalid user id format", http.StatusBadRequest)
			return
		}

		user, err := data.GetUserWithoutTokenById(userId, db)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "User with id does not exist", http.StatusBadRequest)
			return
		}

		stats, err := service.CalculateStats(user, db)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(stats)
	}
}
