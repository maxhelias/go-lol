package main

import (
	"log"

	"github.com/maxhelias/golol/internal/app"
)

func main() {
	app, err := app.New(app.WithDebug())
	if err != nil {
		log.Fatal(err)

		return
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
