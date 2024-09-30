package dto

type AuthDTO struct {
	UserID          string
	UserEmail       string
	UserName        string
	UserCurrentRole string
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=100"`
	AsRole   string `json:"as_role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	AsRole   string `json:"as_role"`
}

type LoginDTO struct {
	AccessToken             string `json:"access_token"`
	RefreshToken            string `json:"-"`
	RefreshTokenExpiredUnix int64  `json:"-"`
}
