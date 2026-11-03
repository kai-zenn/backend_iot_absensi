package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/dto"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	jwtpkg "github.com/kai-zenn/backend_iot_absensi/internal/pkgs/jwt"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"
)

type BcryptInterface interface {
	GenerateHashPassword(password string) (string, error)
	CompareHashPassword(hashedPassword, password string) error
}

type JWTInterface interface {
	GenerateToken(userID uuid.UUID, role entity.UserRole) (string, error)
	VerifyToken(token string) (*jwtpkg.Claims, error)
}

type IAuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, tokenHash string) error
}

type AuthUsecase struct {
	userRepo    repository.IUserRepository
	bcrypt      BcryptInterface
	jwtAuth     JWTInterface
}

func NewAuthUseCase(
	ur repository.IUserRepository,
	ttl time.Duration,
	bcrypt BcryptInterface,
	jwt JWTInterface,
) IAuthUsecase {
	return &AuthUsecase{
		userRepo:    ur,
		bcrypt:      bcrypt,
		jwtAuth:     jwt,
	}
}

func (au *AuthUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := au.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("user tidak ada")
	}

	if user.PasswordHash == nil {
		return nil, errors.New("password belum diatur")
	}

	err = au.bcrypt.CompareHashPassword(*user.PasswordHash, req.Password)
	if err != nil {
		return nil, errors.New("password salah")
	}

	token, err := au.jwtAuth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.AuthenticatedUser{
		  ID: user.ID, 
			Name: user.Name, 
			Role: user.Role,
		},
	}, nil
}

func (au *AuthUsecase) Logout(ctx context.Context, tokenHash string) error {
	return nil
}
