package postgres

import (
	"github.com/google/uuid"
	"github.com/kai-zenn/backend_iot_absensi/internal/entity"
	"github.com/kai-zenn/backend_iot_absensi/internal/repository"

	"gorm.io/gorm"

	"context"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.IUserRepository {
  return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
  var user entity.User

  err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
  if err != nil {
    return nil, err
  }
  return &user, err
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
  var user entity.User

  err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, u *entity.User) (id uuid.UUID, err error) {
  if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
    return uuid.Nil, err
  }
  return u.ID, nil
}
