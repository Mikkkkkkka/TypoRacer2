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
	"github.com/Mikkkkkkka/typoracer/pkg/model/requests"
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

	user, err := requestLogin(username, password, cfg)
	if err != nil {
		return err
	}
	quote, err := requestQuote(cfg)
	if err != nil {
		return err
	}

	_, err = playScene(user, quote)
	if err != nil {
		return err
	}

	return nil
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
	json.NewDecoder(res.Body).Decode(&quote)

	return &quote, nil
}

func requestLogin(username, password string, cfg *config.CliConfig) (*model.User, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}

	body := requests.LoginInfo{
		Username: username,
		Password: password,
	}
	// TODO: add password hashing

	payloadData, _ := json.Marshal(body)

	res, err := http.Post(cfg.Url()+"api/v1/login", "json", bytes.NewBuffer(payloadData))
	if err != nil {
		return nil, err
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("unexpected response code: %s", res.Status)
	}

	var user model.User
	json.NewDecoder(res.Body).Decode(&user)

	return &user, nil
}

func playScene(user *model.User, quote *model.Quote) (*model.PlayRecord, error) {
	record := model.PlayRecord{
		UserId:  user.Id,
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
