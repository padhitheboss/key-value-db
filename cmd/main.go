package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/padhitheboss/key-value-db/pkg/routes"
)

var PORT = "3000"

func main() {
	var Router = chi.NewRouter()
	routes.RegisterRoute(Router)
	fmt.Printf("Starting server on Port %v\n", PORT)
	err := http.ListenAndServe(":"+PORT, Router)
	if err != nil {
		log.Panic(err)
	}
}
