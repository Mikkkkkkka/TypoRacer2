package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

func AddPlay(play *model.Play, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO plays (user_id, quote_id, words_per_minute, accuracy, consistency) VALUES ($1, $2, $3, $4, $5)",
		play.UserId, play.QuoteId, play.WordsPerMinute, play.Accuracy, play.Consistency)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("AddPlay: %w", err)
	}
	return nil
}

func GetPlaysByUserId(id uint, db *sql.DB) (*[]model.Play, error) {
	rows, err := db.Query("SELECT * FROM plays WHERE user_id=$1", id)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetPlaysByUserId: %w", err)
	}
	defer rows.Close()

	plays := make([]model.Play, 0)
	for rows.Next() {
		var play model.Play
		if err = rows.Scan(
			&play.UserId,
			&play.QuoteId,
			&play.WordsPerMinute,
			&play.Accuracy,
			&play.Consistency); err != nil {
			log.Println(err)
			return nil, err
		}
		plays = append(plays, play)
	}

	return &plays, nil
}

func GetAllPlays(db *sql.DB) (*[]model.Play, error) {
	rows, err := db.Query("SELECT * FROM plays")
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GetAllPlays: %w", err)
	}
	defer rows.Close()

	plays := make([]model.Play, 0)
	for rows.Next() {
		var play model.Play
		if err = rows.Scan(
			&play.UserId,
			&play.QuoteId,
			&play.WordsPerMinute,
			&play.Accuracy,
			&play.Consistency); err != nil {
			log.Println(err)
			return nil, err
		}
		plays = append(plays, play)
	}

	return &plays, nil
}
