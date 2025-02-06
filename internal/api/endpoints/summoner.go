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

	var formattedJSON map[string]interface{}
	err = json.Unmarshal(bts, &formattedJSON)
	if err != nil {
		fmt.Println("Erreur lors du parsing JSON:", err)
	}

	// Affichage JSON formaté dans le terminal
	prettyJSON, _ := json.MarshalIndent(formattedJSON, "", "  ")
	fmt.Println(string(prettyJSON))

	data := &models.CurrSummoner{}
	err = json.Unmarshal(bts, data)
	if nil != err {
		return nil, err
	}

	return data, nil
}
