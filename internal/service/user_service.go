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

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

const tokenChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func generateToken() string {
	token := make([]byte, 20)
	for i := range 20 {
		token[i] = tokenChars[rand.IntN(len(tokenChars))]
	}
	return string(token)
}

func (service UserService) LoginUser(username, password string) (*model.User, error) {
	user, err := data.GetUserWithoutTokenByUsername(username, service.db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("LoginUser: no user with username \"%s\"", username)
	}

	for range 5 {
		token := generateToken()
		expiration := time.Now().Add(15 * time.Minute)
		if err = data.AddTokenToUser(token, expiration, user.Id, service.db); err == nil {
			return data.GetUserById(user.Id, service.db)
		}
	}

	log.Println(err)
	return nil, fmt.Errorf("LoginUser: failed to generate token")
}

func (service UserService) RegisterUser(username, password string) error {
	if err := data.AddUser(username, password, service.db); err != nil {
		return fmt.Errorf("UserService.RegisterUser.data.AddUser: %w", err)
	}
	return nil
}

func (service UserService) AuthorizeUser(token string) (*model.User, error) {
	user, err := data.GetUserFromToken(token, service.db)
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

func (service UserService) CalculateStats(userId uint) (*model.UserStats, error) {
	plays, err := data.GetPlaysByUserId(userId, service.db)
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
