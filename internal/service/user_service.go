package service

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

func CalculateStats(user *model.User, db *sql.DB) (*model.UserStats, error) {
	plays, err := data.GetPlaysByUserId(user.Id, db)
	if err != nil {
		return nil, fmt.Errorf("CalculateStats: %w", err)
	}
	var stats model.UserStats
	quotesPlayed := make(map[int]int, 10)
	for _, play := range *plays {
		stats.WordsPerMinute += play.WordsPerMinute
		stats.Consistency += play.Consistency
		stats.Accuracy += play.Accuracy
		quotesPlayed[play.QuoteId]++
	}
	plays_count := float32(len(*plays))
	if plays_count != 0 {
		stats.WordsPerMinute /= plays_count
		stats.Consistency /= plays_count
		stats.Accuracy /= plays_count
	}

	keys := make([]int, 0, len(quotesPlayed))
	for k := range quotesPlayed {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return quotesPlayed[i] > quotesPlayed[j]
	})
	if len(keys) != 0 {
		stats.FavoriteQuote = keys[0]
	}

	return &stats, nil
}
