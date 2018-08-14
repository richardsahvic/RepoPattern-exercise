package main

import (
	"log"
	"net/http"

	"routes"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	route := routes.Routes()

	http.Handle("/", route)
	log.Println("SERVER STARTED")

	http.ListenAndServe(":8080", route)
}
