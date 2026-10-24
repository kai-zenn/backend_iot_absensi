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

func (ur *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
  var user entity.User

  err := ur.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
  if err != nil {
    return nil, err
  }
  return &user, err
}

func (ur *userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
  var user entity.User

  err := ur.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (ur *userRepositoryImpl) Create(ctx context.Context, u *entity.User) (id uuid.UUID, err error) {
  if err := ur.db.WithContext(ctx).Create(u).Error; err != nil {
    return uuid.Nil, err
  }
  return u.ID, nil
}

// Storing User Session 
type sessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repository.ISessionRepository {
  return &sessionRepositoryImpl{db: db}
}

func (ser *sessionRepositoryImpl) Create(ctx context.Context, s *entity.Session) error {
  return ser.db.WithContext(ctx).Create(entity.Session{}).Error
}

func (ser *sessionRepositoryImpl)	FindByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error){
  var session entity.Session

  err := ser.db.WithContext(ctx).Where("token_has = ?", tokenHash).First(&session).Error
  if err != nil {
    return nil, err
  }

  return &session, nil
}

func (ser *sessionRepositoryImpl)	Delete(ctx context.Context, tokenHash string) error {
  return ser.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Delete(&entity.Session{}).Error
}
