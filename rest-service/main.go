package main

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stackpath/backend-developer-tests/rest-service/pkg/api"
	"log"
	"net/http"
)

var listenAddr = ":8080"

func main() {
	flag.Parse()

	fmt.Println("SP// Backend Developer Test - RESTful Service")
	fmt.Println()

	restAPI := api.New()

	router := httprouter.New()
	router.GET("/people", restAPI.RequestLogger(restAPI.SearchPeople))
	router.GET("/people/:id", restAPI.RequestLogger(restAPI.GetPerson))

	log.Fatalln(http.ListenAndServe(listenAddr, router))
}

func init() {
	flag.StringVar(&listenAddr, "listenAddr", listenAddr, "The address to listen on passed into ListenAndServe.")
}
