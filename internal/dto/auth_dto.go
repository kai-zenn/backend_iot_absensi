package dto

import (
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type LoginRequest struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

type AuthenticatedUser struct {
  ID   uuid.UUID `json:"id"`
  Name string `json:"name"`
  Role entity.UserRole `json:"role"`
}

type LoginResponse struct {
  Token string `json:"token"`
  User  AuthenticatedUser `json:"user"`
}
