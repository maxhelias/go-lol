package api

import (
	"encoding/json"
	"fmt"

	"github.com/maxhelias/golol/pkg/client"
)

type (
	CurrSummoner struct {
		AccountID                   int64  `json:"accountId"`
		DisplayName                 string `json:"displayName"`
		InternalName                string `json:"internalName"`
		NameChangeFlag              bool   `json:"nameChangeFlag"`
		PercentCompleteForNextLevel int    `json:"percentCompleteForNextLevel"`
		ProfileIconID               int    `json:"profileIconId"`
		Puuid                       string `json:"puuid"`
		RerollPoints                struct {
			CurrentPoints    int `json:"currentPoints"`
			MaxRolls         int `json:"maxRolls"`
			NumberOfRolls    int `json:"numberOfRolls"`
			PointsCostToRoll int `json:"pointsCostToRoll"`
			PointsToReroll   int `json:"pointsToReroll"`
		} `json:"rerollPoints"`
		SummonerID       int64 `json:"summonerId"`
		SummonerLevel    int   `json:"summonerLevel"`
		Unnamed          bool  `json:"unnamed"`
		XpSinceLastLevel int   `json:"xpSinceLastLevel"`
		XpUntilNextLevel int   `json:"xpUntilNextLevel"`
	}
)

func GetCurrSummoner() (*CurrSummoner, error) {
	bts, err := client.Get("/lol-summoner/v1/current-summoner")
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	data := &CurrSummoner{}
	err = json.Unmarshal(bts, data)
	if nil != err {
		return nil, err
	}

	return data, nil
}
