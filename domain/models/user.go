package models

import "time"

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Role      UserRole  `json:"role" gorm:"default:'customer'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleSeller   UserRole = "seller"
	RoleCustomer UserRole = "customer"
)
