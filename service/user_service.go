package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"repo"

	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repo.UserRepository
}

type Token struct {
	jwt.StandardClaims
	Role int `json:"role"`
}

var mySigningKey []byte

func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

// NewUserService create new instance of UserService implementation
func NewUserService(userRepo repo.UserRepository) UserService {
	s := userService{userRepo: userRepo}
	return &s
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *userService) Login(username string, password string, role int) (token string, err error) {
	mySigningKey := []byte("IDKWhatThisIs")

	userData, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
		} else {
			log.Println("Error at finding user's data", err)
		}
		return
	}

	match := CheckPasswordHash(password, userData.Password)
	if !match {
		log.Println("Wrong password")
	}

	loginRole, err := s.userRepo.FindExactRole(userData.ID, role)
	if len(loginRole.RoleID) == 0 {
		log.Println("User has no such role")
		return
	}

	claims := Token{
		jwt.StandardClaims{
			Subject:   userData.ID,
			ExpiresAt: time.Now().Add(15 * time.Second).Unix(),
		},
		role,
	}

	signing := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ = signing.SignedString(mySigningKey)
	if len(token) == 0 {
		log.Println("Failed to generate token")
		return
	}

	return
}

func (s *userService) Register(userRegister repo.User, role int) (registered bool, err error) {
	registered = false

	reEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	emailValid := reEmail.MatchString(userRegister.Email)
	if !emailValid {
		log.Println("Email format is not valid.")
		return
	}

	checkEmail, err := s.userRepo.FindByEmail(userRegister.Email)
	if len(checkEmail.Email) != 0 {
		checkRole, err := s.userRepo.FindUserRole(checkEmail.ID)
		if checkRole.Role == role {
			registered = false
			log.Println("User registered with an existing role,    ", err)
			return registered, err
		} else if userRegister.Username == checkEmail.Username {
			node, err := snowflake.NewNode(1)
			if err != nil {
				fmt.Println("Fail to generate snowflake id,    ", err)
				return registered, err
			}

			id := node.Generate().String()
			newRole := repo.UserRole{
				RoleID: id,
				UserID: checkRole.UserID,
				Role:   role,
			}
			registered, err = s.userRepo.InsertToRole(newRole)
			return registered, err
		}
	}

	checkUsername, err := s.userRepo.FindByUsername(userRegister.Username)
	if len(checkUsername.Username) != 0 {
		registered = false
		log.Println("Username exist on another account,    ", err)
		return
	}

	checkMsisdn, err := s.userRepo.FindByMsisdn(userRegister.Msisdn)
	if len(checkMsisdn.Msisdn) != 0 {
		registered = false
		log.Println("Phone number exist on another account,   ", err)
		return
	}

	userRegister.Password, err = HashPassword(userRegister.Password)
	if err != nil {
		log.Println("Failed encrypting password,  ", err)
		return
	}

	_, err = s.userRepo.InsertNewUser(userRegister)
	if err != nil {
		log.Println("Failed registering,    ", err)
		return
	} else {
		registered = true
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("Failed generating snowflake id,    ", err)
		return registered, err
	}
	id := node.Generate().String()

	newInsertRole := repo.UserRole{
		RoleID: id,
		UserID: userRegister.ID,
		Role:   role,
	}

	_, err = s.userRepo.InsertToRole(newInsertRole)
	if err != nil {
		log.Println("Failed registering new role by request,    ", err)
		return
	} else {
		registered = true
	}

	return
}

func (s *userService) ViewProfile(email string, token string) (userProfile repo.User, err error) {
	at(time.Unix(0, 0), func() {
		tokenClaims, err := jwt.ParseWithClaims(token, &Token{}, func(tokenClaims *jwt.Token) (interface{}, error) {
			return []byte("IDKWhatThisIs"), nil
		})

		if claims, _ := tokenClaims.Claims.(*Token); claims.ExpiresAt > time.Now().Unix() {
			fmt.Printf("%v %v", claims.Role, claims.StandardClaims.ExpiresAt)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	})

	reEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	emailValid := reEmail.MatchString(email)
	if !emailValid {
		return
	}

	userProfile, err = s.userRepo.FindByEmail(email)
	if err != nil {
		log.Println("Error at finding user's profile,	", err)
	}
	return
}
