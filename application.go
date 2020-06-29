package main

import (
	"fam-photos-server/router"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {

	r := router.Router()
	fmt.Println("Starting server on the port 5000...")

	corsHandler := handlers.CORS(
		handlers.AllowedHeaders(
			[]string{"X-Requested-With", "Content-Type", "Authorization"},
		),
		handlers.AllowedMethods(
			[]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		),
		handlers.AllowedOrigins(
			[]string{"*"},
		),
	)

	log.Fatal(http.ListenAndServe(":5000", corsHandler(r)))
}
