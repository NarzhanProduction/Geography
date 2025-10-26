package request_models

type UserCreateRequest struct {
	Name     string
	Email    string
	Password string
}

type UserLoginRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginResponse struct {
	Token string `json:"token"`
}
