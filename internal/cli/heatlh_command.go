package cli

import (
	"fmt"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/config"
)

func Health(cfg *config.CliConfig) error {
	res, err := http.Get(cfg.Url() + "/health")

	if err != nil {
		return err
	}

	fmt.Println(res.Status)
	return nil
}
