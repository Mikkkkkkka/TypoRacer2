package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Mikkkkkkka/typoracer/internal/config"
	"github.com/Mikkkkkkka/typoracer/pkg/model"
	"github.com/Mikkkkkkka/typoracer/pkg/utils"
	"github.com/eiannone/keyboard"
	"github.com/muesli/termenv"
)

func Play(cfg *config.CliConfig) error {
	var username, password string

	fmt.Print("Username: ")
	if _, err := fmt.Scan(&username); err != nil {
		panic(err)
	}

	fmt.Print("Password: ")
	if err := utils.ScanSecret(&password); err != nil {
		panic(err)
	}

	token, err := requestLogin(username, password, cfg)
	if err != nil {
		return fmt.Errorf("requestLogin: %w", err)
	}
	quote, err := requestQuote(cfg)
	if err != nil {
		return fmt.Errorf("requestQuote: %w", err)
	}

	record, err := playScene(quote)
	if err != nil {
		return fmt.Errorf("playScene: %w", err)
	}

	play, err := sendPlayRecord(token, record, cfg)
	if err != nil {
		return fmt.Errorf("sendPlayRecord: %w", err)
	}

	fmt.Println("Stats:")
	fmt.Println()
	fmt.Println("QuoteId:", quote.Id)
	fmt.Println("WPM:", play.WordsPerMinute)
	fmt.Println("Accuracy:", play.Accuracy)
	fmt.Println("Consistency:", play.Consistency)

	return nil
}

func sendPlayRecord(token *string, record *model.PlayRecord, cfg *config.CliConfig) (*model.Play, error) {
	payloadData, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.Url()+"api/v1/plays", bytes.NewBuffer(payloadData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", *token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("unexpected response code: %s", res.Status)
	}

	responseData := model.Play{
		UserId:         0,
		QuoteId:        0,
		WordsPerMinute: 0,
		Accuracy:       0,
		Consistency:    0,
	}

	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return &responseData, err
}

func requestQuote(cfg *config.CliConfig) (*model.Quote, error) {
	res, err := http.Get(cfg.Url() + "api/v1/quotes?random=true")
	if err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("unexpected response code: %s", res.Status)
	}

	var quote model.Quote
	if err := json.NewDecoder(res.Body).Decode(&quote); err != nil {
		return nil, err
	}

	return &quote, nil
}

func requestLogin(username, password string, cfg *config.CliConfig) (*string, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}

	body := model.LoginInfo{
		Username: username,
		Password: password,
	}
	// TODO: add password hashing

	payloadData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(cfg.Url()+"api/v1/login", "json", bytes.NewBuffer(payloadData))
	if err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("unexpected response code: %s", res.Status)
	}

	token := res.Header.Get("Authorization")

	return &token, nil
}

func playScene(quote *model.Quote) (*model.PlayRecord, error) {
	record := model.PlayRecord{
		QuoteId: quote.Id,
	}

	if err := keyboard.Open(); err != nil {
		return nil, err
	}
	defer keyboard.Close()

	output := termenv.NewOutput(os.Stdout)
	output.AltScreen()
	defer output.ExitAltScreen()

	var builder strings.Builder
	for {
		drawScreen(quote, &builder, output)
		prevKeyTime := time.Now()
		char, key, err := keyboard.GetSingleKey()
		record.KeyStream = append(record.KeyStream, model.KeyPress{
			KeyEvent:    keyboard.KeyEvent{Key: key, Rune: char, Err: err},
			ElapsedTime: time.Since(prevKeyTime),
		})
		if char == 0 {
			if key == keyboard.KeyBackspace2 {
				restorer := builder.String()
				builder.Reset()
				builder.WriteString(restorer[:len(restorer)-1])
			} else if key == keyboard.KeySpace {
				builder.WriteRune(' ')
			} else if key == keyboard.KeyEnter {
				break
			}
			continue
		}
		builder.WriteRune(char)
	}

	record.KeyStream[0].ElapsedTime = 0
	return &record, nil
}

func drawScreen(quote *model.Quote, builder *strings.Builder, output *termenv.Output) {
	input := []rune(builder.String())
	text := []rune(quote.Text)

	var correct, mistake, unwritten strings.Builder
	output.ClearScreen()
	fmt.Println(quote.Text)
	lastCorrect := -1
	for i, v := range input {
		if v != text[i] {
			break
		}
		lastCorrect = i
		correct.WriteRune(v)
	}
	if lastCorrect != len(input) {
		mistake.WriteString(string(input[lastCorrect+1:]))
	}
	if lastCorrect != len(text) {
		unwritten.WriteString(string(text[lastCorrect+1:]))
	}
	fmt.Printf("%s%s%s",
		output.String(correct.String()).Foreground(output.Color("#ffffff")),
		output.String(mistake.String()).Foreground(output.Color("#ff0000")),
		output.String(unwritten.String()).Foreground(output.Color("#888888")),
	)
}
