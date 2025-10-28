package request_models

type UserCreateRequest struct {
	Name     string
	Email    string
	Password string
}

type UserLoginRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
