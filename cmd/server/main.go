package main

// TODO: Refactor (internal/api, /internal/service, /internal/data) to implement Controller-Service-Dao architecture more clearly

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/api"
	"github.com/Mikkkkkkka/typoracer/internal/data"
)

var (
	host     string
	port     string
	user     string
	password string
	db       string
)

func initialiseFlags() {
	flag.StringVar(&host, "server", "localhost", "The adress of DB")
	flag.StringVar(&port, "port", "5432", "The port of DB")
	flag.StringVar(&user, "user", "postgres", "The user of DB to be logged in as")
	flag.StringVar(&password, "password", "mysecretpassword", "The password of said user")
	flag.StringVar(&db, "db", "typoracer", "The name of the db to use")
	flag.Parse()
}

func createConnectionString() string {
	return fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		user, db, password, host, port)
}

func main() {
	initialiseFlags()

	db, err := data.ConnectDB(createConnectionString())
	if err != nil {
		log.Fatal(err)
		return
	}
	mux := api.NewRouter(db)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
