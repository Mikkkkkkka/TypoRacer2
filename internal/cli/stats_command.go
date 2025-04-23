package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/config"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
)

func Stats(cfg *config.CliConfig) error {

	var userId int
	fmt.Print("Enter UserId: ")
	fmt.Scan(&userId)

	resp, err := http.Get(cfg.Url() + "api/v1/users/" + fmt.Sprintf("%d", userId))
	if err != nil {
		return fmt.Errorf("Stats.http.Get: %w", err)
	}

	if resp.Status != "200 OK" {
		return fmt.Errorf("unexpected response code: %s", resp.Status)
	}

	var stats model.UserStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return fmt.Errorf("Stats.json.Decode: %w", err)
	}

	fmt.Println("UserId: " + fmt.Sprintf("%d", userId))
	fmt.Println("Average WPM: " + fmt.Sprintf("%f", stats.WordsPerMinute))
	fmt.Println("Average Accuracy: " + fmt.Sprintf("%f", stats.Accuracy))
	fmt.Println("Average Consistency: " + fmt.Sprintf("%f", stats.Consistency))
	fmt.Println("Favorite quote: " + fmt.Sprintf("%d", stats.FavoriteQuote))
	return nil
}
