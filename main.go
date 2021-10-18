package main

import (
  "github.com/mirzaakhena/gogen/application"
  "github.com/mirzaakhena/gogen/application/registry"
)

func main() {
  application.Run(registry.NewGogen2()())
}
