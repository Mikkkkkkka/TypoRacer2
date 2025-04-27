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

type PlayService struct {
	Db *sql.DB
}

func NewPlayService(Db *sql.DB) *PlayService {
	return &PlayService{Db: Db}
}

func (service *PlayService) GetPlaysByUserId(userId uint) (*[]model.Play, error) {
	return data.GetPlaysByUserId(userId, service.Db)
}

func (service *PlayService) GetAllPlays() (*[]model.Play, error) {
	return data.GetAllPlays(service.Db)
}

func (service *PlayService) AddPlay(play *model.Play) error {
	return data.AddPlay(play, service.Db)
}

func (service PlayService) CalculatePlayResults(userId uint, record *model.PlayRecord) (*model.Play, error) {
	play := model.Play{
		UserId:  userId,
		QuoteId: record.QuoteId,
	}
	quote, err := data.GetQuote(record.QuoteId, service.Db)
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
			!(record.KeyStream[j].KeyEvent.Rune == 0 &&
				record.KeyStream[j].KeyEvent.Key == keyboard.KeyBackspace2) { // TODO: refactor (this is ridiculous)
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

func calculateConsistency(record *model.PlayRecord) float32 {
	n := float64(len(record.KeyStream))
	msTimes := make([]time.Duration, 0, int64(n))
	for _, keyPress := range record.KeyStream {
		msTimes = append(msTimes, keyPress.ElapsedTime/time.Millisecond)
	}
	deviations := make([]float64, 0, int(n))
	averageTime := utils.Average(msTimes)
	for _, time := range msTimes {
		deviations = append(deviations, math.Pow(float64(time)-averageTime, 2))
	}
	return 1.0 - float32(math.Sqrt(utils.Sum(deviations)/(n))/averageTime)
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
