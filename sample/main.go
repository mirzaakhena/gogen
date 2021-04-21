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


	}
}
