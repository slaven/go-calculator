package main

import (
	"go-calculator/calcserver"
	"log"
	"net/http"
	"time"
)

func main() {

	calcHandler := calcserver.Create()

	srv := http.Server{
		Addr:        ":8080",
		Handler:     calcHandler,
		ReadTimeout: time.Duration(30 * time.Second),
		IdleTimeout: time.Duration(30 * time.Second),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}