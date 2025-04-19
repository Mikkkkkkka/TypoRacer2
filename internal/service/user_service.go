package service

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

func LoginUser(username, password string, db *sql.DB) (*model.User, error) {
	user, err := data.GetUserWithoutTokenByUsername(username, db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("LoginUser: no user with username \"%s\"", username)
	}

	for range 5 {
		token := generateToken()
		expiration := time.Now().Add(15 * time.Minute)
		if err = data.AddTokenToUser(token, expiration, user.Id, db); err == nil {
			log.Println(err)
			return data.GetUserById(user.Id, db)
		}
	}

	return nil, fmt.Errorf("LoginUser: failed to generate token")
}

func AuthorizeUser(token string, db *sql.DB) (*model.User, error) {
	user, err := data.GetUserFromToken(token, db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("AuthorizeUser: no user with token")
	}
	if time.Now().After(user.TokenExpiration) {
		log.Println("token has expired")
		return nil, fmt.Errorf("AuthorizeUser: token has expired")
	}
	return user, nil
}

const tokenChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func generateToken() string {
	token := make([]byte, 20)
	for i := range 20 {
		token[i] = tokenChars[rand.IntN(len(tokenChars))]
	}
	return string(token)
}

func CalculateStats(user *model.User, db *sql.DB) (*model.UserStats, error) {
	plays, err := data.GetPlaysByUserId(user.Id, db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("CalculateStats: %w", err)
	}
	var stats model.UserStats
	quotesPlayed := make(map[uint]uint, 10)
	for _, play := range *plays {
		stats.WordsPerMinute += play.WordsPerMinute
		stats.Consistency += play.Consistency
		stats.Accuracy += play.Accuracy
		quotesPlayed[play.QuoteId]++
	}
	if playsCount := float32(len(*plays)); playsCount != 0 {
		stats.WordsPerMinute /= playsCount
		stats.Consistency /= playsCount
		stats.Accuracy /= playsCount
	}

	keys := make([]uint, 0, len(quotesPlayed))
	for k := range quotesPlayed {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return quotesPlayed[uint(i)] > quotesPlayed[uint(j)]
	})
	if len(keys) != 0 {
		stats.FavoriteQuote = keys[0]
	}

	return &stats, nil
}
