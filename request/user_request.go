package request

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type Response struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Msisdn   string `json:"msisdn"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Role     int    `json:"role"`
}
