package main

// TODO: Refactor (internal/api, /internal/service, /internal/data) to implement Controller-Service-Dao architecture more clearly

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/api"
	"github.com/Mikkkkkkka/typoracer/internal/api/handlers"
	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/internal/service"
)

var (
	host     string
	port     string
	user     string
	password string
	db       string
)

func initializeFlags() {
	flag.StringVar(&host, "host", "localhost", "The address of DB")
	flag.StringVar(&port, "port", "5432", "The port of DB")
	flag.StringVar(&user, "user", "postgres", "The user of DB to be logged in as")
	flag.StringVar(&password, "password", "mysecretpassword", "The password of said user")
	flag.StringVar(&db, "db", "typoracer", "The name of the db to use")
	flag.Parse()
}

func createConnectionString() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, db)
}

func main() {
	initializeFlags()

	db, err := data.ConnectDB(createConnectionString())
	if err != nil {
		log.Fatalln(err)
		return
	}

	userService := service.NewUserService(db)

	loginHandler := handlers.NewLoginHandler(&userService)
	registerHandler := handlers.NewRegisterHandler(&userService)
	usersHandler := handlers.NewUsersHandler(userService)
	quotesHandler := handlers.NewQuotesHandler(db)
	playsHandler := handlers.NewPlaysHandler(db)

	mux := api.NewRouter(db)

	loginHandler.RegisterRoutes(mux)
	registerHandler.RegisterRoutes(mux)
	usersHandler.RegisterRoutes(mux)
	quotesHandler.RegisterRoutes(mux)
	playsHandler.RegisterRoutes(mux)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
