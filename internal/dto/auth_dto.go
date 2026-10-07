package dto

import "github.com/google/uuid"

type LoginRequest struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

type AuthenticatedUser struct {
  ID   uuid.UUID `json:"id"`
  Name string `json:"name"`
  Role string `json:"role"` // ADMIN | STAFF
}

type LoginResponse struct {
  User AuthenticatedUser `json:"user"`
}
