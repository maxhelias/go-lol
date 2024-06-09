package endpoints

import (
	"encoding/json"
	"fmt"

	"github.com/maxhelias/golol/internal/api/models"
	"github.com/maxhelias/golol/pkg/client"
)

func GetCurrSummoner() (*models.CurrSummoner, error) {
	bts, err := client.Get("/lol-summoner/v1/current-summoner")
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	data := &models.CurrSummoner{}
	err = json.Unmarshal(bts, data)
	if nil != err {
		return nil, err
	}

	return data, nil
}
