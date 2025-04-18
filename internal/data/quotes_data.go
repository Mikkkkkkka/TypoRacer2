package data

import (
	"database/sql"
	"errors"
	"math/rand/v2"

	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

func GetQuote(db *sql.DB, id int) (*model.Quote, error) {
	var quote model.Quote
	err := db.QueryRow("SELECT * FROM quotes WHERE id=$1", id).Scan(&quote.Id, &quote.Text)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func GetRandomQuote(db *sql.DB) (*model.Quote, error) {
	var quotesCount int
	err := db.QueryRow("SELECT count(*) FROM quotes").Scan(&quotesCount)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id FROM quotes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("no quotes are in the database yet")
	}

	quoteIndex := rand.IntN(quotesCount)
	for i := 0; i < quoteIndex && rows.Next(); i++ {
	}

	var randomQuoteId int
	err = rows.Scan(&randomQuoteId)
	if err != nil {
		return nil, err
	}

	return GetQuote(db, randomQuoteId)
}

func GetAllQuotes(db *sql.DB) (*[]model.Quote, error) {
	rows, err := db.Query("SELECT * FROM quotes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quotes := make([]model.Quote, 0)
	for rows.Next() {
		var quote model.Quote
		if err = rows.Scan(&quote.Id, &quote.Text); err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return &quotes, nil
}
