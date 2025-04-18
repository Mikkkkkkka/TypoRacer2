package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
)

func ScanSecret(secret *string) error {

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	var builder strings.Builder

	for {
		char, key, err := keyboard.GetKey()

		if err != nil {
			return err
		}
		if key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
			return errors.New("escaped")
		}
		if key == keyboard.KeyEnter {
			break
		}
		if char == 0 {
			continue
		}
		fmt.Print("*")
		builder.WriteRune(char)
	}

	*secret = builder.String()
	fmt.Println()
	return nil
}
