package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

func AddUser(username, password string, db *sql.DB) error {
	if !areSecure(username, password) {
		return errors.New("AddUser: username or password contain insecure characters")
	}
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	return err
}

func GetUserWithoutTokenById(userId int, db *sql.DB) (*model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE id=$1", userId).
		Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("GetUserById: %w", err)
	}
	return &user, nil
}

// Attempting to prevent sql-injections
func areSecure(username, password string) bool {
	const potentiallyUnsafe string = ";\"[]() *"
	var isUsernameSpoiled bool = strings.ContainsAny(username, potentiallyUnsafe)
	var isPasswordSpoiled bool = strings.ContainsAny(password, potentiallyUnsafe)
	return !(isUsernameSpoiled || isPasswordSpoiled)
}
