package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}
