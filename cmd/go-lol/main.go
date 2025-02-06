package main

import (
	"log"

	"github.com/maxhelias/golol/internal/api/endpoints"
	"github.com/maxhelias/golol/internal/logger"
	"github.com/maxhelias/golol/pkg/client"
	"github.com/maxhelias/golol/pkg/process"
)

func main() {
	/*app, err := app.New(app.WithDebug())
	if err != nil {
		log.Fatal(err)

		return
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}*/

	appLogger, err := logger.NewLogger(true)
	if err != nil {
		log.Fatal(err)
	}

	lcuInfo, err := process.FindLcuConnectInfo()
	if err != nil {
		appLogger.Error(err.Error())

		return
	}

	// Create a new LCU client
	options := []client.Option{}
	options = append(options, client.WithPort(lcuInfo.Port))
	options = append(options, client.WithAuthToken(lcuInfo.AuthToken))

	client.NewClient(options...)

	endpoints.GetCurrSummoner()

	// Dump currentSummoner
	//fmt.Println(currentSummoner)
}
