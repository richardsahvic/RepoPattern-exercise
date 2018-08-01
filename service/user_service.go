package service

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"repo"

	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repo.UserRepository
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

func (s *userService) Login(username string, password string) (valid bool, err error) {
	valid = true
	result, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
		} else {
			log.Println("Error at Login", err)
		}
		valid = false
		return
	}

	match := CheckPasswordHash(password, result.Password)
	if !match {
		valid = false
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

func (s *userService) ViewProfile(email string) (userProfile repo.User, err error) {
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
