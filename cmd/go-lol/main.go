package main

import (
	"fmt"

	"github.com/maxhelias/golol/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		fmt.Println(err)

		return
	}

	app.Run()

	/*data, err := endpoints.GetCurrSummoner()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(data)*/
}
