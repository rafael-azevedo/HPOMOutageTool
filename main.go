package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rafael-azevedo/HPOMOutageTool/router"
)

func main() {
	router := router.BuildRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:1234",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
