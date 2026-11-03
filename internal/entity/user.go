package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleStaff   UserRole = "staff"
)

type User struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `gorm:"column:username"`
	PasswordHash *string   `gorm:"column:password_hash"`
	Name         string    `gorm:"column:name"`
	Role         UserRole  `gorm:"column:role"`
	IsActive     bool      `gorm:"column:is_active"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (User) TableName() string { return "users" }
