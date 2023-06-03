package internal

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Platform int    `json:"platform" validate:"required,oneof=1 2"` // 1.web端 2.平台端(安卓、ios、pc)
}
