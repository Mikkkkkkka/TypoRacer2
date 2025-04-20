package service

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
	"github.com/Mikkkkkkka/typoracer/pkg/utils"
	"github.com/eiannone/keyboard"
)

func CalculatePlayResults(userId uint, record *model.PlayRecord, db *sql.DB) (*model.Play, error) {
	play := model.Play{
		UserId:  userId,
		QuoteId: record.QuoteId,
	}
	quote, err := data.GetQuote(record.QuoteId, db)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("CalculatePlayResults.data.GetQuote: %w", err)
	}
	play.WordsPerMinute = calculateWordsPerMinute(record)
	play.Accuracy = calculateAccuracy(quote, record) * 100.0
	play.Consistency = calculateConsistency(record) * 100.0
	return &play, nil
}

func calculateWordsPerMinute(record *model.PlayRecord) float32 {
	times := make([]time.Duration, len(record.KeyStream))
	for _, v := range record.KeyStream {
		times = append(times, v.ElapsedTime)
	}
	totalTime := utils.Sum(times)
	return float32(countWordsInInput(record)) / (float32(totalTime) / float32(time.Minute))
}

func calculateAccuracy(quote *model.Quote, record *model.PlayRecord) float32 {
	mistakeCount := 0
	text := []rune(quote.Text)
	for i, j := 0, 0; i < len(quote.Text) && j < len(record.KeyStream); j++ {
		if record.KeyStream[j].KeyEvent.Rune != text[i] &&
			record.KeyStream[j].KeyEvent.Rune == 0 &&
			record.KeyStream[j].KeyEvent.Key != keyboard.KeyBackspace2 { // TODO: refactor (this is ridiculous)
			mistakeCount++
			continue
		}
		i++
	}
	if mistakeCount > len(quote.Text) {
		return 0
	}
	return 1.0 - float32(mistakeCount)/float32(len(text))
}

// FIXME: unexpectedly returns values over 100
func calculateConsistency(record *model.PlayRecord) float32 {
	n := len(record.KeyStream)
	times := make([]time.Duration, n)
	for _, v := range record.KeyStream {
		times = append(times, v.ElapsedTime)
	}
	averageTime := time.Duration(math.Round(utils.Average(times)))
	for i, v := range times {
		times[i] = averageTime - v
	}
	return 1.0 - (float32(utils.Sum(times))/float32(n-1))/float32(n)
}

func countWordsInInput(record *model.PlayRecord) int {
	currentWord := 0
	wordLengths := make([]int, 0, 20)
	for _, keyPress := range record.KeyStream {
		if keyPress.KeyEvent.Rune == ' ' {
			wordLengths = append(wordLengths, currentWord)
			currentWord = 0
			continue
		}
		if keyPress.KeyEvent.Rune == 0 && keyPress.KeyEvent.Key == keyboard.KeyBackspace {
			currentWord--
			if currentWord < 0 {
				currentWord = wordLengths[len(wordLengths)-1]
				wordLengths = slices.Delete(wordLengths, len(wordLengths)-1, len(wordLengths))
			}
			continue
		}
		currentWord++
	}
	if currentWord != 0 {
		wordLengths = append(wordLengths, currentWord)
	}
	return len(wordLengths)
}
