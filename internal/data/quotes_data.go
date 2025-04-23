package data

import (
	"database/sql"
	"fmt"
	"log"

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
