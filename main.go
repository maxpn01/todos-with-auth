package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-app-with-auth/postgres"
	"todo-app-with-auth/router"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := router.Router(db)

	fmt.Println("Starting server at http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
