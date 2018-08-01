package repo

import (
	"time"
)

// User is user_auth's model
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Msisdn    string    `json:"msisdn" db:"msisdn"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	Status    int       `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole is user_role's model
type UserRole struct {
	RoleID    string    `json:"role_id" db:"id"`
	UserID    string    `json:"id" db:"user_id"`
	Role      int       `json:"role" db:"role"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
