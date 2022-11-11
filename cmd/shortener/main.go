package main

import (
	"log"
	"net/http"

	"github.com/izaake/go-shortener-tpl/internal/handlers"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}