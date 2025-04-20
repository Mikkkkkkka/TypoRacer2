package utils

import (
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
			return fmt.Errorf("escaped")
		}
		if key == keyboard.KeyEnter {
			break
		}
		if char == 0 {
			if length := builder.Len(); key == keyboard.KeyBackspace2 && length > 0 {
				restorer := builder.String()
				builder.Reset()
				builder.WriteString(restorer[:length-1])
			}
			continue
		}
		builder.WriteRune(char)
	}

	*secret = builder.String()
	fmt.Println()
	return nil
}
