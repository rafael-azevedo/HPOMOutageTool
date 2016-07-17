package main

import (
	"log"
	"net/http"
)

func main() {
	router := BuildRouter()
	log.Fatal(http.ListenAndServe(":1234", router))
}
