package main

import (
	"accounting/application"
	"accounting/application/registry"
	"flag"
	"fmt"
)

func main() {

	appMap := map[string]func() application.RegistryContract{
		"appGinGonic":      registry.NewAppGinGonic(),
		"appNetHttp":       registry.NewAppNetHttp(),
		"appLabstackEcho":  registry.NewAppLabstackEcho(),
		"appGorrilaMux":    registry.NewAppGorrilaMux(),
		"appJulienSchmidt": registry.NewAppJulienSchmidt(),
		"appGoChi":         registry.NewAppGoChi(),
	}

	flag.Parse()

	app, exist := appMap[flag.Arg(0)]
	if exist {
		application.Run(app())
	} else {
		fmt.Println("You may try this app name:")
		for appName, _ := range appMap {
			fmt.Printf("%s\n", appName)
		}
	}

}
