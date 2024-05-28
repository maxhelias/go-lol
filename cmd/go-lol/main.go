package main

import (
	"fmt"
	"runtime"

	api "github.com/maxhelias/golol/internal"
	"github.com/maxhelias/golol/pkg/client"
	"github.com/maxhelias/golol/pkg/process"
)

var (
	version   string
	commit    string
	buildDate string
)

func main() {
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Build Date: %s\n", buildDate)

	fmt.Println("Architecture actuelle:", runtime.GOARCH)
	fmt.Println("Système d'exploitation actuel:", runtime.GOOS)

	lcuInfo, err := process.FindLcuConnectInfo()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("Processus trouvé: Id=%d, authToken=%s, port=%d\n", lcuInfo.ProcessID, lcuInfo.AuthToken, lcuInfo.Port)

	client.NewClient(lcuInfo.Port, lcuInfo.AuthToken)

	data, err := api.GetCurrSummoner()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(data)
}
