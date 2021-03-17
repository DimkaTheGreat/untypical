package main

import (
	"flag"
	"untypicalCompanyTestTask/pkg/repository"
	"untypicalCompanyTestTask/pkg/routing"
)

//Порт, на котором разворачивается веб-сервис
var port = flag.String("port", "8086", "server port")

func main() {

	flag.Parse()

	storage := repository.NewStorage()

	routing.Run(storage, *port)

}
