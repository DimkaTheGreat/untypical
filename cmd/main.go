package main

import (
	"flag"

	"github.com/DimkaTheGreat/untypical/pkg/repository"
	"github.com/DimkaTheGreat/untypical/pkg/routing"
)

//Порт, на котором разворачивается веб-сервис
var port = flag.String("port", "8086", "server port")

func main() {

	flag.Parse()

	storage := repository.NewStorage()

	routing.Run(storage, *port)

}
