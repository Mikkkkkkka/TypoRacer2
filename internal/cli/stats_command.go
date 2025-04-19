package cli

import (
	"net/http"

	"github.com/Mikkkkkkka/typoracer/internal/config"
)

func Stats(cfg *config.CliConfig) {

	http.Get(cfg.Url() + "/api/v1/users")
}
