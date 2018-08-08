package repo

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	conn              *sqlx.DB
	findAllStmt       *sqlx.Stmt
	findEmailStmt     *sqlx.Stmt
	findMsisdnStmt    *sqlx.Stmt
	findUsrnameStmt   *sqlx.Stmt
	findRoleStmt      *sqlx.Stmt
	findExactRoleStmt *sqlx.Stmt
	insertUser        *sqlx.NamedStmt
	insertToRole      *sqlx.NamedStmt
}

func (db *userRepository) MustPrepareStmt(query string) *sqlx.Stmt {
	stmt, err := db.conn.Preparex(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %s\n", err)
		os.Exit(1)
	}
	return stmt
}

func (db *userRepository) MustPrepareNamedStmt(query string) *sqlx.NamedStmt {
	stmt, err := db.conn.PrepareNamed(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %s\n", err)
		os.Exit(1)
	}
	return stmt
}

// NewRepository create new instance of UserRepository implementation
func NewRepository(db *sqlx.DB) UserRepository {
	r := userRepository{conn: db}
	r.findAllStmt = r.MustPrepareStmt("SELECT * FROM user_auth")
	r.findMsisdnStmt = r.MustPrepareStmt("SELECT * FROM user_auth WHERE msisdn=?")
	r.findEmailStmt = r.MustPrepareStmt("SELECT * FROM user_auth WHERE email=?")
	r.findUsrnameStmt = r.MustPrepareStmt("SELECT * FROM user_auth WHERE username=?")
	r.findRoleStmt = r.MustPrepareStmt("SELECT * FROM user_role WHERE user_id =?")
	r.findExactRoleStmt = r.MustPrepareStmt("SELECT * FROM user_role WHERE user_id =? AND role=?")
	r.insertUser = r.MustPrepareNamedStmt("INSERT INTO user_auth (id, email, msisdn, username, password, status) VALUES (:id, :email, :msisdn, :username, :password, :status)")
	r.insertToRole = r.MustPrepareNamedStmt("INSERT INTO user_role (id, user_id, role) VALUES (:id, :user_id, :role)")
	return &r
}

func (db *userRepository) FindProfiles() (usr []User, err error) {
	err = db.findAllStmt.Select(&usr)
	if err != nil {
		log.Println("Error at finding profiles,    ", err)
	}
	return
}

func (db *userRepository) FindByMsisdn(msisdn string) (usr User, err error) {
	err = db.findMsisdnStmt.Get(&usr, msisdn)
	if err != nil {
		log.Println("Error at Finding phone number,    ", err)
	}
	return
}

func (db *userRepository) FindByEmail(email string) (usr User, err error) {
	var user []User
	err = db.findEmailStmt.Select(&user, email)
	if err != nil {
		log.Println("Error at finding email,    ", err)
	}

	if len(user) != 0 {
		usr = user[0]
	}

	return
}

func (db *userRepository) FindByUsername(usrname string) (usr User, err error) {
	var u User
	err = db.findUsrnameStmt.Get(&u, usrname)
	usr = u
	if err != nil {
		log.Println("Error at finding username,    ", err)
	}
	return
}

func (db *userRepository) FindUserRole(userID string) (userRole UserRole, err error) {
	err = db.findRoleStmt.Get(&userRole, userID)
	if err != nil {
		log.Println("Error while finding user role,    ", err)
	}
	return
}

func (db *userRepository) FindExactRole(userID string, role int) (userRole UserRole, err error) {
	err = db.findExactRoleStmt.Get(&userRole, userID, role)
	if err != nil {
		log.Println("Error at finding exat user role,    ", err)
	}
	return
}

func (db *userRepository) InsertNewUser(user User) (lastID string, err error) {
	_, err = db.insertUser.Exec(user)
	if err != nil {
		log.Println("Error inserting new user,    ", err)
	}

	lastID = user.ID

	return
}

func (db *userRepository) InsertToRole(newRole UserRole) (success bool, err error) {
	_, err = db.insertToRole.Exec(newRole)
	if err != nil {
		log.Println("Error inserting new role,    ", err)
	} else {
		success = true
	}
	return
}
