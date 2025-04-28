package data

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

// Attempting to prevent sql-injections
func isSecure(str string) bool {
	const potentiallyUnsafe = ";\"[]() *"
	return !strings.ContainsAny(str, potentiallyUnsafe)
}

func AddUser(username, password string, db *sql.DB) error {
	if !isSecure(username) || !isSecure(password) {
		log.Println("username or password contain insecure characters")
		return fmt.Errorf("AddUser: username or password contain insecure characters")
	}
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("AddUser: %w", err)
	}
	return nil
}

func GetUserById(userId uint, db *sql.DB) (*model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT id, username, password, token, token_expiration FROM users WHERE id=$1", userId).
		Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetUserById.db.QueryRow.Scan: %w", err)
	}
	return &user, nil
}

func GetUserByUsername(username string, db *sql.DB) (*model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).
		Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetUserByUsername.db.QueryRow.Scan: %w", err)
	}
	return &user, nil
}

func GetUserFromToken(token string, db *sql.DB) (*model.User, error) {
	if !isSecure(token) {
		log.Println("the provided token contains insecure characters")
		return nil, fmt.Errorf("GetUserFromToken: the provided token contains insecure characters")
	}
	var user model.User
	err := db.QueryRow("SELECT id, username, password, token, token_expiration FROM users WHERE token=$1", token).
		Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetUserFromToken.db.QueryRow.Scan: %w", err)
	}
	return &user, nil
}

func AddTokenToUser(token string, datetime time.Time, userId uint, db *sql.DB) error {
	if !isSecure(token) {
		log.Println("the generated token contains insecure characters")
		return fmt.Errorf("AddTokenToUser: the generated token contains insecure characters")
	}
	_, err := db.Exec("UPDATE users SET token=$1, token_expiration=$2 WHERE id=$3", token, datetime, userId)
	return err
}

func GetTokenExpiration(token string, db *sql.DB) (time.Time, error) {
	if !isSecure(token) {
		log.Println("the generated token contains insecure characters")
		return time.Now(), fmt.Errorf("GetTokenExpiration: the generated token contains insecure characters")
	}
	var expiration time.Time
	err := db.QueryRow("SELECT token_expiration FROM users WHERE token=$1", token).
		Scan(&expiration)
	if err != nil {
		return time.Now(), fmt.Errorf("GetTokenExpiration: %w", err)
	}
	return expiration, nil
}
