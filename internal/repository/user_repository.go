package repository

import (
  "context"
  "github.com/kai-zenn/backend_iot_absensi/internal/entity"
)

type IUserRepository interface {
  FindByUsername(ctx context.Context, username string) (*entity.User, error)
  FindByID(ctx context.Context, id uint64) (*entity.User, error)
  Create(ctx context.Context, u *entity.User) (id uint64, err error)
}

type ISessionRepository interface {
  Create(ctx context.Context, s *entity.Session) error
  FindByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error)
  Delete(ctx context.Context, tokenHash string) error
}
