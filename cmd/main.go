package main

import (
	"github.com/erazr/omnibot/cmd/handlers"
)

func main() {

	err := handlers.RegisterCommands()
	if err != nil {
		return
	}

}
