package main

import (
	"accounting/application"
	"accounting/application/registry"
	"flag"
)

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "appGinGonic":
		application.Run(registry.NewAppGinGonic())

	case "appNetHttp":
		application.Run(registry.NewAppNetHttp())

	case "appLabstackEcho":
		application.Run(registry.NewAppLabstackEcho())

	case "appGorrilaMux":
		application.Run(registry.NewAppGorrilaMux())

	case "appJulienSchmidt":
		application.Run(registry.NewAppJulienSchmidt())

	case "appGoChi":
		application.Run(registry.NewAppGoChi())
	}
}
