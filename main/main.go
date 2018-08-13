package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"repo"
	"service"

	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type response struct {
	Message string `json:"message"`
}

type registerRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Msisdn   string `json:"msisdn"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Role     int    `json:"role"`
}

var userService service.UserService

func main() {
	db, err := sqlx.Connect("mysql", "dev:dev@(localhost:3306)/myapp?parseTime=true")
	if err != nil {
		log.Fatalln("Failed to connect to database,    ", err)
	}

	r := repo.NewRepository(db)

	userService = service.NewUserService(r)

	route := mux.NewRouter()
	route.HandleFunc("/login", loginHandler).Methods("POST")
	route.HandleFunc("/register", registerHandler).Methods("POST")
	route.HandleFunc("/viewprofile", profileHandler).Methods("POST")

	http.Handle("/", route)
	log.Println("SERVER STARTED")

	http.ListenAndServe(":8080", route)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 5000))

	var loginReq loginRequest
	json.Unmarshal(body, &loginReq)

	loginResult, err := userService.Login(loginReq.Username, loginReq.Password, loginReq.Role)
	if err != nil {
		log.Println("Failed at login,   ", err)
	}

	var loginResp response

	if len(loginResult) == 0 {
		loginResp.Message = "Login failed"
	} else {
		loginResp.Message = "Login Success"
	}

	js, err := json.Marshal(loginResp)
	if err != nil {
		log.Println("ERROR at login marshal,    ", err)
	}

	w.Header().Set("token", loginResult)
	w.Write(js)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 5000))

	var regRequest registerRequest
	json.Unmarshal(body, &regRequest)

	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("Fail to generate snowflake id,    ", err)
		return
	}

	id := node.Generate().String()

	userRegister := repo.User{
		ID:       id,
		Email:    regRequest.Email,
		Msisdn:   regRequest.Msisdn,
		Username: regRequest.Username,
		Password: regRequest.Password,
		Status:   0,
	}

	role := regRequest.Role

	registerResult, err := userService.Register(userRegister, role)
	if err != nil {
		log.Println("failed to register,    ", err)
	}

	var regResponse response

	if !registerResult {
		regResponse.Message = "Register failed"
	} else {
		regResponse.Message = "Register success"
	}

	json.NewEncoder(w).Encode(regResponse)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenHeader := r.Header.Get("token")
	emailHeader := r.Header.Get("email")

	profile, err := userService.ViewProfile(emailHeader, tokenHeader)
	if err != nil {
		log.Println("Failed to view profile,    ", err)
	}

	profile.Password = "-"

	json.NewEncoder(w).Encode(profile)
}
