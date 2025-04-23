package service

import (
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/Mikkkkkkka/typoracer/internal/data"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

type QuoteService struct {
	db *sql.DB
}

func NewQuoteService(db *sql.DB) QuoteService {
	return QuoteService{db: db}
}

func (service *QuoteService) GetAllQuotes() (*[]model.Quote, error) {
	return data.GetAllQuotes(service.db)
}

func (service *QuoteService) GetRandomQuote() (*model.Quote, error) {
	quotes, err := data.GetAllQuotes(service.db)
	if err != nil {
		return nil, fmt.Errorf("QuoteService.GetRandomQuote.data.GetAllQuotes: %w", err)
	}
	return &(*quotes)[rand.IntN(len(*quotes))], nil
}

func (service *QuoteService) GetQuote(id uint) (*model.Quote, error) {
	return data.GetQuote(id, service.db)
}
