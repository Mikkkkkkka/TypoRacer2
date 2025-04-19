package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/config"
	"github.com/Mikkkkkkka/typoracer/pkg/model/requests"
	"github.com/Mikkkkkkka/typoracer/pkg/utils"
)

func Register(cfg *config.CliConfig) error {

	var username, password string

	fmt.Print("Username: ")
	if _, err := fmt.Scan(&username); err != nil {
		panic(err)
	}

	fmt.Print("Password: ")
	if err := utils.ScanSecret(&password); err != nil {
		panic(err)
	}

	err := requestRegistration(username, password, cfg)
	if err != nil {
		return err
	}

	fmt.Printf("User \"%s\" was registered successfully!\n", username)
	return nil
}

func requestRegistration(username, password string, cfg *config.CliConfig) error {

	if username == "" || password == "" {
		return errors.New("username or password is empty")
	}

	body := requests.LoginInfo{
		Username: username,
		Password: password,
	}
	// TODO: add password hashing

	jsonBody, _ := json.Marshal(body)

	res, err := http.Post(cfg.Url()+"api/v1/register", "json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	if res.Status != "200 OK" {
		return fmt.Errorf("unexpected response code: %s", res.Status)
	}

	return nil
}
