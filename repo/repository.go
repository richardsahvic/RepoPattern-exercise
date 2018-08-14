package repo

// UserRepository is a contract to persist user data to database
type UserRepository interface {
	FindProfiles() ([]User, error)
	FindByID(id string) (User, error)
	FindByEmail(email string) (User, error)
	FindByMsisdn(msisdn string) (User, error)
	FindByUsername(usrname string) (User, error)
	FindUserRole(userID string) (UserRole, error)
	FindExactRole(userID string, role int) (UserRole, error)
	InsertNewUser(user User) (string, error)
	InsertToRole(userRole UserRole) (bool, error)
}
