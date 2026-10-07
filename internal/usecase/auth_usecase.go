package usecase

import (
    "context"
    "github.com/kai-zenn/backend_iot_absensi/internal/dto"
)

type IAuthUsecase interface {
  Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
  Logout(ctx context.Context, tokenHash string) error
  ValidateSession(ctx context.Context, rawToken string) (dto.AuthenticatedUser, error)
}
