package data

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

func GetQuote(id uint, db *sql.DB) (*model.Quote, error) {
	var quote model.Quote
	err := db.QueryRow("SELECT * FROM quotes WHERE id=$1", id).Scan(&quote.Id, &quote.Text)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetQuote: %w", err)
	}
	return &quote, nil
}

func GetRandomQuote(db *sql.DB) (*model.Quote, error) {
	var quotesCount int
	err := db.QueryRow("SELECT count(*) FROM quotes").Scan(&quotesCount)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetRandomQuote.db.QueryRow: %w", err)
	}

	rows, err := db.Query("SELECT id FROM quotes")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	quoteIndex := rand.IntN(quotesCount)
	for i := 0; rows.Next() && i < quoteIndex; i++ {
	}

	var randomQuoteId uint
	err = rows.Scan(&randomQuoteId)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetRandomQuote.rows.Scan: %w", err)
	}

	return GetQuote(randomQuoteId, db)
}

func GetAllQuotes(db *sql.DB) (*[]model.Quote, error) {
	rows, err := db.Query("SELECT * FROM quotes")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	quotes := make([]model.Quote, 0)
	for rows.Next() {
		var quote model.Quote
		if err = rows.Scan(&quote.Id, &quote.Text); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("GetAllQuotes.rows.Scan: %w", err)
		}
		quotes = append(quotes, quote)
	}

	return &quotes, nil
}
