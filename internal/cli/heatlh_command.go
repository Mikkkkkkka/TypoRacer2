package cli

import (
	"fmt"
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/config"
)

func Health(cfg *config.CliConfig) {
	res, err := http.Get(cfg.Url() + "/health")

	if err != nil {
		panic(err)
	}

	fmt.Println(res.Status)
}
