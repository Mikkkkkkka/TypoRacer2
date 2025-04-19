package api

import (
	"database/sql"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/api/handlers"
)

type HttpHandlerFuncWithDB func(db *sql.DB) http.HandlerFunc

// Works kind of like Currying
func injectDB(handlerWithDB HttpHandlerFuncWithDB, db *sql.DB) http.HandlerFunc {
	return handlerWithDB(db)
}

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/v1/login", injectDB(handlers.LoginHandlerWithDB, db))
	mux.HandleFunc("/api/v1/register", injectDB(handlers.RegisterHandlerWithDB, db))
	mux.HandleFunc("/api/v1/users/{id}", injectDB(handlers.UsersHandlerWithDB, db))
	mux.HandleFunc("/api/v1/quotes", injectDB(handlers.QuotesHandlerWithDB, db))
	mux.HandleFunc("/api/v1/plays", injectDB(handlers.PlaysHandlerWithDB, db))

	return mux
}
