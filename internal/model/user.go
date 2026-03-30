package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"size:16;not null" json:"role"` // "admin" | "super_admin"
}

const (
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"
)
