package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/pkg/model/requests"
	"github.com/Mikkkkkkka/typoracer/pkg/utils"
)

func Register() error {

	var username, password string

	fmt.Print("Username: ")
	if _, err := fmt.Scan(&username); err != nil {
		panic(err)
	}

	fmt.Print("Password: ")
	if err := utils.ScanSecret(&password); err != nil {
		panic(err)
	}

	err := requestRegistration(username, password)
	if err != nil {
		return err
	}

	return nil
}

func requestRegistration(username, password string) error {

	if username == "" || password == "" {
		return errors.New("username or password is empty")
	}

	body := requests.RegistrationRequestBody{
		Username: username,
		Password: password,
	}
	// TODO: add password hashing

	jsonBody, _ := json.Marshal(body)

	// TODO: make server uri configurable from client
	res, err := http.Post("http://localhost:8080/api/v1/register", "json", bytes.NewBuffer(jsonBody))
	fmt.Println("Warning! using localhost:8080")

	if err != nil {
		return err
	}

	if res.Status != "200 OK" {
		return fmt.Errorf("unexpected response code: %s", res.Status)
	}

	fmt.Printf("User \"%s\" was registered successfully!\n", username)
	return nil
}
