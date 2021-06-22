package main

import (
  "flag"
  "fmt"

  "github.com/mirzaakhena/gogen/gogencommand"
)

func main() {

  flag.Parse()

  cmds := map[string]func() (gogencommand.Commander, error){
    "init":        gogencommand.NewInitializeModel,
    "usecase":     gogencommand.NewUsecaseModel,
    "test":        gogencommand.NewTestModel,
    "entity":      gogencommand.NewEntityModel,
    "method":      gogencommand.NewMethodModel,
    "enum":        gogencommand.NewEnumModel,
    "error":       gogencommand.NewErrorModel,
    "valueobject": gogencommand.NewValueObjectModel,
    "valuestring": gogencommand.NewValueStringModel,
    "repository":  gogencommand.NewRepositoryModel,
    "service":     gogencommand.NewServiceModel,
    "gateway":     gogencommand.NewGatewayModel,
    "controller":  gogencommand.NewControllerModel,
    "registry":    gogencommand.NewRegistryModel,
  }

  var obj gogencommand.Commander
  var err error

  f, ok := cmds[flag.Arg(0)]
  if !ok {
    fmt.Printf("ERROR : %v", "Command is not recognized. Maybe you can try : \n")
    fmt.Println("gogen usecase CreateOrder")
    fmt.Println("gogen repository SaveOrder Order")
    fmt.Println("gogen gateway inmemory CreateOrder")
    fmt.Println("gogen controller openapi CreateOrder")
    fmt.Println("gogen registry appone openapi CreateOrder inmemory")

    return
  }

  obj, err = f()
  if err != nil {
    fmt.Printf("ERROR : %v\n", err.Error())
    return
  }

  err = obj.Run()
  if err != nil {
    fmt.Printf("ERROR : %v\n", err.Error())
    return
  }

}
