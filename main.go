package main

import (
	"flag"
	"untypicalCompanyTestTask/repository"
	"untypicalCompanyTestTask/routing"
)

var port = flag.String("port", "8086", "server port")

func main() {

	flag.Parse()

	storage := repository.NewStorage()

	routing.Run(storage, *port)

}
