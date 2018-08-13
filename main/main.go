package main

import (
	"log"
	"net/http"

	"datasource"
	"repo"
	"routes"
	"service"

	_ "github.com/go-sql-driver/mysql"
)

var userService service.UserService

func main() {
	db := datasource.InitConnection()

	r := repo.NewRepository(db)

	_ = service.NewUserService(r)

	route := routes.Routes()

	http.Handle("/", route)
	log.Println("SERVER STARTED")

	http.ListenAndServe(":8080", route)
}
