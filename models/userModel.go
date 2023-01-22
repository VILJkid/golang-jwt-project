package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	User_id       *string `json:"user_id"`
	First_name    *string `json:"first_name" validate:"required, min=2, max=100"`
	Last_name     *string `json:"last_name" validate:"required, min=2, max=100"`
	Password      *string `json:"password" validate:"required, min=6"`
	Email         *string `json:"email" validate:"required, email"`
	Phone         *string `json:"phone" validate:"required"`
	Token         *string `json:"token"`
	User_type     *string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Refresh_token *string `json:"refresh_token"`
}
