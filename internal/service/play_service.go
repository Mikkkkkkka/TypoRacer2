package service

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
	"github.com/Mikkkkkkka/typoracer/pkg/utils"
)

func CalculatePlayResults(record *model.PlayRecord, db *sql.DB) (*model.Play, error) {
	play := model.Play{
		UserId:  record.UserId,
		QuoteId: record.QuoteId,
	}
	quote, err := data.GetQuote(record.QuoteId, db)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("CalculatePlayResults.data.GetQuote: %w", err)
	}
	play.WordsPerMinute = calculateWordsPerMinute(quote, record)
	play.Accuracy = calculateAccuracy(quote, record) * 100.0
	play.Consistency = calculateConsistency(record) * 100.0
	return &play, nil
}

// FIXME: total remake needed
func calculateWordsPerMinute(quote *model.Quote, record *model.PlayRecord) float32 {
	durations := make([]time.Duration, len(record.KeyStream))
	for _, v := range record.KeyStream {
		durations = append(durations, v.ElapsedTime)
	}
	return float32(len(strings.Split(string(quote.Text), " "))) / float32(time.Duration(utils.Sum(durations))) / float32(time.Minute)
}

func calculateAccuracy(quote *model.Quote, record *model.PlayRecord) float32 {
	mistakeCount := 0
	text := []rune(quote.Text)
	for i, j := 0, 0; i < len(quote.Text) && j < len(record.KeyStream); j++ {
		if text[i] != record.KeyStream[j].KeyEvent.Rune {
			mistakeCount++
			continue
		}
		i++
	}
	if mistakeCount > len(quote.Text) {
		return 0
	}
	return 1.0 - float32(mistakeCount/len(quote.Text))
}

// FIXME: find out why perfect consistency results in returning 0
func calculateConsistency(record *model.PlayRecord) float32 {
	n := len(record.KeyStream)
	durations := make([]time.Duration, n)
	for _, v := range record.KeyStream {
		durations = append(durations, v.ElapsedTime)
	}
	averageDuration := time.Duration(math.Round(utils.Average(durations)))
	for i, v := range durations {
		durations[i] = averageDuration - v
	}
	return (float32(utils.Sum(durations)) / float32(n-1)) / float32(n)
}
