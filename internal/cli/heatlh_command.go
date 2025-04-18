package cli

import (
	"fmt"
	"net/http"
)

func Health() {
	res, err := http.Get("http://localhost:8080/health") // TODO: make uri configurable from client

	if err != nil {
		panic(err)
	}

	fmt.Println(res.Status)
}
